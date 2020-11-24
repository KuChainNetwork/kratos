package staking

import (
	"fmt"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	statkingexport "github.com/KuChainNetwork/kuchain/x/staking/exported"
	"github.com/KuChainNetwork/kuchain/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func genesisAppendValidator(ctx sdk.Context, keeper Keeper, isExport bool, validator types.Validator) {
	keeper.SetValidator(ctx, validator)

	// Manually set indices for the first time
	keeper.SetValidatorByConsAddr(ctx, validator)
	keeper.SetValidatorByPowerIndex(ctx, validator)

	// Call the creation hook if not exported
	if !isExport {
		keeper.AfterValidatorCreated(ctx, validator.OperatorAccount)
	}

	// update timeslice if necessary
	if validator.IsUnbonding() {
		keeper.InsertValidatorQueue(ctx, validator)
	}
}

// InitGenesis sets the pool and parameters for the provided keeper.  For each
// validator in data, it sets that validator in the keeper along with manually
// setting the indexes. In addition, it also sets any delegations found in
// data. Finally, it updates the bonded validators.
// Returns final validator set after applying all declaration and delegations
func InitGenesis(
	ctx sdk.Context, keeper Keeper, accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper, supplyKeeper types.SupplyKeeper, data types.GenesisState,
) (res []abci.ValidatorUpdate) {
	bondedTokens := sdk.ZeroInt()
	notBondedTokens := sdk.ZeroInt()

	// We need to pretend to be "n blocks before genesis", where "n" is the
	// validator update delay, so that e.g. slashing periods are correctly
	// initialized for the validator set e.g. with a one-block offset - the
	// first TM block is at height 1, so state updates applied from
	// genesis.json are in block 0.
	ctx = ctx.WithBlockHeight(1 - statkingexport.ValidatorUpdateDelay)

	keeper.SetParams(ctx, data.Params)
	keeper.SetLastTotalPower(ctx, data.LastTotalPower)

	for _, validator := range data.Validators {
		genesisAppendValidator(ctx, keeper, data.Exported, validator)

		switch validator.GetStatus() {
		case statkingexport.Bonded:
			bondedTokens = bondedTokens.Add(validator.GetTokens())
		case statkingexport.Unbonding, statkingexport.Unbonded:
			notBondedTokens = notBondedTokens.Add(validator.GetTokens())
		default:
			panic("invalid validator status")
		}
	}

	for _, delegation := range data.Delegations {
		// Call the before-creation hook if not exported
		if !data.Exported {
			keeper.BeforeDelegationCreated(ctx, delegation.DelegatorAccount, delegation.ValidatorAccount)
		}

		keeper.SetDelegation(ctx, delegation)

		// Call the after-modification hook if not exported
		if !data.Exported {
			keeper.AfterDelegationModified(ctx, delegation.DelegatorAccount, delegation.ValidatorAccount)
		}
	}

	for _, ubd := range data.UnbondingDelegations {
		keeper.SetUnbondingDelegation(ctx, ubd)
		for _, entry := range ubd.Entries {
			keeper.InsertUBDQueue(ctx, ubd, entry.CompletionTime)
			notBondedTokens = notBondedTokens.Add(entry.Balance)
		}
	}

	for _, red := range data.Redelegations {
		keeper.SetRedelegation(ctx, red)
		for _, entry := range red.Entries {
			keeper.InsertRedelegationQueue(ctx, red, entry.CompletionTime)
		}
	}

	bondedCoins := chainTypes.NewCoins(chainTypes.NewCoin(data.Params.BondDenom, bondedTokens))
	notBondedCoins := chainTypes.NewCoins(chainTypes.NewCoin(data.Params.BondDenom, notBondedTokens))

	// check if the unbonded and bonded pools accounts exists
	bondedPool := keeper.GetBondedPool(ctx)
	if bondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	// TODO remove with genesis 2-phases refactor https://github.com/cosmos/cosmos-sdk/issues/2862
	// add coins if not provided on genesis
	if bankKeeper.GetCoinPowers(ctx, bondedPool.GetID()).IsZero() {
		ctx.Logger().Info("genesis module account asset", "account", bondedPool.GetID(), "asset", bondedCoins)
		supplyKeeper.SetModuleAccount(ctx, bondedPool)
	}

	notBondedPool := keeper.GetNotBondedPool(ctx)
	if notBondedPool == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	if bankKeeper.GetCoinPowers(ctx, notBondedPool.GetID()).IsZero() {
		ctx.Logger().Info("genesis module account asset", "account", notBondedPool.GetID(), "asset", notBondedCoins)
		supplyKeeper.SetModuleAccount(ctx, notBondedPool)
	}

	// check if the module account exists
	ctx.Logger().Info("genesis module account", "name", types.ModuleAccountName)
	if err := supplyKeeper.InitModuleAccount(ctx, types.ModuleName); err != nil {
		panic(err)
	}

	// don't need to run Tendermint updates if we exported
	if data.Exported {
		for _, lv := range data.LastValidatorPowers {
			valAccountID := chainTypes.NewAccountIDFromAccAdd(sdk.AccAddress(lv.Address))
			keeper.SetLastValidatorPower(ctx, valAccountID, lv.Power)
			validator, found := keeper.GetValidator(ctx, valAccountID)
			if !found {
				panic(fmt.Sprintf("validator %s not found", lv.Address))
			}
			update := validator.ABCIValidatorUpdate()
			update.Power = lv.Power // keep the next-val-set offset, use the last power for the first block
			res = append(res, update)
		}
	} else {
		res = keeper.ApplyAndReturnValidatorSetUpdates(ctx)
	}

	return res
}

// ExportGenesis returns a GenesisState for a given context and keeper. The
// GenesisState will contain the pool, params, validators, and bonds found in
// the keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	lastTotalPower := keeper.GetLastTotalPower(ctx)
	validators := keeper.GetAllValidators(ctx)
	delegations := keeper.GetAllDelegations(ctx)
	var unbondingDelegations []types.UnbondingDelegation
	keeper.IterateUnbondingDelegations(ctx, func(_ int64, ubd types.UnbondingDelegation) (stop bool) {
		unbondingDelegations = append(unbondingDelegations, ubd)
		return false
	})
	var redelegations []types.Redelegation
	keeper.IterateRedelegations(ctx, func(_ int64, red types.Redelegation) (stop bool) {
		redelegations = append(redelegations, red)
		return false
	})
	var lastValidatorPowers []types.LastValidatorPower
	keeper.IterateLastValidatorPowers(ctx, func(addr chainTypes.AccountID, power int64) (stop bool) {
		accAddr, _ := addr.ToAccAddress()
		lastValidatorPowers = append(lastValidatorPowers, types.LastValidatorPower{Address: sdk.ValAddress(accAddr), Power: power})
		return false
	})

	return types.GenesisState{
		Params:               params,
		LastTotalPower:       lastTotalPower,
		LastValidatorPowers:  lastValidatorPowers,
		Validators:           validators,
		Delegations:          delegations,
		UnbondingDelegations: unbondingDelegations,
		Redelegations:        redelegations,
		Exported:             true,
	}
}

// WriteValidators returns a slice of bonded genesis validators.
func WriteValidators(ctx sdk.Context, keeper Keeper) (vals []tmtypes.GenesisValidator) {
	keeper.IterateLastValidators(ctx, func(_ int64, validator statkingexport.ValidatorI) (stop bool) {
		vals = append(vals, tmtypes.GenesisValidator{
			PubKey: validator.GetConsPubKey(),
			Power:  validator.GetConsensusPower(),
			Name:   validator.GetMoniker(),
		})

		return false
	})

	return
}
