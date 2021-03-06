package ethbridge

import (
	"strings"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
	"github.com/swishlabsco/cosmos-ethereum-bridge/x/ethbridge/types"
	keeperLib "github.com/swishlabsco/cosmos-ethereum-bridge/x/oracle/keeper"
)

func TestBasicMsgs(t *testing.T) {
	//Setup
	cdc := codec.New()
	ctx, _, keeper := keeperLib.CreateTestKeepers(t, false, 0, nil)
	handler := NewHandler(keeper, cdc, types.DefaultCodespace)

	//Unrecognized type
	res := handler(ctx, sdk.NewTestMsg())
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Unrecognized ethbridge message type: "))

	//Normal Creation
	normalCreateMsg := types.CreateTestEthMsg(t)
	res = handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())

	//Bad Creation
	badCreateMsg := types.CreateTestEthMsg(t)
	badCreateMsg.Nonce = -1
	res = handler(ctx, badCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "invalid ethereum nonce provided"))

	badCreateMsg = types.CreateTestEthMsg(t)
	badCreateMsg.EthereumSender = "badAddress"
	res = handler(ctx, badCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "invalid ethereum address provided"))
}

func TestDuplicateMsgs(t *testing.T) {
	//Setup
	cdc := codec.New()
	ctx, _, keeper := keeperLib.CreateTestKeepers(t, false, 0, nil)
	handler := NewHandler(keeper, cdc, types.DefaultCodespace)
	normalCreateMsg := types.CreateTestEthMsg(t)
	res := handler(ctx, normalCreateMsg)
	require.True(t, res.IsOK())
	require.Equal(t, res.Log, "pending")

	//Duplicate message from same validator
	res = handler(ctx, normalCreateMsg)
	require.False(t, res.IsOK())
	require.True(t, strings.Contains(res.Log, "Not yet implemented"))

}
