package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/slashing/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestSigningInfo(t *testing.T) {
	wallet := simapp.NewWallet()
	Convey("TestGetSetValidatorSigningInfo", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.SlashKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		pk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepqn4usdx22zdntysj7n795xj77wrc62sytheeevr7zlna4yhwppdrs8mpds3")
		info, found := keeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(pk.Address()))
		require.False(t, found)
		newInfo := types.NewValidatorSigningInfo(
			sdk.ConsAddress(pk.Address()),
			int64(4),
			int64(3),
			time.Unix(2, 0),
			false,
			int64(10),
		)
		keeper.SetValidatorSigningInfo(ctx, sdk.ConsAddress(pk.Address()), newInfo)
		info, found = keeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(pk.Address()))
		require.True(t, found)
		require.Equal(t, info.StartHeight, int64(4))
		require.Equal(t, info.IndexOffset, int64(3))
		require.Equal(t, info.JailedUntil, time.Unix(2, 0).UTC())
		require.Equal(t, info.MissedBlocksCounter, int64(10))
	})
	Convey("TestGetSetValidatorSigningInfo", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.SlashKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		pk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepqn4usdx22zdntysj7n795xj77wrc62sytheeevr7zlna4yhwppdrs8mpds3")
		missed := keeper.GetValidatorMissedBlockBitArray(ctx, sdk.ConsAddress(pk.Address()), 0)
		require.False(t, missed) // treat empty key as not missed
		keeper.SetValidatorMissedBlockBitArray(ctx, sdk.ConsAddress(pk.Address()), 0, true)
		missed = keeper.GetValidatorMissedBlockBitArray(ctx, sdk.ConsAddress(pk.Address()), 0)
		require.True(t, missed) // now should be missed
	})
	Convey("TestTombstoned", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.SlashKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		pk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepqn4usdx22zdntysj7n795xj77wrc62sytheeevr7zlna4yhwppdrs8mpds3")
		require.Panics(t, func() { keeper.Tombstone(ctx, sdk.ConsAddress(pk.Address())) })
		require.False(t, keeper.IsTombstoned(ctx, sdk.ConsAddress(pk.Address())))

		newInfo := types.NewValidatorSigningInfo(
			sdk.ConsAddress(pk.Address()),
			int64(4),
			int64(3),
			time.Unix(2, 0),
			false,
			int64(10),
		)
		keeper.SetValidatorSigningInfo(ctx, sdk.ConsAddress(pk.Address()), newInfo)

		require.False(t, keeper.IsTombstoned(ctx, sdk.ConsAddress(pk.Address())))
		keeper.Tombstone(ctx, sdk.ConsAddress(pk.Address()))
		require.True(t, keeper.IsTombstoned(ctx, sdk.ConsAddress(pk.Address())))
		require.Panics(t, func() { keeper.Tombstone(ctx, sdk.ConsAddress(pk.Address())) })
	})
	Convey("TestJailUntil", t, func() {
		_, _, _, _, _, _, app := NewTestApp(wallet)
		keeper := app.SlashKeeper()
		ctx := app.BaseApp.NewContext(true, abci.Header{Height: app.LastBlockHeight() + 1})
		pk, _ := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, "kuchainvalconspub1zcjduepqn4usdx22zdntysj7n795xj77wrc62sytheeevr7zlna4yhwppdrs8mpds3")
		require.Panics(t, func() { keeper.JailUntil(ctx, sdk.ConsAddress(pk.Address()), time.Now()) })

		newInfo := types.NewValidatorSigningInfo(
			sdk.ConsAddress(pk.Address()),
			int64(4),
			int64(3),
			time.Unix(2, 0),
			false,
			int64(10),
		)
		keeper.SetValidatorSigningInfo(ctx, sdk.ConsAddress(pk.Address()), newInfo)
		keeper.JailUntil(ctx, sdk.ConsAddress(pk.Address()), time.Unix(253402300799, 0).UTC())

		info, ok := keeper.GetValidatorSigningInfo(ctx, sdk.ConsAddress(pk.Address()))
		require.True(t, ok)
		require.Equal(t, time.Unix(253402300799, 0).UTC(), info.JailedUntil)
	})
}
