package main

import (
	"flag"
	"fmt"
	"github.com/iwvelando/tesla-energy-stats-collector/config"
	"github.com/iwvelando/tesla-energy-stats-collector/connect"
	"github.com/iwvelando/tesla-energy-stats-collector/influxdb"
	"go.uber.org/zap"
	"os"
	"time"
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
			zap.String("op", "config.LoadConfiguration"),
			zap.Error(err),
		)
	}

	tesla, refreshTime, err := connect.Auth(configuration)
	if err != nil {
		logger.Fatal("failed to authenticate to Tesla energy gateway",
			zap.String("op", "connect.Auth"),
			zap.Error(err),
		)
	}

	influxClient, writeAPI, err := influxdb.Connect(configuration)
	if err != nil {
		logger.Fatal("failed to authenticate to InfluxDB",
			zap.String("op", "influxdb.Connect"),
			zap.Error(err),
		)
	}
	defer influxClient.Close()

	for {

		if time.Now().After(refreshTime) {
			tesla, refreshTime, err = connect.Auth(configuration)
			if err != nil {
				logger.Fatal("failed to refresh authentication to Tesla energy gateway",
					zap.String("op", "connect.Auth"),
					zap.Error(err),
				)
			}
		}

		pollStartTime := time.Now()

		metrics, err := connect.GetAll(configuration, tesla)
		if err != nil {
			logger.Error("failed to query all metrics, waiting for next poll",
				zap.String("op", "connect.GetAll"),
				zap.Error(err),
			)
		} else {
			influxdb.WriteAll(configuration, writeAPI, metrics)
		}

		timeRemaining := configuration.Polling.Interval*time.Second - time.Since(pollStartTime)
		time.Sleep(time.Duration(timeRemaining))
		continue

	}

}
