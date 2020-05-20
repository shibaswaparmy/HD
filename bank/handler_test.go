package bank_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bank"
	"github.com/maticnetwork/heimdall/bank/types"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	handler        sdk.Handler
	contractCaller mocks.IContractCaller
}

// SetupTest setup all necessary things for querier tesing
func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)

	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = bank.NewHandler(suite.app.BankKeeper, &suite.contractCaller)
}

// TestHandlerTestSuite
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleMsgUnknown() {
	t, _, ctx := suite.T(), suite.app, suite.ctx

	result := suite.handler(ctx, nil)
	require.False(t, result.IsOK())
}

func (suite *HandlerTestSuite) TestHandlerMsgSend() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	amount := int64(10000000)
	from := hmTypes.HexToHeimdallAddress("123")
	app.BankKeeper.AddCoins(ctx, from, sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(amount*10))))
	msgSend := types.NewMsgSend(
		from,
		hmTypes.HexToHeimdallAddress("456"),
		sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(amount))),
	)
	result := suite.handler(ctx, msgSend)
	require.True(t, result.IsOK(), "Expected New msg to be sent")

}
