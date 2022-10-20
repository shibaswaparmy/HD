package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/shibaswaparmy/bor/crypto"
	"github.com/shibaswaparmy/bor/rlp"
)

const (
	// PulpHashLength pulp hash length
	PulpHashLength int = 4
)

// Pulp codec for RLP
type Pulp struct {
	typeInfos map[string]reflect.Type
}

var once sync.Once
var pulp *Pulp

// GetPulpInstance gets new pulp codec
func GetPulpInstance() *Pulp {
	once.Do(func() {
		pulp = NewPulp()
	})
	return pulp
}

// NewPulp creates new pulp codec
func NewPulp() *Pulp {
	p := &Pulp{}
	p.typeInfos = make(map[string]reflect.Type)
	return p
}

// GetPulpHash returns string hash
func GetPulpHash(msg sdk.Msg) []byte {
	return crypto.Keccak256([]byte(fmt.Sprintf("%s::%s", msg.Route(), msg.Type())))[:PulpHashLength]
}

// RegisterConcrete should be used to register concrete types that will appear in
// interface fields/elements to be encoded/decoded by pulp.
func (p *Pulp) RegisterConcrete(msg sdk.Msg) {
	rtype := reflect.TypeOf(msg)
	p.typeInfos[hex.EncodeToString(GetPulpHash(msg))] = rtype
}

// GetMsgTxInstance get new instance associated with base tx
func (p *Pulp) GetMsgTxInstance(hash []byte) interface{} {
	rtype := p.typeInfos[hex.EncodeToString(hash[:PulpHashLength])]
	return reflect.New(rtype).Elem().Interface().(sdk.Msg)
}

// EncodeToBytes encodes msg to bytes
func (p *Pulp) EncodeToBytes(tx StdTx) ([]byte, error) {
	msg := tx.GetMsgs()[0]
	txBytes, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return nil, err
	}

	return append(GetPulpHash(msg), txBytes[:]...), nil
}

// DecodeBytes decodes bytes to msg
func (p *Pulp) DecodeBytes(data []byte) (interface{}, error) {
	var txRaw StdTxRaw

	if len(data) <= PulpHashLength {
		return nil, errors.New("Invalid data length, should be greater than PulpPrefix")
	}

	if err := rlp.DecodeBytes(data[PulpHashLength:], &txRaw); err != nil {
		return nil, err
	}

	rtype := p.typeInfos[hex.EncodeToString(data[:PulpHashLength])]
	newMsg := reflect.New(rtype).Interface()
	if err := rlp.DecodeBytes(txRaw.Msg[:], newMsg); err != nil {
		return nil, err
	}

	// change pointer to non-pointer
	vptr := reflect.New(reflect.TypeOf(newMsg).Elem()).Elem()
	vptr.Set(reflect.ValueOf(newMsg).Elem())
	// return vptr.Interface(), nil

	result := StdTx{
		Msg:       vptr.Interface().(sdk.Msg),
		Signature: txRaw.Signature,
		Memo:      txRaw.Memo,
	}
	return result, nil
}
