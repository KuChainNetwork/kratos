package types

import (
	"github.com/KuChainNetwork/kuchain/chain/types"
	accountExported "github.com/KuChainNetwork/kuchain/x/account/exported"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
)

type (
	Account   = accountExported.Account
	AccountID = types.AccountID
)

var (
	MustName             = types.MustName
	NewAccountIDFromName = types.NewAccountIDFromName
)

type (
	ModuleAccount         = accountTypes.ModuleAccount
	PermissionsForAddress = accountTypes.PermissionsForAddress
)

var (
	NewEmptyModuleAccount    = accountTypes.NewEmptyModuleAccount
	NewPermissionsForAddress = accountTypes.NewPermissionsForAddress
)

const (
	Minter  = accountTypes.Minter
	Burner  = accountTypes.Burner
	Staking = accountTypes.Staking
)
