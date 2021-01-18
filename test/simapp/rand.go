package simapp

import (
	"math/rand"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
)

func RandCommonName(app *SimApp) types.AccountID {
	const (
		numsInName  = "0123456789"
		charsInName = "abcdefghijklmnopqrstuvwxyz"
		useInName   = numsInName + charsInName
	)

	// just use a fixed rand seed
	r := rand.New(rand.NewSource(app.RandSeed()))
	app.seed++

	// rand a 12 len name
	str := make([]byte, 0, 32)
	str = append(str, useInName[r.Intn(len(useInName))])
	for i := 1; i < constants.CommonAccountNameLen; i++ {
		str = append(str, useInName[r.Intn(len(useInName))])
	}

	id, err := types.NewAccountIDFromStr(string(str))
	if err != nil {
		panic(err)
	}

	return id
}
