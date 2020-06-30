package keeper

// nolint:deadcode,unused
// DONTCOVER
// noalias

import (
	chaintype "github.com/KuChain-io/kuchain/chain/types"
	"github.com/KuChain-io/kuchain/x/slashing/external"
	"github.com/KuChain-io/kuchain/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
)

var (
	InitTokens = external.TokensFromConsensusPower(200)
)

// Have to change these parameters for tests
// lest the tests take forever
func TestParams() types.Params {
	params := types.DefaultParams()
	params.SignedBlocksWindow = 1000
	params.DowntimeJailDuration = 60 * 60
	return params
}

func NewTestMsgCreateValidator(acc chaintype.AccountID, pubKey crypto.PubKey, amt sdk.Int) external.StakingMsgCreateValidator {
	return external.StakingNewMsgCreateValidator(
		acc, pubKey,
		external.StakingDescription{}, sdk.ZeroDec(), acc,
	)
}

func NewTestMsgDelegate(delAddr chaintype.AccountID, valAddr chaintype.AccountID, delAmount sdk.Int) external.StakingMsgDelegate {
	amount := sdk.NewCoin(external.DefaultBondDenom, delAmount)
	return external.StakingNewMsgDelegate(delAddr, valAddr, amount)
}
