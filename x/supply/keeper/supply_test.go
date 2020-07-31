package keeper

import (
	"fmt"
	"github.com/KuChainNetwork/kuchain/chain/constants"
	chaintypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/supply/types"
	"testing"

	"gopkg.in/yaml.v2"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/stretchr/testify/require"
)

func TestSupplyMarshalYAML(t *testing.T) {
	supply := types.DefaultSupply()
	coins := chaintypes.NewCoins(chaintypes.NewCoin(constants.DefaultBondDenom, sdk.OneInt()))
	supply.SetTotal(coins)

	bz, err := yaml.Marshal(supply)
	require.NoError(t, err)
	bzCoins, err := yaml.Marshal(coins)
	require.NoError(t, err)

	want := fmt.Sprintf(`total:
%s`, string(bzCoins))

	require.Equal(t, want, string(bz))
	require.Equal(t, want, supply.String())
}
