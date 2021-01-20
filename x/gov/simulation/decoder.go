package simulation

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/KuChainNetwork/kuchain/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

// DecodeStore unmarshals the KVPair's Value to the corresponding gov type
func DecodeStore(cdc *codec.Codec, kvA, kvB kv.Pair) string {
	switch {
	case bytes.Equal(kvA.Key[:1], types.ProposalsKeyPrefix):
		var proposalA, proposalB types.Proposal
		cdc.MustUnmarshalBinaryBare(kvA.Value, &proposalA)
		cdc.MustUnmarshalBinaryBare(kvB.Value, &proposalB)
		return fmt.Sprintf("%v\n%v", proposalA, proposalB)

	case bytes.Equal(kvA.Key[:1], types.ActiveProposalQueuePrefix),
		bytes.Equal(kvA.Key[:1], types.InactiveProposalQueuePrefix),
		bytes.Equal(kvA.Key[:1], types.ProposalIDKey):
		proposalIDA := binary.LittleEndian.Uint64(kvA.Value)
		proposalIDB := binary.LittleEndian.Uint64(kvB.Value)
		return fmt.Sprintf("proposalIDA: %d\nProposalIDB: %d", proposalIDA, proposalIDB)

	case bytes.Equal(kvA.Key[:1], types.DepositsKeyPrefix):
		var depositA, depositB types.Deposit
		cdc.MustUnmarshalBinaryBare(kvA.Value, &depositA)
		cdc.MustUnmarshalBinaryBare(kvB.Value, &depositB)
		return fmt.Sprintf("%v\n%v", depositA, depositB)

	case bytes.Equal(kvA.Key[:1], types.VotesKeyPrefix):
		var voteA, voteB types.Vote
		cdc.MustUnmarshalBinaryBare(kvA.Value, &voteA)
		cdc.MustUnmarshalBinaryBare(kvB.Value, &voteB)
		return fmt.Sprintf("%v\n%v", voteA, voteB)

	default:
		panic(fmt.Sprintf("invalid governance key prefix %X", kvA.Key[:1]))
	}
}
