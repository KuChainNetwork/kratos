package genutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/KuChainNetwork/kuchain/chain/client/txutil"
	"github.com/KuChainNetwork/kuchain/x/asset"
	"github.com/KuChainNetwork/kuchain/x/genutil/types"

	cfg "github.com/tendermint/tendermint/config"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/codec"
)

// GenAppStateFromConfig gets the genesis app state from the config
func GenAppStateFromConfig(cdc *codec.Codec, config *cfg.Config, initCfg InitConfig, genDoc tmtypes.GenesisDoc,
	genBalIterator types.GenesisBalancesIterator, stakingFuncManager types.StakingFuncManager,
) (appState json.RawMessage, err error) {

	// process genesis transactions, else create default genesis.json
	appGenTxs, persistentPeers, err := CollectStdTxs(cdc, config.Moniker, initCfg.GenTxsDir, genDoc,
		genBalIterator, stakingFuncManager)
	if err != nil {
		return appState, err
	}

	config.P2P.PersistentPeers = persistentPeers
	cfg.WriteConfigFile(filepath.Join(config.RootDir, "config", "config.toml"), config)

	// if there are no gen txs to be processed, return the default empty state
	if len(appGenTxs) == 0 {
		return appState, errors.New("there must be at least one genesis tx")
	}

	// create the app state
	appGenesisState, err := GenesisStateFromGenDoc(cdc, genDoc)
	if err != nil {
		return appState, err
	}

	appGenesisState, err = SetGenTxsInAppGenesisState(cdc, appGenesisState, appGenTxs)
	if err != nil {
		return appState, err
	}

	appState, err = codec.MarshalJSONIndent(cdc, appGenesisState)
	if err != nil {
		return appState, err
	}

	genDoc.AppState = appState
	err = ExportGenesisFile(&genDoc, config.GenesisFile())

	return appState, err
}

// CollectStdTxs processes and validates application's genesis StdTxs and returns
// the list of appGenTxs, and persistent peers required to generate genesis.json.
func CollectStdTxs(cdc *codec.Codec, moniker, genTxsDir string, genDoc tmtypes.GenesisDoc,
	genBalIterator types.GenesisBalancesIterator, stakingFuncManager types.StakingFuncManager,
) (appGenTxs []txutil.StdTx, persistentPeers string, err error) {

	var fos []os.FileInfo
	fos, err = ioutil.ReadDir(genTxsDir)
	if err != nil {
		return appGenTxs, persistentPeers, err
	}

	// prepare a map of all balances in genesis state to then validate
	// against the validators addresses
	var appState map[string]json.RawMessage
	if err := cdc.UnmarshalJSON(genDoc.AppState, &appState); err != nil {
		return appGenTxs, persistentPeers, err
	}

	balancesMap := make(map[string]asset.GenesisAsset)
	genBalIterator.IterateGenesisBalances(
		cdc, appState,
		func(balance asset.GenesisAsset) (stop bool) {
			balancesMap[balance.GetID().String()] = balance
			return false
		},
	)

	// addresses and IPs (and port) validator server info
	var addressesIPs []string

	for _, fo := range fos {
		filename := filepath.Join(genTxsDir, fo.Name())
		if !fo.IsDir() && (filepath.Ext(filename) != ".json") {
			continue
		}

		// get the genStdTx
		var jsonRawTx []byte
		if jsonRawTx, err = ioutil.ReadFile(filename); err != nil {
			return appGenTxs, persistentPeers, err
		}

		var genStdTx txutil.StdTx
		if err = cdc.UnmarshalJSON(jsonRawTx, &genStdTx); err != nil {
			return appGenTxs, persistentPeers, err
		}

		appGenTxs = append(appGenTxs, genStdTx)

		// the memo flag is used to store
		// the ip and node-id, for example this may be:
		// "528fd3df22b31f4969b05652bfe8f0fe921321d5@192.168.2.37:26656"
		nodeAddrIP := genStdTx.GetMemo()
		if len(nodeAddrIP) == 0 {
			return appGenTxs, persistentPeers, fmt.Errorf("failed to find node's address and IP in %s", fo.Name())
		}

		// genesis transactions must be single-message
		msgs := genStdTx.GetMsgs()

		err := stakingFuncManager.MsgDelegateWithBalance(msgs[1], balancesMap)

		if err != nil {
			return appGenTxs, persistentPeers, err
		}

		msgMoniker, err := stakingFuncManager.GetMsgCreateValidatorMoniker(msgs[0])

		if err != nil {
			return appGenTxs, persistentPeers, err
		}

		// exclude itself from persistent peers
		if msgMoniker != moniker {
			addressesIPs = append(addressesIPs, nodeAddrIP)
		}

	}

	sort.Strings(addressesIPs)
	persistentPeers = strings.Join(addressesIPs, ",")

	return appGenTxs, persistentPeers, nil
}
