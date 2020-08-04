package keeper_test

import (
	"encoding/json"
	chainType "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/supply/exported"
	"github.com/KuChainNetwork/kuchain/x/supply/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"gopkg.in/yaml.v2"
	"strings"
	"testing"
)

func TestModuleAccountMarshalYAML(t *testing.T) {
	name := "test"
	moduleAcc := types.NewEmptyModuleAccount(name, types.Minter, types.Burner, types.Staking)
	app, ctx := createTestApp(false)

	mAccI := (app.AccountKeeper().NewAccount(ctx, moduleAcc)).(exported.ModuleAccountI) // set the account number

	address := chainType.AccAddress([]byte("cosmos18x0gmwnjyq1qw"))
	mAccI.SetAuth(address)

	app.SupplyKeeper().SetModuleAccount(ctx, mAccI)

	bs, err := yaml.Marshal(moduleAcc)
	require.NoError(t, err)

	want := "|\n  address: cosmos1vdhhxmt0wvcns7psvakhwmn209cnzuthn4xfwn\n  public_key: \"\"\n  account_number: 0\n  name: test\n  permissions:\n  - minter\n  - burner\n  - staking\n"

	require.Equal(t, want, string(bs))
}

func TestHasPermissions(t *testing.T) {
	name := "test"
	macc := types.NewEmptyModuleAccount(name, types.Staking, types.Minter, types.Burner)
	cases := []struct {
		permission string
		expectHas  bool
	}{
		{types.Staking, true},
		{types.Minter, true},
		{types.Burner, true},
		{"other", false},
	}

	for i, tc := range cases {
		hasPerm := macc.HasPermission(tc.permission)
		if tc.expectHas {
			require.True(t, hasPerm, "test case #%d", i)
		} else {
			require.False(t, hasPerm, "test case #%d", i)
		}
	}
}

func TestValidate(t *testing.T) {

	tests := []struct {
		name   string
		acc    exported.ModuleAccountI
		expErr error
	}{
		{
			"valid module account",
			types.NewEmptyModuleAccount("test", types.Minter),
			nil,
		},
	}
	app, ctx := createTestApp(false)
	for _, tt := range tests {
		tt := tt
		mAccI := (app.AccountKeeper().NewAccount(ctx, tt.acc)).(exported.ModuleAccountI) // set the account number
		addr := sdk.AccAddress(secp256k1.GenPrivKey().PubKey().Address())
		mAccI.SetAuth(addr)

		t.Run(tt.name, func(t *testing.T) {
			b := tt.acc.HasPermission(types.Minter)
			require.Equal(t, true, b)
		})
	}
}

func TestModuleAccountJSON(t *testing.T) {
	pubkey := secp256k1.GenPrivKey().PubKey()
	addr := sdk.AccAddress(pubkey.Address())

	baseAcc := types.NewEmptyModuleAccount("test", "burner")
	app, ctx := createTestApp(false)
	mAccI := (app.AccountKeeper().NewAccount(ctx, baseAcc)).(exported.ModuleAccountI) // set the account number
	mAccI.SetAuth(addr)

	bz, err := json.Marshal(mAccI)
	require.NoError(t, err)

	require.NotEqual(t, strings.Index(string(bz), mAccI.GetAddress().String()), -1)
	require.NotEqual(t, strings.Index(string(bz), mAccI.GetID().String()), -1)
	require.NotEqual(t, strings.Index(string(bz), mAccI.GetPermissions()[0]), -1)

	var a types.ModuleAccount
	require.NoError(t, json.Unmarshal(bz, &a))
	require.Equal(t, mAccI.String(), a.String())
}
