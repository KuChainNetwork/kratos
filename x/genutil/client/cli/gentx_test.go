package cli

import (
	"testing"

	"github.com/KuChainNetwork/kuchain/chain/config"
	"github.com/KuChainNetwork/kuchain/test/simapp"
	"github.com/KuChainNetwork/kuchain/x/asset"
	"github.com/KuChainNetwork/kuchain/x/staking"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/server"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
)

func TestGenTxCmdCmd(t *testing.T) {
	Convey("TestGenTxCmdCmd", t, func() {
		cmd := &cobra.Command{}

		wallet := simapp.NewWallet()
		key := wallet.NewAccAddress()

		SetupViper()
		keysHome, cleanup := simapp.NewTestCaseDir(t)
		defer cleanup()
		viper.Set(flags.FlagHome, keysHome)

		addKeyCmd := keys.AddKeyCommand()
		err := addKeyCmd.RunE(cmd, []string{"alice"})
		So(err, ShouldBeNil)

		home, cleanup := simapp.NewTestCaseDir(t)
		defer cleanup()
		viper.Set(flags.FlagHome, home)
		viper.Set(flagClientHome, keysHome)

		logger := log.NewNopLogger()
		cfg, err := tcmd.ParseConfig()
		So(err, ShouldBeNil)
		ctx := server.NewContext(cfg, logger)
		cdc := makeCodec()

		initCmd := InitCmd(ctx, cdc, simapp.ModuleBasics, home)
		err = initCmd.RunE(nil, []string{"moniker"})
		So(err, ShouldBeNil)

		gentxCmd := GenTxCmd(ctx, cdc, simapp.ModuleBasics, staking.AppModuleBasic{}, asset.GenesisBalancesIterator{},
			home, keysHome, staking.NewFuncManager())
		err = gentxCmd.RunE(cmd, []string{"moniker", key.String()})
		So(err, ShouldBeNil)
	})
}

func SetupViper() {
	viper.Set(flags.FlagName, "alice")
	viper.Set(flags.FlagKeyringBackend, "test")
	viper.Set(cli.OutputFlag, "json")
	viper.Set(keys.FlagBechPrefix, config.PrefixAccount)
}
