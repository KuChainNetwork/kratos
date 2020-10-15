package block

import (
	"encoding/json"

	"github.com/KuChainNetwork/kuchain/chain/client/utils"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/pkg/errors"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
	tmliteProxy "github.com/tendermint/tendermint/lite/proxy"
	sdk "github.com/tendermint/tendermint/types"
)

type DecodeData struct {
	Txs     []json.RawMessage  `json:"txs"`
	TxsHash []tmbytes.HexBytes `json:"txs_hash"`

	// Volatile
	hash tmbytes.HexBytes
}

type DecodeBlock struct {
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

	txs := res.Block.Txs
	stdTxs := make([]json.RawMessage, 0, len(txs))
	txsHash := make([]tmbytes.HexBytes, 0, len(txs))
	for _, tx := range txs {
		var stdTx types.StdTx
		err = cliCtx.Codec.UnmarshalBinaryLengthPrefixed(tx, &stdTx)
		if err != nil {
			return nil, errors.Wrapf(err, "unmarshal stdtx error")
		}

		datas, err := utils.OutputPrettifyJSON(cliCtx, stdTx)
		if err != nil {
			return nil, errors.Wrapf(err, "stdTx to json error")
		}

		stdTxs = append(stdTxs, datas)
		txsHash = append(txsHash, tx.Hash())
	}

	decodeRes := &DecodeResultBlock{
		BlockID: res.BlockID,
		DecodeBlock: &DecodeBlock{
			Header:     res.Block.Header,
			Evidence:   res.Block.Evidence,
			LastCommit: res.Block.LastCommit,
			DecodeData: DecodeData{
				Txs:     stdTxs,
				TxsHash: txsHash,
			}}}

	return utils.OutputPrettifyJSON(cliCtx, decodeRes)
}
