package govcodec

import (
	"fmt"

	"github.com/KuChainNetwork/kuchain/x/gov/types"
	gov "github.com/KuChainNetwork/kuchain/x/gov/types"
	"github.com/KuChainNetwork/kuchain/x/params/types/proposal"
	"github.com/cosmos/cosmos-sdk/codec"
)

// CodecProposal defines the application-level concrete proposal type used in governance
// proposals.
type CodecProposal struct {
	types.ProposalBase `json:"base" yaml:"base"`
	Content            CodecContent `json:"content" yaml:"content"`
}

// Content defines the application-level allowed Content to be included in a
// governance proposal.
type CodecContent struct {
	Sum CodecContentInterface
}

func (m *CodecContent) GetText() *types.TextProposal {
	if x, ok := m.Sum.(*CodecContentText); ok {
		return x.Text
	}
	return nil
}

func (m *CodecContent) GetParameterChange() *proposal.ParameterChangeProposal {
	if x, ok := m.Sum.(*CodecContentParameterChange); ok {
		return x.ParameterChange
	}
	return nil
}

func (cc *CodecContent) GetContent() types.Content {
	if x := cc.GetText(); x != nil {
		return x
	}
	if x := cc.GetParameterChange(); x != nil {
		return x
	}
	return nil
}

func (cc *CodecContent) SetContent(value types.Content) error {
	if value == nil {
		cc.Sum = nil
		return nil
	}
	switch vt := value.(type) {
	case *types.TextProposal:
		cc.Sum = &CodecContentText{vt}
		return nil
	case types.TextProposal:
		cc.Sum = &CodecContentText{&vt}
		return nil
	case *proposal.ParameterChangeProposal:
		cc.Sum = &CodecContentParameterChange{vt}
		return nil
	case proposal.ParameterChangeProposal:
		cc.Sum = &CodecContentParameterChange{&vt}
		return nil
	}
	return fmt.Errorf("can't encode value of type %T as message CodecContent", value)
}

type CodecContentInterface interface {
}

type CodecContentText struct {
	Text *types.TextProposal `json:"text,omitempty"`
}
type CodecContentParameterChange struct {
	ParameterChange *proposal.ParameterChangeProposal `json:"parameter_change,omitempty"`
}

// MarshalProposal marshals a Proposal. It accepts a Proposal defined by the x/gov
// module and uses the application-level Proposal type which has the concrete
// Content implementation to serialize.
func MarshalProposal(cdc *codec.Codec, p gov.Proposal) ([]byte, error) {
	proposal := &CodecProposal{ProposalBase: p.ProposalBase}
	if err := proposal.Content.SetContent(p.Content); err != nil {
		return nil, err
	}

	return cdc.MarshalBinaryBare(proposal)
}

// UnmarshalProposal decodes a Proposal defined by the x/gov module and uses the
// application-level Proposal type which has the concrete Content implementation
// to deserialize.
func UnmarshalProposal(cdc *codec.Codec, bz []byte) (gov.Proposal, error) {
	proposal := &CodecProposal{}
	if err := cdc.UnmarshalBinaryBare(bz, proposal); err != nil {
		return gov.Proposal{}, err
	}

	return gov.Proposal{
		Content:      proposal.Content.GetContent(),
		ProposalBase: proposal.ProposalBase,
	}, nil
}
