package types

import (
	"github.com/KuChainNetwork/kuchain/x/evidence/exported"
	"github.com/cosmos/cosmos-sdk/codec"
)

// Codec defines the interface required to serialize evidence
type Codec interface {
	codec.Marshaler

	MarshalEvidence(exported.Evidence) ([]byte, error)
	UnmarshalEvidence([]byte) (exported.Evidence, error)
	MarshalEvidenceJSON(exported.Evidence) ([]byte, error)
	UnmarshalEvidenceJSON([]byte) (exported.Evidence, error)
}

type EvidenceCdc struct {
	codec.Marshaler

	// Keep reference to the amino codec to allow backwards compatibility along
	// with type, and interface registration.
	amino *codec.Codec
}

func NewEveidenceCodec(amino *codec.Codec) *EvidenceCdc {
	return &EvidenceCdc{Marshaler: codec.NewHybridCodec(amino), amino: amino}
}

// MarshalEvidence marshals an Evidence interface. If the given type implements
// the Marshaler interface, it is treated as a Proto-defined message and
// serialized that way. Otherwise, it falls back on the internal Amino codec.
func (c *EvidenceCdc) MarshalEvidence(evidenceI exported.Evidence) ([]byte, error) {
	evidence := &Evidence{}
	if err := evidence.SetEvidence(evidenceI); err != nil {
		return nil, err
	}

	return c.Marshaler.MarshalBinaryBare(evidence)
}

// UnmarshalEvidence returns an Evidence interface from raw encoded evidence
// bytes of a Proto-based Evidence type. An error is returned upon decoding
// failure.
func (c *EvidenceCdc) UnmarshalEvidence(bz []byte) (exported.Evidence, error) {
	evidence := &Evidence{}
	if err := c.Marshaler.UnmarshalBinaryBare(bz, evidence); err != nil {
		return nil, err
	}

	return evidence.GetEvidence(), nil
}

// MarshalEvidenceJSON JSON encodes an evidence object implementing the Evidence
// interface.
func (c *EvidenceCdc) MarshalEvidenceJSON(evidence exported.Evidence) ([]byte, error) {
	return c.Marshaler.MarshalJSON(evidence)
}

// UnmarshalEvidenceJSON returns an Evidence from JSON encoded bytes
func (c *EvidenceCdc) UnmarshalEvidenceJSON(bz []byte) (exported.Evidence, error) {
	evidence := &Evidence{}
	if err := c.Marshaler.UnmarshalJSON(bz, evidence); err != nil {
		return nil, err
	}

	return evidence.GetEvidence(), nil
}

// RegisterCodec registers all the necessary types and interfaces for the
// evidence module.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterInterface((*exported.Evidence)(nil), nil)
	cdc.RegisterConcrete(MsgSubmitEvidenceBase{}, "kuchain/MsgSubmitEvidenceBase", nil)
	cdc.RegisterConcrete(Equivocation{}, "kuchain/Equivocation", nil)
}

var (
	amino = codec.New()

	// ModuleCdc references the global x/evidence module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/evidence and
	// defined at the application level.
	Evidence_Cdc = NewEveidenceCodec(amino)
)

func init() {
	RegisterCodec(amino)
	codec.RegisterCrypto(amino)
	amino.Seal()
}
