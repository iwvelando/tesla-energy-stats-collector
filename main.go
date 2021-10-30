package main

import (
	"flag"
	"fmt"
	"github.com/iwvelando/tesla-energy-stats-collector/config"
	"os"
	//        "github.com/iwvelando/tesla-energy-stats-collector/influxdb"
	"github.com/iwvelando/tesla-energy-stats-collector/connect"
	"go.uber.org/zap"
)

var BuildVersion = "UNKNOWN"

// CliInputs holds the data passed in via CLI parameters
type CliInputs struct {
	BuildVersion string
	Config       string
	ShowVersion  bool
}

func main() {

	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Println("{\"op\": \"main\", \"level\": \"fatal\", \"msg\": \"failed to initiate logger\"}")
		panic(err)
	}
	defer logger.Sync()

	cliInputs := CliInputs{
		BuildVersion: BuildVersion,
	}
	flags := flag.NewFlagSet("tesla-energy-stats-collector", 0)
	flags.StringVar(&cliInputs.Config, "config", "config.yaml", "Set the location for the YAML config file")
	flags.BoolVar(&cliInputs.ShowVersion, "version", false, "Print the version of modem-script")
	flags.Parse(os.Args[1:])

	if cliInputs.ShowVersion {
		fmt.Println(cliInputs.BuildVersion)
		os.Exit(0)
	}

	configuration, err := config.LoadConfiguration(cliInputs.Config)
	if err != nil {
		logger.Fatal("failed to parse configuration",
			zap.String("op", "main"),
			zap.Error(err),
		)
		panic(err)
	}

	client, refreshTime, err := connect.Auth(configuration)
	if err != nil {
		panic(err)
	}
	fmt.Println(refreshTime)

	err = connect.GetAll(configuration, client)
	if err != nil {
		panic(err)
	}

}
