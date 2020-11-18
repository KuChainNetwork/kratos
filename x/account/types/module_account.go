package types

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account/exported"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"gopkg.in/yaml.v2"
)

var (
	_ exported.GenesisAccount = (*ModuleAccount)(nil)
	_ exported.Account        = (*ModuleAccount)(nil)
)

// ModuleAccount defines an account for modules that holds coins on a pool
type ModuleAccount struct {
	KuAccount
	Permissions []string `json:"permissions,omitempty"`
}

// NewModuleAddress creates an AccAddress from the hash of the module's name
func NewModuleAddress(name string) sdk.AccAddress {
	return sdk.AccAddress(crypto.AddressHash([]byte(name)))
}

// NewEmptyModuleAccount creates a empty ModuleAccount from a string
func NewEmptyModuleAccount(name string, permissions ...string) *ModuleAccount {
	n, err := types.NewName(name)
	if err != nil {
		panic(err)
	}

	baseAccount := NewKuAccount(types.NewAccountIDFromName(n))

	if err := validatePermissions(permissions...); err != nil {
		panic(err)
	}

	return &ModuleAccount{
		KuAccount:   *baseAccount,
		Permissions: permissions,
	}
}

// NewModuleAccount creates a new ModuleAccount instance
func NewModuleAccount(ba KuAccount, name string, permissions ...string) *ModuleAccount {
	if err := validatePermissions(permissions...); err != nil {
		panic(err)
	}

	return &ModuleAccount{
		KuAccount:   ba,
		Permissions: permissions,
	}
}

// HasPermission returns whether or not the module account has permission.
func (ma ModuleAccount) HasPermission(permission string) bool {
	for _, perm := range ma.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// GetPermissions returns permissions granted to the module account
func (ma ModuleAccount) GetPermissions() []string {
	return ma.Permissions
}

// GetAddress for some old API calls
func (ma ModuleAccount) GetAddress() sdk.AccAddress {
	return ma.GetAuth()
}

// SetPubKey - Implements Account
func (ma ModuleAccount) SetPubKey(pubKey crypto.PubKey) error {
	return fmt.Errorf("not supported for module accounts")
}

// SetSequence - Implements Account
func (ma ModuleAccount) SetSequence(seq uint64) error {
	return fmt.Errorf("not supported for module accounts")
}

// Validate checks for errors on the account fields
func (ma ModuleAccount) Validate() error {
	if ma.GetName().Empty() {
		return errors.New("module account name cannot be blank")
	}

	return ma.KuAccount.Validate()
}

type moduleAccountPretty struct {
	Address       sdk.AccAddress `json:"address" yaml:"address"`
	PubKey        string         `json:"public_key" yaml:"public_key"`
	AccountNumber uint64         `json:"account_number" yaml:"account_number"`
	Name          string         `json:"name" yaml:"name"`
	Permissions   []string       `json:"permissions" yaml:"permissions"`
}

func (ma ModuleAccount) String() string {
	out, _ := ma.MarshalYAML()
	return out.(string)
}

// MarshalYAML returns the YAML representation of a ModuleAccount.
func (ma ModuleAccount) MarshalYAML() (interface{}, error) {
	add := types.AccAddress{}
	if len(ma.Auths) > 0 {
		add = ma.GetAuth()
	}

	bs, err := yaml.Marshal(moduleAccountPretty{
		Address:       add,
		PubKey:        "",
		AccountNumber: ma.AccountNumber,
		Name:          ma.ID.String(),
		Permissions:   ma.Permissions,
	})

	if err != nil {
		return nil, err
	}

	return string(bs), nil
}

// MarshalJSON returns the JSON representation of a ModuleAccount.
func (ma ModuleAccount) MarshalJSON() ([]byte, error) {
	add := types.AccAddress{}
	if len(ma.Auths) > 0 {
		add = ma.GetAuth()
	}
	return json.Marshal(moduleAccountPretty{
		Address:       add,
		PubKey:        "",
		AccountNumber: ma.AccountNumber,
		Name:          ma.ID.String(),
		Permissions:   ma.Permissions,
	})
}

// UnmarshalJSON unmarshal raw JSON bytes into a ModuleAccount.
func (ma *ModuleAccount) UnmarshalJSON(bz []byte) error {
	var alias moduleAccountPretty
	if err := json.Unmarshal(bz, &alias); err != nil {
		return err
	}

	name, err := types.NewName(alias.Name)
	if err != nil {
		return err
	}

	newAccount := NewKuAccount(NewAccountIDFromName(name))
	if err := newAccount.SetAuth(alias.Address); err != nil {
		panic(err)
	}

	if err := newAccount.SetAccountNumber(alias.AccountNumber); err != nil {
		panic(err)
	}

	ma.KuAccount = *newAccount
	ma.Permissions = alias.Permissions

	return nil
}
