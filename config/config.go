// Package config defines the data structures related to configuration and
// includes functions for modifying the loading and parsing the config.
package config

import (
	"fmt"
	"github.com/spf13/viper"
	"time"
)

// Configuration holds all configuration for finance-forecast.
type Configuration struct {
	TeslaGateway TeslaGateway
	InfluxDB     InfluxDB
	Polling      Polling
}

// TeslaGateway holds the Tesla Gateway connection parameters
type TeslaGateway struct {
	Email         string
	Password      string
	Address       string
	SkipVerifySsl bool
}

// InfluxDB holds the connection parameters for InfluxDB
type InfluxDB struct {
	Address           string
	Username          string
	Password          string
	MeasurementPrefix string
	Database          string
	RetentionPolicy   string
	Token             string
	Organization      string
	Bucket            string
	SkipVerifySsl     bool
	FlushInterval     uint
}

// Polling holds parameters related to how we poll the Tesla Gateway
type Polling struct {
	Interval   time.Duration
	ExitOnFail bool
}

// LoadConfiguration takes a file path as input and loads the YAML-formatted
// configuration there.
func LoadConfiguration(configPath string) (*Configuration, error) {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()

	viper.SetConfigType("yml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	var configuration Configuration
	err := viper.Unmarshal(&configuration)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %s", err)
	}

	return &configuration, nil
}
