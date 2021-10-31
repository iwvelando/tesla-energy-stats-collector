package influxdb

import (
	"crypto/tls"
	"fmt"
	influx "github.com/influxdata/influxdb-client-go/v2"
	influxAPI "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/iwvelando/tesla-energy-stats-collector/config"
	"github.com/iwvelando/tesla-energy-stats-collector/model"
)

type InfluxWriteConfigError struct{}

func (r *InfluxWriteConfigError) Error() string {
	return "must configure at least one of bucket or database/retention policy"
}

func Connect(config *config.Configuration) (influx.Client, influxAPI.WriteAPI, error) {
	var auth string
	if config.InfluxDB.Token != "" {
		auth = config.InfluxDB.Token
	} else if config.InfluxDB.Username != "" && config.InfluxDB.Password != "" {
		auth = fmt.Sprintf("%s:%s", config.InfluxDB.Username, config.InfluxDB.Password)
	} else {
		auth = ""
	}

	var writeDest string
	if config.InfluxDB.Bucket != "" {
		writeDest = config.InfluxDB.Bucket
	} else if config.InfluxDB.Database != "" && config.InfluxDB.RetentionPolicy != "" {
		writeDest = fmt.Sprintf("%s/%s", config.InfluxDB.Database, config.InfluxDB.RetentionPolicy)
	} else {
		return nil, nil, &InfluxWriteConfigError{}
	}

	options := influx.DefaultOptions().
		SetTLSConfig(&tls.Config{
			InsecureSkipVerify: config.InfluxDB.SkipVerifySsl,
		})
	client := influx.NewClientWithOptions(config.InfluxDB.Address, auth, options)

	writeAPI := client.WriteAPI(config.InfluxDB.Organization, writeDest)

	return client, writeAPI, nil
}

func WriteAll(config *config.Configuration, writeAPI influxAPI.WriteAPI, metrics model.Teg) error {

	p := influx.NewPoint(
		config.InfluxDB.MeasurementPrefix+"energy_meters",
		map[string]string{
			"gateway_id":        metrics.Status.GatewayId,
			"firmware_version":  metrics.Status.FirmwareVersion,
			"firmware_git_hash": metrics.Status.FirmwareGitHash,
			"sync_type":         metrics.Status.SyncType,
			"meter_serial":      metrics.MetersStatus.Serial,
			"site_name":         metrics.SiteInfo.SiteName,
			"site_grid_code":    metrics.SiteInfo.GridCode.GridCode,
			"site_country":      metrics.SiteInfo.GridCode.Country,
			"site_state":        metrics.SiteInfo.GridCode.State,
			"site_utility":      metrics.SiteInfo.GridCode.Utility,
		},
		map[string]interface{}{
			"meter_status":                    metrics.MetersStatus.Status,
			"site_last_comm_time":             metrics.Meters.Site.LastCommunicationTime.UnixNano(),
			"site_instant_power":              metrics.Meters.Site.InstantPowerWatts,
			"site_instant_reactive_power":     metrics.Meters.Site.InstantReactivePowerWatts,
			"site_instant_apparent_power":     metrics.Meters.Site.InstantApparentPowerWatts,
			"site_frequency":                  metrics.Meters.Site.Frequency,
			"site_energy_exported":            metrics.Meters.Site.EnergyExportedWatts,
			"site_energy_imported":            metrics.Meters.Site.EnergyImportedWatts,
			"site_instant_average_voltage":    metrics.Meters.Site.InstantAverageVoltage,
			"site_instant_average_current":    metrics.Meters.Site.InstantAverageCurrent,
			"site_instant_total_current":      metrics.Meters.Site.InstantTotalCurrent,
			"battery_last_comm_time":          metrics.Meters.Battery.LastCommunicationTime.UnixNano(),
			"battery_instant_power":           metrics.Meters.Battery.InstantPowerWatts,
			"battery_instant_reactive_power":  metrics.Meters.Battery.InstantReactivePowerWatts,
			"battery_instant_apparent_power":  metrics.Meters.Battery.InstantApparentPowerWatts,
			"battery_frequency":               metrics.Meters.Battery.Frequency,
			"battery_energy_exported":         metrics.Meters.Battery.EnergyExportedWatts,
			"battery_energy_imported":         metrics.Meters.Battery.EnergyImportedWatts,
			"battery_instant_average_voltage": metrics.Meters.Battery.InstantAverageVoltage,
			"battery_instant_average_current": metrics.Meters.Battery.InstantAverageCurrent,
			"battery_instant_total_current":   metrics.Meters.Battery.InstantTotalCurrent,
			"load_last_comm_time":             metrics.Meters.Load.LastCommunicationTime.UnixNano(),
			"load_instant_power":              metrics.Meters.Load.InstantPowerWatts,
			"load_instant_reactive_power":     metrics.Meters.Load.InstantReactivePowerWatts,
			"load_instant_apparent_power":     metrics.Meters.Load.InstantApparentPowerWatts,
			"load_frequency":                  metrics.Meters.Load.Frequency,
			"load_energy_exported":            metrics.Meters.Load.EnergyExportedWatts,
			"load_energy_imported":            metrics.Meters.Load.EnergyImportedWatts,
			"load_instant_average_voltage":    metrics.Meters.Load.InstantAverageVoltage,
			"load_instant_average_current":    metrics.Meters.Load.InstantAverageCurrent,
			"load_instant_total_current":      metrics.Meters.Load.InstantTotalCurrent,
			"solar_last_comm_time":            metrics.Meters.Solar.LastCommunicationTime.UnixNano(),
			"solar_instant_power":             metrics.Meters.Solar.InstantPowerWatts,
			"solar_instant_reactive_power":    metrics.Meters.Solar.InstantReactivePowerWatts,
			"solar_instant_apparent_power":    metrics.Meters.Solar.InstantApparentPowerWatts,
			"solar_frequency":                 metrics.Meters.Solar.Frequency,
			"solar_energy_exported":           metrics.Meters.Solar.EnergyExportedWatts,
			"solar_energy_imported":           metrics.Meters.Solar.EnergyImportedWatts,
			"solar_instant_average_voltage":   metrics.Meters.Solar.InstantAverageVoltage,
			"solar_instant_average_current":   metrics.Meters.Solar.InstantAverageCurrent,
			"solar_instant_total_current":     metrics.Meters.Solar.InstantTotalCurrent,
		},
		metrics.Meters.Timestamp)

	writeAPI.WritePoint(p)

	return nil
}
