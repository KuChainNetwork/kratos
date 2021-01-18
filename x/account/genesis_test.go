package account_test

import (
	"math/rand"
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/constants"
	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/account"
	accountTypes "github.com/KuChainNetwork/kuchain/x/account/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestGenesisExport(t *testing.T) {
	asset1 := types.NewInt64Coins(constants.DefaultBondDenom, 10000000000)
	genAcc := simapp.NewSimGenesisAccount(account1, addr1).WithAsset(asset1)
	genAccs := simapp.NewGenesisAccounts(wallet.GetRootAuth(), genAcc)
	app := simapp.SetupWithGenesisAccounts(genAccs).WithWallet(wallet)

	// create some accounts
	type tmpAccouts struct {
		ID   types.AccountID
		Auth types.AccAddress
	}

	Convey("test genesis export", t, func() {
		r := rand.New(rand.NewSource(app.RandSeed()))

		auths := make([]types.AccAddress, 0, 1024)
		accountCreated := make([]tmpAccouts, 0)

		const (
			createAccountNum = 1000
			createAuths      = 1000
		)

		// close log
		Convey("create accounts", func() {
			// create some auths
			for i := 0; i < createAccountNum; i++ {
				if r.Intn(10) > 6 {
					simapp.AfterBlockCommitted(app, 1)
				}

				id := simapp.RandCommonName(app)
				auth := app.GetWallet().NewAccAddressByName(id.MustName())

				accountCreated = append(accountCreated, tmpAccouts{
					ID:   id,
					Auth: auth,
				})

				err := testAccountCreate(t, app, app.GetWallet(), true, constants.SystemAccountID, id.MustName(), auth)
				if err != nil {
					panic(err)
				}
			}
		})

		// sign some auths
		Convey("create auths", func() {
			// create some auths
			for i := 0; i < createAuths; i++ {
				if r.Intn(10) > 6 {
					simapp.AfterBlockCommitted(app, 1)
				}

				auth := app.GetWallet().NewAccAddress()
				auths = append(auths, auth)

				err := simapp.CommitTransferTx(t, app, wallet,
					true,
					constants.SystemAccountID, types.NewAccountIDFromAccAdd(auth),
					types.NewInt64CoreCoins(1), constants.SystemAccountID)

				if err != nil {
					panic(err)
				}
			}
		})

		// export account genesis
		genesis := account.ExportGenesis(app.NewTestContext(), *app.AccountKeeper())
		So(accountTypes.ValidateGenesis(genesis), ShouldBeNil)
	})
}
