package lcd

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/log"
	tmrpcserver "github.com/tendermint/tendermint/rpc/jsonrpc/server"

	kuLog "github.com/KuChainNetwork/kuchain/utils/log"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
)

// RestServer represents the Light Client Rest server
type RestServer struct {
	Mux    *mux.Router
	CliCtx context.CLIContext

	log      log.Logger
	listener net.Listener
}

// NewRestServer creates a new rest server instance
func NewRestServer(cdc *codec.Codec) *RestServer {
	r := mux.NewRouter()
	cliCtx := context.NewCLIContext().WithCodec(cdc)
	logger := kuLog.NewLoggerByZap(viper.GetBool(cli.TraceFlag), viper.GetString("log_level")).With("module", "rest-server")

	return &RestServer{
		Mux:    r,
		CliCtx: cliCtx,
		log:    logger,
	}
}

// Start starts the rest server
func (rs *RestServer) Start(listenAddr string, maxOpen int, readTimeout, writeTimeout uint, cors bool) (err error) {
	server.TrapSignal(func() {
		err := rs.listener.Close()
		rs.log.Error("error closing listener", "err", err)
	})

	cfg := tmrpcserver.DefaultConfig()
	cfg.MaxOpenConnections = maxOpen
	cfg.ReadTimeout = time.Duration(readTimeout) * time.Second
	cfg.WriteTimeout = time.Duration(writeTimeout) * time.Second

	rs.listener, err = tmrpcserver.Listen(listenAddr, cfg)
	if err != nil {
		return
	}
	rs.log.Info(
		fmt.Sprintf(
			"Starting application REST service (chain-id: %q)...",
			viper.GetString(flags.FlagChainID),
		),
	)

	var h http.Handler = rs.Mux
	if cors {
		allowAllCORS := handlers.CORS(handlers.AllowedHeaders([]string{"Content-Type"}))
		h = allowAllCORS(h)
	}

	return tmrpcserver.Serve(rs.listener, h, rs.log, cfg)
}

// ServeCommand will start the application REST service as a blocking process. It
// takes a codec to create a RestServer object and a function to register all
// necessary routes.
func ServeCommand(cdc *codec.Codec, registerRoutesFn func(*RestServer)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rest-server",
		Short: "Start LCD (light-client daemon), a local REST server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			rs := NewRestServer(cdc)

			registerRoutesFn(rs)

			// Start the rest server and return error if one exists
			err = rs.Start(
				viper.GetString(flags.FlagListenAddr),
				viper.GetInt(flags.FlagMaxOpenConnections),
				uint(viper.GetInt(flags.FlagRPCReadTimeout)),
				uint(viper.GetInt(flags.FlagRPCWriteTimeout)),
				viper.GetBool(flags.FlagUnsafeCORS),
			)

			return err
		},
	}

	return registerRestServerFlags(cmd)
}

// registerRestServerFlags registers the flags required for rest server
func registerRestServerFlags(cmd *cobra.Command) *cobra.Command {
	cmd = flags.GetCommands(cmd)[0]

	cmd.Flags().String(flags.FlagListenAddr, "tcp://localhost:1317", "The address for the server to listen on")
	cmd.Flags().Uint(flags.FlagMaxOpenConnections, 1000, "The number of maximum open connections")
	cmd.Flags().Uint(flags.FlagRPCReadTimeout, 10, "The RPC read timeout (in seconds)")
	cmd.Flags().Uint(flags.FlagRPCWriteTimeout, 10, "The RPC write timeout (in seconds)")
	cmd.Flags().Bool(flags.FlagUnsafeCORS, false, "Allows CORS requests from all domains. For development purposes only, use it at your own risk.")
	cmd.Flags().String("log_level", "*:info", "Log Level for rest server")

	return cmd
}
