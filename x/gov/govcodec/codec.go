package govcodec

import (
	gov "github.com/KuChain-io/kuchain/x/gov/types"
	"github.com/cosmos/cosmos-sdk/codec"
)

var (
	_ gov.Codec = (*Codec)(nil)
)

type Codec struct {
	codec.Marshaler

	// Keep reference to the amino codec to allow backwards compatibility along
	// with type, and interface registration.
	amino *codec.Codec
}

func NewAppCodec(amino *codec.Codec) *Codec {
	return &Codec{Marshaler: codec.NewHybridCodec(amino), amino: amino}
}

// MarshalProposal marshals a Proposal. It accepts a Proposal defined by the x/gov
// module and uses the application-level Proposal type which has the concrete
// Content implementation to serialize.
func (c *Codec) MarshalProposal(p gov.Proposal) ([]byte, error) {
	proposal := &CodecProposal{ProposalBase: p.ProposalBase}
	if err := proposal.Content.SetContent(p.Content); err != nil {
		return nil, err
	}

	return c.Marshaler.MarshalBinaryBare(proposal)
}

// UnmarshalProposal decodes a Proposal defined by the x/gov module and uses the
// application-level Proposal type which has the concrete Content implementation
// to deserialize.
func (c *Codec) UnmarshalProposal(bz []byte) (gov.Proposal, error) {
	proposal := &CodecProposal{}
	if err := c.Marshaler.UnmarshalBinaryBare(bz, proposal); err != nil {
		return gov.Proposal{}, err
	}

	return gov.Proposal{
		Content:      proposal.Content.GetContent(),
		ProposalBase: proposal.ProposalBase,
	}, nil
}

func NewGovCodec(amino *codec.Codec) *Codec {
	return &Codec{Marshaler: codec.NewHybridCodec(amino), amino: amino}
}

var (
	Gov_Cdc = NewGovCodec(gov.Amino)
)
