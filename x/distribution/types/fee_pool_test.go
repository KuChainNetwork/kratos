package types

import (
	"github.com/KuChainNetwork/kuchain/chain/constants"
	chaintypes "github.com/KuChainNetwork/kuchain/chain/types"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestValidateGenesis(t *testing.T) {

	fp := InitialFeePool()
	require.Nil(t, fp.ValidateGenesis())

	fp2 := FeePool{CommunityPool: chaintypes.DecCoins{{Denom: constants.DefaultBondDenom, Amount: sdk.NewDec(-1)}}}
	require.NotNil(t, fp2.ValidateGenesis())

}
