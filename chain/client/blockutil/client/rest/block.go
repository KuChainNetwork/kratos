package rest

import (
	"sync"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	sdk "github.com/tendermint/tendermint/types"
)

type DecodeData struct {
	//Txs []string `json:"txs"`
	Txs []types.StdTx `json:"txs"`

	// Volatile
	hash tmbytes.HexBytes
}

type DecodeBlock struct {
	mtx        sync.Mutex
	sdk.Header `json:"header"`
	DecodeData `json:"data"`
	Evidence   sdk.EvidenceData `json:"evidence"`
	LastCommit *sdk.Commit      `json:"last_commit"`
}

type DecodeResultBlock struct {
	BlockID     sdk.BlockID  `json:"block_id"`
	DecodeBlock *DecodeBlock `json:"block"`
}

func getBlock(cliCtx context.CLIContext, height *int64) ([]byte, error) {
	// get the node
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	// header -> BlockchainInfo
	// header, tx -> Block
	// results -> BlockResults
	res, err := node.Block(height)
	if err != nil {
		return nil, err
	}

	if !cliCtx.TrustNode {
		check, err := cliCtx.Verify(res.Block.Height)
		if err != nil {
			return nil, err
		}

		if err := tmliteProxy.ValidateHeader(&res.Block.Header, check); err != nil {
			return nil, err
		}

		if err = tmliteProxy.ValidateBlock(res.Block, check); err != nil {
			return nil, err
		}
	}

	var stdTxs []types.StdTx
	txs := res.Block.Txs
	for _, tx := range txs {
		var stdTx types.StdTx
		err = cliCtx.Codec.UnmarshalBinaryBare(tx, &stdTx)
		if err != nil {
			return nil, err
		}
		stdTxs = append(stdTxs, stdTx)
	}

	decodeRes := &DecodeResultBlock{BlockID: res.BlockID, DecodeBlock: &DecodeBlock{Header: res.Block.Header,
		Evidence: res.Block.Evidence, LastCommit: res.Block.LastCommit, DecodeData: DecodeData{Txs: stdTxs}}}

	if cliCtx.Indent {
		return cliCtx.Codec.MarshalJSONIndent(decodeRes, "", "  ")
	}

	return cliCtx.Codec.MarshalJSON(decodeRes)
}
