package version

import (
	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	largeIn239 = chainTypes.MustAccountID("kukis.valida")
)

var (
	toLargeNumMap map[int64]bool
)

var (
	toLargeNums = [128]int64{
		27670, 35031, 42045, 239435, 556854, 557509, 565837, 567074, 236995, 402979, 657925, 43334,
		66704, 309472, 381238, 400717, 503001, 537823, 550991, 578261, 35281,
		239788, 288656, 557458, 558664, 559222, 60793, 240407,
		288707, 316611, 551002, 241083, 315926, 550750, 225078,
		238177, 317134, 357189, 567077, 286107, 425864, 27356,
		223503, 550951, 558715, 26958, 239384, 315977, 357588,
		546145, 245097, 286844, 378255, 402495, 557283, 238960,
		378617, 493893, 522680, 556905, 27009, 316560, 381713,
		558331, 240023, 286793, 314843, 315127, 566959, 238228, 378668, 400768, 225027, 356662, 381289,
		653428, 239972, 245148, 315076, 309523, 314792, 317083, 317295, 558335, 558894, 651955, 241134,
		315568, 385891, 504940, 558386, 555199, 558945, 317346, 331992, 402155, 553375, 237946, 384283,
		456427, 239011, 239737, 241876, 286224, 384653, 566060, 241825, 284687, 309426, 407858, 46088,
		237431, 243653, 359331, 378204, 554800, 554749, 27305, 422815, 553426, 557232, 26762, 559171,
		31483, 237046, 237380, 237895, 240458,
	}
)

func init() {
	toLargeNumMap = make(map[int64]bool, 128)
	for _, k := range toLargeNums {
		toLargeNumMap[k] = true
	}
}

func toNormalAccountID(acc chainTypes.AccountID) chainTypes.AccountID {
	val := make([]byte, 17)
	copy(val, acc.Value[:17])
	return chainTypes.AccountID{val}
}

func ProcessValidatorID(ctx sdk.Context, val types.Validator) {
	blockNum := ctx.BlockHeight()

	if blockNum > 657926 {
		val.OperatorAccount.Value = val.OperatorAccount.StoreKey()
		return
	}

	// Both In 239737 large [kukis.valida kukis.valida] <--> less [mathwalletbp]
	if blockNum == 239737 {
		if val.OperatorAccount.Eq(largeIn239) {
			val.OperatorAccount.Value = val.OperatorAccount.StoreKey()
		} else {
			val.OperatorAccount = toNormalAccountID(val.OperatorAccount)
		}
		return
	}

	if _, ok := toLargeNumMap[blockNum]; ok {
		val.OperatorAccount.Value = val.OperatorAccount.StoreKey()
	} else {
		val.OperatorAccount = toNormalAccountID(val.OperatorAccount)
	}
}
