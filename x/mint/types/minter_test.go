package types

import (
	"flag"
	"fmt"
	StakingExported "github.com/KuChainNetwork/kuchain/x/staking/exported"
	"log"
	"math/rand"
	"os"
	"runtime"
	"testing"

	chainTypes "github.com/KuChainNetwork/kuchain/chain/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestNextInflation(t *testing.T) {
	minter := InitialMinter(sdk.NewDecWithPrec(9, 2))
	params := NewParams(StakingExported.DefaultBondDenom,
		sdk.NewDecWithPrec(9, 2),
		sdk.NewDecWithPrec(12, 2),
		sdk.NewDecWithPrec(6, 2),
		sdk.NewDecWithPrec(67, 2),
		uint64(60*60*8766/3),
	)
	blocksPerYr := chainTypes.NewDec(int64(params.BlocksPerYear))

	// Governing Mechanism:
	//    inflationRateChangePerYear = (1- BondedRatio/ GoalBonded) * MaxInflationRateChange

	tests := []struct {
		bondedRatio, setInflation, expChange sdk.Dec
	}{
		//with 0% bonded atom kusupply the inflation should increase by InflationRateChange
		{sdk.ZeroDec(), sdk.NewDecWithPrec(6, 2), params.InflationRateChange.Quo(blocksPerYr)},

		// 100% bonded, starting at 20% inflation and being reduced
		// (1 - (1/0.67))*(0.13/8667)
		{sdk.OneDec(), sdk.NewDecWithPrec(12, 2),
			sdk.OneDec().Sub(sdk.OneDec().Quo(params.GoalBonded)).Mul(params.InflationRateChange).Quo(blocksPerYr)},

		// 50% bonded, starting at 10% inflation and being increased
		{sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(10, 2),
			sdk.OneDec().Sub(sdk.NewDecWithPrec(5, 1).Quo(params.GoalBonded)).Mul(params.InflationRateChange).Quo(blocksPerYr)},

		//test 6% minimum stop (testing with 100% bonded)
		{sdk.OneDec(), sdk.NewDecWithPrec(6, 2), sdk.ZeroDec()},
		{sdk.OneDec(), sdk.NewDecWithPrec(600000001, 10), sdk.NewDecWithPrec(-1, 10)},

		// test 20% maximum stop (testing with 0% bonded)
		{sdk.ZeroDec(), sdk.NewDecWithPrec(12, 2), sdk.ZeroDec()},
		{sdk.ZeroDec(), sdk.NewDecWithPrec(1199999999, 10), sdk.NewDecWithPrec(1, 10)},

		// perfect balance shouldn't change inflation
		{sdk.NewDecWithPrec(67, 2), sdk.NewDecWithPrec(12, 2), sdk.ZeroDec()},
	}
	for i, tc := range tests {
		minter.Inflation = tc.setInflation

		inflation := minter.NextInflationRate(params, tc.bondedRatio)
		diffInflation := inflation.Sub(tc.setInflation)

		require.True(t, diffInflation.Equal(tc.expChange),
			"Test Index: %v\nDiff:  %v\nExpected: %v\n", i, diffInflation, tc.expChange)
	}
}

func TestBlockProvision(t *testing.T) {
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := DefaultParams()

	secondsPerYear := int64(60 * 60 * 8766)

	tests := []struct {
		annualProvisions int64
		expProvisions    int64
	}{
		{secondsPerYear / 3, 1},
		{secondsPerYear/3 + 1, 1},
		{(secondsPerYear / 3) * 2, 2},
		{(secondsPerYear / 3) / 2, 0},
	}
	for i, tc := range tests {
		minter.AnnualProvisions = chainTypes.NewDec(tc.annualProvisions)
		provisions := minter.BlockProvision(params)

		expProvisions := chainTypes.NewCoin(params.MintDenom,
			chainTypes.NewInt(tc.expProvisions))

		require.True(t, expProvisions.IsEqual(provisions),
			"test: %v\n\tExp: %v\n\tGot: %v\n",
			i, tc.expProvisions, provisions)
	}
}

// Benchmarking :)
// previously using sdk.Int operations:
// BenchmarkBlockProvision-4 5000000 220 ns/op
//
// using sdk.Dec operations: (current implementation)
// BenchmarkBlockProvision-4 3000000 429 ns/op
func BenchmarkBlockProvision(b *testing.B) {
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := DefaultParams()

	s1 := rand.NewSource(100)
	r1 := rand.New(s1)
	minter.AnnualProvisions = chainTypes.NewDec(r1.Int63n(1000000))

	// run the BlockProvision function b.N times
	for n := 0; n < b.N; n++ {
		minter.BlockProvision(params)
	}
}

// Next inflation benchmarking
// BenchmarkNextInflation-4 1000000 1828 ns/op
func BenchmarkNextInflation(b *testing.B) {
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := DefaultParams()
	bondedRatio := sdk.NewDecWithPrec(1, 1)

	// run the NextInflationRate function b.N times
	for n := 0; n < b.N; n++ {
		minter.NextInflationRate(params, bondedRatio)
	}

}

// Next annual provisions benchmarking
// BenchmarkNextAnnualProvisions-4 5000000 251 ns/op
func BenchmarkNextAnnualProvisions(b *testing.B) {
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := DefaultParams()
	totalSupply := chainTypes.NewInt(100000000000000)

	// run the NextAnnualProvisions function b.N times
	for n := 0; n < b.N; n++ {
		minter.NextAnnualProvisions(params, totalSupply)
	}

}

func BenchmarkTestMinting(t *testing.B) {
	var (
		logFileName = flag.String("log", "BenchmarkTestMinting.log", "Log file name")
	)

	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	//set logfile Stdout
	logFile, logErr := os.OpenFile(*logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if logErr != nil {
		fmt.Println("Fail to find", *logFile, "cServer start Failed")
		os.Exit(1)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	minter := DefaultInitialMinter()
	params := NewParams(StakingExported.DefaultBondDenom,
		sdk.NewDecWithPrec(14, 2),
		sdk.NewDecWithPrec(21, 2),
		sdk.NewDecWithPrec(8, 2),
		sdk.NewDecWithPrec(67, 2),
		uint64(60*60*8766/3),
	)
	tests := []struct {
		bondedRatio sdk.Dec
	}{
		{sdk.ZeroDec(),},
		{sdk.NewDecWithPrec(33, 2),},
		{sdk.NewDecWithPrec(5, 1),},
		{sdk.NewDecWithPrec(67, 2),},
		{sdk.OneDec(),},
	}
	for _, tc := range tests {
		totalStakingSupply := sdk.NewInt(140000000)
		minter.Inflation = params.InflationRateChange
		for i := uint64(0); i < params.BlocksPerYear; i++ {
			minter.Inflation = minter.NextInflationRate(params, tc.bondedRatio)
			minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)

			bAmount := minter.BlockProvision(params).Amount
			totalStakingSupply = totalStakingSupply.Add(bAmount)
			//write log
			if bAmount.BigInt().Int64() <= 0 {
				log.Println("bondedRatio", tc.bondedRatio, minter.Inflation, bAmount, totalStakingSupply)
			}
		}
		fmt.Println(totalStakingSupply)
	}
}

func TestMinting(t *testing.T) {
	minter := DefaultInitialMinter()
	params := DefaultParams()

	totalStakingSupply := sdk.NewInt(140000000)
	tests := []struct {
		inf sdk.Dec
	}{
		{params.InflationMin},
		{params.InflationMax},
	}
	for _, tc := range tests {
		for i := uint64(0); i < params.BlocksPerYear; i++ {
			minter.Inflation = tc.inf
			minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)
			totalStakingSupply = totalStakingSupply.Add(minter.BlockProvision(params).Amount)
		}
		fmt.Println(totalStakingSupply)
	}
}
