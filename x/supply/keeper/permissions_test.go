package keeper_test

import (
	"github.com/KuChainNetwork/kuchain/x/account/types"
	"github.com/KuChainNetwork/kuchain/x/supply/exported"
	types2 "github.com/KuChainNetwork/kuchain/x/supply/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHasPermission(t *testing.T) {
	emptyPermAddr := types2.NewPermissionsForAddress("empty", []string{})
	has := emptyPermAddr.HasPermission(types2.Minter)
	require.False(t, has)

	cases := []struct {
		permission string
		expectHas  bool
	}{
		{types2.Minter, true},
		{types2.Burner, true},
		{types2.Staking, true},
		{"random", false},
		{"", false},
	}
	permAddr := types2.NewPermissionsForAddress("test", []string{types2.Minter, types2.Burner, types2.Staking})
	for i, tc := range cases {
		has = permAddr.HasPermission(tc.permission)
		require.Equal(t, tc.expectHas, has, "test case #%d", i)
	}

}

func TestPermissions(t *testing.T) {
	cases := []struct {
		name        string
		permissions []string
		expectPass  bool
	}{
		{"no permissions", []string{}, true},
		{"valid permission", []string{types2.Minter}, true},
		//{"invalid permission", []string{""}, false},
		{"invalid and valid permission", []string{types2.Staking}, true},
	}

	app, ctx := createTestApp(false)

	for i, tc := range cases {
		acc := types.NewEmptyModuleAccount("test", tc.permissions...)
		mAccI := (app.AccountKeeper().NewAccount(ctx, acc)).(exported.ModuleAccountI) // set the account number
		addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		mAccI.SetAuth(addr)

		i, tc := i, tc
		t.Run(tc.name, func(t *testing.T) {
			err := app.SupplyKeeper().ValidatePermissions(mAccI)
			if tc.expectPass {
				require.NoError(t, err, "test case #%d", i)
			} else {
				require.Error(t, err, "test case #%d", i)
			}
		})
	}
}
