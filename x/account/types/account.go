package types

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"gopkg.in/yaml.v2"
)

var _ exported.Account = (*KuAccount)(nil)

var (
	RootAuthName = types.MustName("root")
)

// KuAccount defines a account for kuchain
type KuAccount struct {
	ID            AccountID     `json:"id" yaml:"id"`
	AccountNumber uint64        `json:"account_number,omitempty" yaml:"account_number"`
	Auths         []AccountAuth `json:"auths" yaml:"auths"`
}

// NewKuAccount new KuAccount by name
func NewKuAccount(id types.AccountID) *KuAccount {
	acc := &KuAccount{
		ID: id,
	}

	return acc
}

// NewProtoKuAccount - a prototype function for KuAccount
func NewProtoKuAccount() exported.Account {
	return &KuAccount{}
}

// GetName - implements exported.Account
func (m KuAccount) GetName() types.Name {
	if n, ok := m.ID.ToName(); ok {
		return n
	}
	return types.Name{}
}

// SetName - implements exported.Account
func (m *KuAccount) SetName(n types.Name) error {
	m.ID = types.NewAccountIDFromName(n)
	return nil
}

// GetID - implements exported.Account
func (m KuAccount) GetID() types.AccountID {
	return m.ID
}

// SetID - implements exported.Account
func (m *KuAccount) SetID(id types.AccountID) {
	m.ID = id
}

// GetAuth - implements exported.Account
func (m KuAccount) GetAuth() types.AccAddress {
	// if KuAccount ID is just a account address, directly return
	if accAddress, ok := m.ID.ToAccAddress(); ok {
		return accAddress
	}

	if len(m.Auths) == 0 {
		panic(sdkerrors.Wrapf(types.ErrMissingAuth, "no auth for account %s", m.ID))
	}

	return m.Auths[0].Address
}

// SetAuth - implements exported.Account
func (m *KuAccount) SetAuth(auth types.AccAddress) {
	if len(m.Auths) == 0 {
		m.Auths = []AccountAuth{
			{
				Name:    RootAuthName,
				Address: auth,
			},
		}
	}

	m.Auths[0].Address = auth
}

// GetAccountNumber - implements exported.Account
func (m KuAccount) GetAccountNumber() uint64 {
	return m.AccountNumber
}

// SetAccountNumber - implements exported.Account
func (m *KuAccount) SetAccountNumber(n uint64) {
	m.AccountNumber = n
}

// Validate - implements exported.GenesisAccount
func (m KuAccount) Validate() error { return nil }

// String - implements exported.Account
func (m KuAccount) String() string {
	out, _ := m.MarshalYAML()
	return out.(string)
}

func (m *KuAccount) makeStr(marFunc func(in interface{}) (out []byte, err error)) ([]byte, error) {
	alias := struct {
		ID            string        `json:"id" yaml:"id"`
		Auths         []AccountAuth `json:"auths" yaml:"auths"`
		AccountNumber uint64        `json:"account_number" yaml:"account_number"`
	}{
		ID:            m.ID.String(),
		Auths:         m.Auths,
		AccountNumber: m.AccountNumber,
	}

	return marFunc(alias)
}

// MarshalYAML returns the YAML representation of an account.
func (m KuAccount) MarshalYAML() (interface{}, error) {
	bz, err := m.makeStr(yaml.Marshal)
	return string(bz), err
}

type AccountAuth struct {
	Name    types.Name     `json:"name" yaml:"name"`
	Address sdk.AccAddress `json:"address,omitempty" yaml:"address"`
}

// String - implements exported.Account
func (m AccountAuth) String() string {
	out, _ := m.MarshalYAML()
	return out.(string)
}

func (m *AccountAuth) makeStr(marFunc func(in interface{}) (out []byte, err error)) ([]byte, error) {
	alias := struct {
		Name string `json:"name" yaml:"name"`
		Auth string `json:"auth" yaml:"auth"`
	}{
		Name: m.Name.String(),
		Auth: m.Address.String(),
	}

	return marFunc(alias)
}

// MarshalYAML returns the YAML representation of an account.
func (m AccountAuth) MarshalYAML() (interface{}, error) {
	bz, err := m.makeStr(yaml.Marshal)
	return string(bz), err
}
