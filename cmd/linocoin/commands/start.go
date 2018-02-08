package commands

import (
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/lino-network/lino/app"

	"github.com/tendermint/abci/server"
	eyesApp "github.com/tendermint/merkleeyes/app"
	eyes "github.com/tendermint/merkleeyes/client"
	"github.com/tendermint/tmlibs/cli"
	cmn "github.com/tendermint/tmlibs/common"

	tcmd "github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/proxy"
	"github.com/tendermint/tendermint/types"

)

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start linocoin",
	RunE:  startCmd,
}

// TODO: move to config file
const EyesCacheSize = 10000

//nolint
const (
	FlagAddress           = "address"
	FlagEyes              = "eyes"
	FlagWithoutTendermint = "without-tendermint"
)

func init() {
	flags := StartCmd.Flags()
	flags.String(FlagAddress, "tcp://0.0.0.0:46658", "Listen address")
	flags.String(FlagEyes, "local", "MerkleEyes address, or 'local' for embedded")
	flags.Bool(FlagWithoutTendermint, false, "Only run linocoin abci app, assume external tendermint process")
	// add all standard 'tendermint node' flags
	tcmd.AddNodeFlags(StartCmd)
}

func startCmd(cmd *cobra.Command, args []string) error {
	rootDir := viper.GetString(cli.HomeFlag)
	meyes := viper.GetString(FlagEyes)

	// Connect to MerkleEyes
	var eyesCli *eyes.Client
	if meyes == "local" {
		eyesApp.SetLogger(logger.With("module", "merkleeyes"))
		eyesCli = eyes.NewLocalClient(path.Join(rootDir, "data", "merkleeyes.db"), EyesCacheSize)
	} else {
		var err error
		eyesCli, err = eyes.NewClient(meyes)
		if err != nil {
			return errors.Errorf("Error connecting to MerkleEyes: %v\n", err)
		}
	}

	// Create Linocoin app
	linocoinApp := app.NewLinocoin(eyesCli)
	linocoinApp.SetLogger(logger.With("module", "app"))

	// register IBC plugn
	//linocoinApp.RegisterPlugin(NewIBCPlugin())

	// register all other plugins
	for _, p := range plugins {
		linocoinApp.RegisterPlugin(p.newPlugin())
	}

	// if chain_id has not been set yet, load the genesis.
	// else, assume it's been loaded
	if linocoinApp.GetState().GetChainID() == "" {
		// If genesis file exists, set key-value options
		genesisFile := path.Join(rootDir, "genesis.json")
		if _, err := os.Stat(genesisFile); err == nil {
			err := linocoinApp.LoadGenesis(genesisFile)
			if err != nil {
				return errors.Errorf("Error in LoadGenesis: %v\n", err)
			}
		} else {
			fmt.Printf("No genesis file at %s, skipping...\n", genesisFile)
		}
	}

	chainID := linocoinApp.GetState().GetChainID()
	if viper.GetBool(FlagWithoutTendermint) {
		logger.Info("Starting Linocoin without Tendermint", "chain_id", chainID)
		// run just the abci app/server
		return startLinocoinABCI(linocoinApp)
	} else {
		logger.Info("Starting Linocoin with Tendermint", "chain_id", chainID)
		// start the app with tendermint in-process
		return startTendermint(rootDir, linocoinApp)
	}
}

func startLinocoinABCI(linocoinApp *app.Linocoin) error {
	// Start the ABCI listener
	addr := viper.GetString(FlagAddress)
	svr, err := server.NewServer(addr, "socket", linocoinApp)
	if err != nil {
		return errors.Errorf("Error creating listener: %v\n", err)
	}
	svr.SetLogger(logger.With("module", "abci-server"))
	svr.Start()

	// Wait forever
	cmn.TrapSignal(func() {
		// Cleanup
		svr.Stop()
	})
	return nil
}

func startTendermint(dir string, linocoinApp *app.Linocoin) error {
	cfg, err := tcmd.ParseConfig()
	if err != nil {
		return err
	}

	// Create & start tendermint node
	privValidator := types.LoadOrGenPrivValidator(cfg.PrivValidatorFile(), logger)
	n := node.NewNode(cfg, privValidator, proxy.NewLocalClientCreator(linocoinApp), logger.With("module", "node"))

	_, err = n.Start()
	if err != nil {
		return err
	}

	// Trap signal, run forever.
	n.RunForever()
	return nil
}
