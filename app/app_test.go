package app

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"
	abci "github.com/tendermint/tendermint/abci/types"
)

/*
func TestKuchainAppExport(t *testing.T) {
	db := tmdb.NewMemDB()
	kuApp := NewKuchainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	err := setGenesis(kuApp)
	require.NoError(t, err)

	// Making a new app object with the db, so that init chain hasn't been called
	newKuApp := NewKuchainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)
	_, _, err = newKuApp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}
*/

/*
// ensure that black listed addresses are properly set in bank keeper
func TestBlackListedAddrs(t *testing.T) {
	db := tmdb.NewMemDB()
	kuApp := NewKuchainApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0)

	for acc := range maccPerms {
		require.True(t, kuApp.assetKeeper.BlacklistedAddr(kuApp.supplyKeeper.GetModuleAddress(acc)))
	}
}
*/
func setGenesis(kuApp *KuchainApp) error {
	genesisState := simapp.NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(kuApp.Codec(), genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	kuApp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	kuApp.Commit()
	return nil
}
