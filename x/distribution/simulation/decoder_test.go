package simulation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/distribution/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/tendermint/tendermint/crypto/ed25519"
)

var (
	delPk1    = ed25519.GenPrivKey().PubKey()
	delAddr1  = chainTypes.NewAccountIDFromByte(delPk1.Address())
	valAddr1  = chainTypes.NewAccountIDFromByte(delPk1.Address())
	consAddr1 = chainTypes.NewAccountIDFromByte(delPk1.Address().Bytes())
)

func makeTestCodec() (cdc *codec.Codec) {
	cdc = codec.New()
	sdk.RegisterCodec(cdc)
	codec.RegisterCrypto(cdc)
	types.RegisterCodec(cdc)
	return
}

func TestDecodeDistributionStore(t *testing.T) {
	cdc := makeTestCodec()

	decCoins := chainTypes.NewDecCoins(chainTypes.NewDecCoin(constants.DefaultBondDenom, sdk.NewInt(sdk.OneDec().Int64())))
	feePool := types.InitialFeePool()
	feePool.CommunityPool = decCoins
	info := types.NewDelegatorStartingInfo(2, sdk.OneDec(), 200)
	outstanding := types.ValidatorOutstandingRewards{Rewards: decCoins}
	commission := types.ValidatorAccumulatedCommission{Commission: decCoins}
	historicalRewards := types.NewValidatorHistoricalRewards(decCoins, 100)
	currentRewards := types.NewValidatorCurrentRewards(decCoins, 5)
	slashEvent := types.NewValidatorSlashEvent(10, sdk.OneDec())

	kvPairs := kv.Pairs{
		kv.Pair{Key: types.FeePoolKey, Value: cdc.MustMarshalBinaryLengthPrefixed(feePool)},
		kv.Pair{Key: types.ProposerKey, Value: consAddr1.Bytes()},
		kv.Pair{Key: types.GetValidatorOutstandingRewardsKey(valAddr1), Value: cdc.MustMarshalBinaryLengthPrefixed(outstanding)},
		kv.Pair{Key: types.GetDelegatorWithdrawAddrKey(delAddr1), Value: delAddr1.Bytes()},
		kv.Pair{Key: types.GetDelegatorStartingInfoKey(valAddr1, delAddr1), Value: cdc.MustMarshalBinaryLengthPrefixed(info)},
		kv.Pair{Key: types.GetValidatorHistoricalRewardsKey(valAddr1, 100), Value: cdc.MustMarshalBinaryLengthPrefixed(historicalRewards)},
		kv.Pair{Key: types.GetValidatorCurrentRewardsKey(valAddr1), Value: cdc.MustMarshalBinaryLengthPrefixed(currentRewards)},
		kv.Pair{Key: types.GetValidatorAccumulatedCommissionKey(valAddr1), Value: cdc.MustMarshalBinaryLengthPrefixed(commission)},
		kv.Pair{Key: types.GetValidatorSlashEventKeyPrefix(valAddr1, 13), Value: cdc.MustMarshalBinaryLengthPrefixed(slashEvent)},
		kv.Pair{Key: []byte{0x99}, Value: []byte{0x99}},
	}

	tests := []struct {
		name        string
		expectedLog string
	}{
		{"FeePool", fmt.Sprintf("%v\n%v", feePool, feePool)},
		{"Proposer", fmt.Sprintf("%v\n%v", consAddr1, consAddr1)},
		{"ValidatorOutstandingRewards", fmt.Sprintf("%v\n%v", outstanding, outstanding)},
		{"DelegatorWithdrawAddr", fmt.Sprintf("%v\n%v", delAddr1, delAddr1)},
		{"DelegatorStartingInfo", fmt.Sprintf("%v\n%v", info, info)},
		{"ValidatorHistoricalRewards", fmt.Sprintf("%v\n%v", historicalRewards, historicalRewards)},
		{"ValidatorCurrentRewards", fmt.Sprintf("%v\n%v", currentRewards, currentRewards)},
		{"ValidatorAccumulatedCommission", fmt.Sprintf("%v\n%v", commission, commission)},
		{"ValidatorSlashEvent", fmt.Sprintf("{%v %v}\n{%v %v}", slashEvent.ValidatorPeriod, slashEvent.Fraction, slashEvent.ValidatorPeriod, slashEvent.Fraction)},
		{"other", ""},
	}
	for i, tt := range tests {
		i, tt := i, tt
		t.Run(tt.name, func(t *testing.T) {
			switch i {
			case len(tests) - 1:
				require.Panics(t, func() { DecodeStore(cdc, kvPairs[i], kvPairs[i]) }, tt.name)
			default:
				require.Equal(t, tt.expectedLog, DecodeStore(cdc, kvPairs[i], kvPairs[i]), tt.name)
			}
		})
	}
}
