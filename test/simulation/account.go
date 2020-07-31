package simulation

import (
	"math/rand"

	"github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RandomFees returns a random fee by selecting a random coin denomination and
// amount from the account's available balance. If the user doesn't have enough
// funds for paying fees, it returns empty coins.
func RandomFees(r *rand.Rand, ctx sdk.Context, spendableCoins types.Coins) (types.Coins, error) {
	if spendableCoins.Empty() {
		return nil, nil
	}

	denomIndex := r.Intn(len(spendableCoins))
	randCoin := spendableCoins[denomIndex]

	if randCoin.Amount.IsZero() {
		return nil, nil
	}

	amt, err := RandPositiveInt(r, randCoin.Amount)
	if err != nil {
		return nil, err
	}

	// Create a random fee and verify the fees are within the account's spendable
	// balance.
	fees := types.NewCoins(types.NewCoin(randCoin.Denom, amt))
	return fees, nil
}
