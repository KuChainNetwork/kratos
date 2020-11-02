package types

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/config"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCoinStoreKey2AccountID(t *testing.T) {
	accID := chainTypes.NewAccountIDFromName(chainTypes.MustName("abcdefg"))
	key := genCoinStoreKey(CoinStatStoreKeyPrefix, accID.StoreKey())

	t.Logf("key %v", key)

	accIDc := coinStoreKey2AccountID(CoinStatStoreKeyPrefix, key)

	t.Logf("accIDc %s -- %s", accIDc, accID)

	Convey("TestCoinStoreKey2AccountID", t, func() {
		So(accIDc.Eq(accID), ShouldBeTrue)
	})
}

func TestCoinStoreKey2AccountIDInAdd(t *testing.T) {
	config.SealChainConfig()

	accAddressStr := "kuchain1xmc2z728py4gtwpc7jgytsan0282ww883qtv07"
	accAddress, _ := sdk.AccAddressFromBech32(accAddressStr)
	accID := chainTypes.NewAccountIDFromAccAdd(accAddress)
	key := genCoinStoreKey(CoinStatStoreKeyPrefix, accID.StoreKey())

	t.Logf("key %v %v", key, accAddress)

	accIDc := coinStoreKey2AccountID(CoinStatStoreKeyPrefix, key)

	t.Logf("accIDc addrs %s --- %s", accIDc, accID)

	Convey("TestCoinStoreKey2AccountIDInAdd", t, func() {
		So(accIDc.Eq(accID), ShouldBeTrue)
	})
}
