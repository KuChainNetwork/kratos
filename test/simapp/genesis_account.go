package simapp

import (
	"errors"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/x/account"
	accountExported "github.com/KuChainNetwork/kuchain/x/account/exported"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
)

var _ accountExported.GenesisAccount = (*SimGenesisAccount)(nil)

// SimGenesisAccount defines a type that implements the GenesisAccount interface
// to be used for simulation accounts in the genesis state.
type SimGenesisAccount struct {
	*accountTypes.KuAccount

	// Assets
	Assets types.Coins `json:"assets" yaml:"assets"` // user coins For genesis

	// vesting account fields
	OriginalVesting  types.Coins `json:"original_vesting" yaml:"original_vesting"`   // total vesting coins upon initialization
	DelegatedFree    types.Coins `json:"delegated_free" yaml:"delegated_free"`       // delegated vested coins at time of delegation
	DelegatedVesting types.Coins `json:"delegated_vesting" yaml:"delegated_vesting"` // delegated vesting coins at time of delegation
	StartTime        int64       `json:"start_time" yaml:"start_time"`               // vesting start time (UNIX Epoch time)
	EndTime          int64       `json:"end_time" yaml:"end_time"`                   // vesting end time (UNIX Epoch time)

	// module account fields
	IsModule          bool     `json:"is_module" yaml:"is_module"`                   // is_module
	ModulePermissions []string `json:"module_permissions" yaml:"module_permissions"` // permissions of module account
}

func NewSimGenesisAccount(id types.AccountID, auth types.AccAddress) SimGenesisAccount {
	acc := account.NewKuAccount(id)
	if err := acc.SetAuth(auth); err != nil {
		panic(err)
	}
	return SimGenesisAccount{
		KuAccount: acc,
	}
}

func NewSimGenesisAccountByWallet(wallet *Wallet, id types.AccountID) SimGenesisAccount {
	addr := wallet.NewAccAddress()

	acc := account.NewKuAccount(id)
	if err := acc.SetAuth(addr); err != nil {
		panic(err)
	}

	return SimGenesisAccount{
		KuAccount: acc,
	}
}

// Validate checks for errors on the vesting and module account parameters
func (sga SimGenesisAccount) Validate() error {
	if !sga.OriginalVesting.IsZero() {
		if sga.OriginalVesting.IsAnyGT(sga.Assets) {
			return errors.New("vesting amount cannot be greater than total amount")
		}
		if sga.StartTime >= sga.EndTime {
			return errors.New("vesting start-time cannot be before end-time")
		}
	}

	if sga.IsModule {
		name, ok := sga.ID.ToName()
		if !ok {
			return errors.New("module account must be a name")
		}

		ma := supply.NewEmptyModuleAccount(name.String(), sga.ModulePermissions...)
		if err := ma.Validate(); err != nil {
			return err
		}
	}

	return sga.KuAccount.Validate()
}

func (sga SimGenesisAccount) WithAccountNumber(i uint64) SimGenesisAccount {
	sga.AccountNumber = i
	return sga
}

func (sga SimGenesisAccount) WithAsset(a types.Coins) SimGenesisAccount {
	sga.Assets = a
	return sga
}

func NewGenesisAccount(id types.AccountID, rootAuth types.AccAddress, num uint64) accountExported.GenesisAccount {
	acc := accountTypes.NewKuAccount(id)

	if err := acc.SetAuth(rootAuth); err != nil {
		panic(err)
	}

	if err := acc.SetAccountNumber(num); err != nil {
		panic(err)
	}

	return acc
}
