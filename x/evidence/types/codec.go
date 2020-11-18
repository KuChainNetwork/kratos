package types

import (
	"github.com/KuChainNetwork/kuchain/x/evidence/exported"
	"github.com/cosmos/cosmos-sdk/codec"
)

// Codec defines the interface required to serialize evidence
type Codec interface {
	MarshalEvidence(exported.Evidence) ([]byte, error)
	UnmarshalEvidence([]byte) (exported.Evidence, error)
	MarshalEvidenceJSON(exported.Evidence) ([]byte, error)
	UnmarshalEvidenceJSON([]byte) (exported.Evidence, error)
}

type EvidenceCodec struct {
	amino *codec.Codec
}

func NewEveidenceCodec(amino *codec.Codec) *EvidenceCodec {
	return &EvidenceCodec{amino: amino}
}

// MarshalEvidence marshals an Evidence interface. If the given type implements
// the Marshaler interface, it is treated as a Proto-defined message and
// serialized that way. Otherwise, it falls back on the internal Amino codec.
func (c *EvidenceCodec) MarshalEvidence(evidenceI exported.Evidence) ([]byte, error) {
	evidence := &Evidence{}
	if err := evidence.SetEvidence(evidenceI); err != nil {
		return nil, err
	}

	return c.amino.MarshalBinaryBare(evidence)
}

// UnmarshalEvidence returns an Evidence interface from raw encoded evidence
// bytes of a Proto-based Evidence type. An error is returned upon decoding
// failure.
func (c *EvidenceCodec) UnmarshalEvidence(bz []byte) (exported.Evidence, error) {
	evidence := &Evidence{}
	if err := c.amino.UnmarshalBinaryBare(bz, evidence); err != nil {
		return nil, err
	}

	return evidence.GetEvidence(), nil
}

// MarshalEvidenceJSON JSON encodes an evidence object implementing the Evidence
// interface.
func (c *EvidenceCodec) MarshalEvidenceJSON(evidence exported.Evidence) ([]byte, error) {
	return c.amino.MarshalJSON(evidence)
}

// UnmarshalEvidenceJSON returns an Evidence from JSON encoded bytes
func (c *EvidenceCodec) UnmarshalEvidenceJSON(bz []byte) (exported.Evidence, error) {
	evidence := &Evidence{}
	if err := c.amino.UnmarshalJSON(bz, evidence); err != nil {
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
	ModuleCdc = codec.New()

	// EvidenceCdc references the global x/evidence module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/evidence and
	// defined at the application level.
	EvidenceCdc = NewEveidenceCodec(ModuleCdc)
)

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
