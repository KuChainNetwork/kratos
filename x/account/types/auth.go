package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"gopkg.in/yaml.v2"
)

type Auth struct {
	Name      Name           `json:"name" yaml:"name"`
	Address   sdk.AccAddress `json:"address,omitempty" yaml:"address"`
	PublicKey []byte         `json:"public_key,omitempty" yaml:"public_key"`
	Number    uint64         `json:"number,omitempty" yaml:"number"`
	Sequence  uint64         `json:"sequence,omitempty" yaml:"sequence"`
}

func NewAuth(address sdk.AccAddress) Auth {
	return Auth{
		Address:  address,
		Sequence: 1,
	}
}

// GetPubKey - Implements sdk.Account.
func (m Auth) GetPubKey() (pk crypto.PubKey) {
	if len(m.PublicKey) == 0 {
		return nil
	}

	codec.Cdc.MustUnmarshalBinaryBare(m.PublicKey, &pk)
	return pk
}

// SetPubKey - Implements sdk.Account.
func (m *Auth) SetPubKey(pubKey crypto.PubKey) {
	if pubKey == nil {
		m.PublicKey = nil
	} else {
		m.PublicKey = pubKey.Bytes()
	}
}

// GetSequence get sequence number for auth data
func (m Auth) GetSequence() uint64 { return m.Sequence }

// GetNumber
func (m Auth) GetNumber() uint64 { return m.Number }

// GetAddress
func (m Auth) GetAddress() sdk.AccAddress { return m.Address }

// SetSequence - implements exported.Account
func (m *Auth) SetSequence(s uint64) {
	m.Sequence = s
}

// SetPubKey - Implements sdk.Account.
func (m *Auth) SetAccountNum(num uint64) {
	m.Number = num
}

// SetAddress
func (m *Auth) SetAddress(address sdk.AccAddress) {
	m.Address = address
}

// String - implements exported.Account
func (m Auth) String() string {
	out, _ := m.MarshalYAML()
	return out.(string)
}

func (m *Auth) makeAuthStrByMarshal(marFunc func(in interface{}) (out []byte, err error)) ([]byte, error) {
	alias := struct {
		Name     string `json:"name" yaml:"name"`
		Auth     string `json:"auth" yaml:"auth"`
		PubKey   string `json:"publicKey" yaml:"publicKey"`
		Number   uint64 `json:"number" yaml:"number"`
		Sequence uint64 `json:"sequence" yaml:"sequence"`
	}{
		Name:     m.Name.String(),
		Auth:     m.Address.String(),
		Number:   m.GetNumber(),
		Sequence: m.GetSequence(),
	}

	pk := m.GetPubKey()

	if pk != nil {
		alias.PubKey = pk.Address().String()
	}

	return marFunc(alias)
}

// MarshalYAML returns the YAML representation of an account.
func (m Auth) MarshalYAML() (interface{}, error) {
	bz, err := m.makeAuthStrByMarshal(yaml.Marshal)
	return string(bz), err
}
