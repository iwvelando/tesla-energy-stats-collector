package main

import (
	"flag"
	"fmt"
	"github.com/iwvelando/tesla-energy-stats-collector/config"
	"os"
	//        "github.com/iwvelando/tesla-energy-stats-collector/influxdb"
	"github.com/iwvelando/tesla-energy-stats-collector/connect"
	"github.com/iwvelando/tesla-energy-stats-collector/model"
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

	client, err := connect.Auth(configuration)
	if err != nil {
		panic(err)
	}

	tegStatus := model.TegStatus{}
	err = connect.GetEndpoint(configuration, client, "/api/status", &tegStatus)
	if err != nil {
		panic(err)
	}
	fmt.Println(tegStatus)

	tegMeters := model.TegMeters{}
	err = connect.GetEndpoint(configuration, client, "/api/meters/aggregates", &tegMeters)
	if err != nil {
		panic(err)
	}
	fmt.Println(tegMeters)

	tegOperation := model.TegOperation{}
	err = connect.GetEndpoint(configuration, client, "/api/operation", &tegOperation)
	if err != nil {
		panic(err)
	}
	fmt.Println(tegOperation)

	tegPowerwalls := model.TegPowerwalls{}
	err = connect.GetEndpoint(configuration, client, "/api/powerwalls", &tegPowerwalls)
	if err != nil {
		panic(err)
	}
	fmt.Println(tegPowerwalls)

	tegSiteInfo := model.TegSiteInfo{}
	err = connect.GetEndpoint(configuration, client, "/api/site_info", &tegSiteInfo)
	if err != nil {
		panic(err)
	}
	fmt.Println(tegSiteInfo)

	tegSolars := []model.TegSolars{}
	err = connect.GetEndpoint(configuration, client, "/api/solars", &tegSolars)
	if err != nil {
		panic(err)
	}
	for _, i := range tegSolars {
		fmt.Println(i)
	}

}
