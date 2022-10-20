package types

import (
	hmTyps "github.com/shibaswaparmy/heimdall/types"
)

const (
	QueryBalance = "balances"
)

// QueryBalanceParams defines the params for querying an account balance.
type QueryBalanceParams struct {
	Address hmTyps.HeimdallAddress
}

// NewQueryBalanceParams creates a new instance of QueryBalanceParams.
func NewQueryBalanceParams(addr hmTyps.HeimdallAddress) QueryBalanceParams {
	return QueryBalanceParams{Address: addr}
}
