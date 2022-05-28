package main

import (
	"flag"
	"fmt"
	"github.com/iwvelando/tesla-energy-stats-collector/config"
	"github.com/iwvelando/tesla-energy-stats-collector/connect"
	"github.com/iwvelando/tesla-energy-stats-collector/influxdb"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// CliInputs holds the data passed in via CLI parameters
type CliInputs struct {
	Config string
}

func main() {

	cliInputs := CliInputs{}
	flags := flag.NewFlagSet("tesla-energy-stats-collector", 0)
	flags.StringVar(&cliInputs.Config, "config", "config.yaml", "Set the location for the YAML config file")
	flags.Parse(os.Args[1:])

	conf, err := config.LoadConfiguration(cliInputs.Config)
	if err != nil {
		log.WithFields(log.Fields{
			"op":    "config.LoadConfiguration",
			"error": err,
		}).Fatal("failed to parse configuration")
	}

	tesla, refreshTime, err := connect.Auth(conf)
	if err != nil {
		log.WithFields(log.Fields{
			"op":    "connect.Auth",
			"error": err,
		}).Fatal("failed to authenticate to Tesla energy gateway")
	}
	defer tesla.CloseIdleConnections()

	influxClient, writeAPI, err := influxdb.Connect(conf)
	if err != nil {
		log.WithFields(log.Fields{
			"op":    "influxdb.Connect",
			"error": err,
		}).Fatal("failed to authenticate to InfluxDB")
	}
	defer influxClient.Close()
	defer writeAPI.Flush()

	errorsCh := writeAPI.Errors()

	// Monitor InfluxDB write errors
	go func() {
		for err := range errorsCh {
			log.WithFields(log.Fields{
				"op":    "influxdb.WriteAll",
				"error": err,
			}).Error("encountered error on writing to InfluxDB")
		}
	}()

	// Look for SIGTERM or SIGINT
	cancelCh := make(chan os.Signal, 1)
	signal.Notify(cancelCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {

			if time.Now().After(refreshTime) {
				tesla, refreshTime, err = connect.Auth(conf)
				if err != nil {
					log.WithFields(log.Fields{
						"op":    "connect.Auth",
						"error": err,
					}).Error("failed to refresh authentication to Tesla energy gateway")
				}
			}

			pollStartTime := time.Now()

			metrics, err := connect.GetAll(conf, tesla)
			if err != nil {
				var msg string
				if conf.Polling.ExitOnFail {
					msg = "failed to query all metrics, exiting"
				} else {
					msg = "failed to query all metrics, waiting for next poll"
				}
				log.WithFields(log.Fields{
					"op":    "connect.GetAll",
					"error": err,
				}).Error(msg)
				if conf.Polling.ExitOnFail {
					os.Exit(1)
				}
			} else {
				influxdb.WriteAll(conf, writeAPI, metrics)
			}

			timeRemaining := conf.Polling.Interval*time.Second - time.Since(pollStartTime)
			time.Sleep(time.Duration(timeRemaining))
			continue

		}
	}()

	sig := <-cancelCh
	log.WithFields(log.Fields{
		"op": "main",
	}).Info(fmt.Sprintf("caught signal %v, flushing data to InfluxDB", sig))
	writeAPI.Flush()

}
