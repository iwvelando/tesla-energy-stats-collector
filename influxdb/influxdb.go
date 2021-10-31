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

	// Meters data
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

	// Overall powerwall info
	p = influx.NewPoint(
		config.InfluxDB.MeasurementPrefix+"energy_powerwalls",
		map[string]string{
			"gateway_id":        metrics.Status.GatewayId,
			"firmware_version":  metrics.Status.FirmwareVersion,
			"firmware_git_hash": metrics.Status.FirmwareGitHash,
			"sync_type":         metrics.Status.SyncType,
			"site_name":         metrics.SiteInfo.SiteName,
			"site_grid_code":    metrics.SiteInfo.GridCode.GridCode,
			"site_country":      metrics.SiteInfo.GridCode.Country,
			"site_state":        metrics.SiteInfo.GridCode.State,
			"site_utility":      metrics.SiteInfo.GridCode.Utility,
		},
		map[string]interface{}{
			"enumerating":                   metrics.Powerwalls.Enumerating,
			"updating":                      metrics.Powerwalls.Updating,
			"checking_if_offgrid":           metrics.Powerwalls.CheckingIfOffgrid,
			"running_phase_detection":       metrics.Powerwalls.RunningPhaseDetection,
			"bubble_shedding":               metrics.Powerwalls.BubbleShedding,
			"grid_qualifying":               metrics.Powerwalls.GridQualifying,
			"grid_code_validating":          metrics.Powerwalls.GridCodeValidating,
			"phase_detection_not_available": metrics.Powerwalls.PhaseDetectionNotAvailable,
			"on_grid_check_error":           metrics.Powerwalls.OnGridCheckError,
			"phase_detection_last_error":    metrics.Powerwalls.PhaseDetectionLastError,
			"sync_updating":                 metrics.Powerwalls.Sync,
			"charge_percent":                metrics.SystemStateOfEnergy.Percentage,
		},
		metrics.Powerwalls.Timestamp)

	writeAPI.WritePoint(p)

	// Overall powerwall sync diagnostics
	p = influx.NewPoint(
		config.InfluxDB.MeasurementPrefix+"energy_powerwalls",
		map[string]string{
			"diagnostic":        metrics.Powerwalls.Sync.CommissioningDiagnostic.Name,
			"category":          metrics.Powerwalls.Sync.CommissioningDiagnostic.Category,
			"gateway_id":        metrics.Status.GatewayId,
			"firmware_version":  metrics.Status.FirmwareVersion,
			"firmware_git_hash": metrics.Status.FirmwareGitHash,
			"sync_type":         metrics.Status.SyncType,
			"site_name":         metrics.SiteInfo.SiteName,
			"site_grid_code":    metrics.SiteInfo.GridCode.GridCode,
			"site_country":      metrics.SiteInfo.GridCode.Country,
			"site_state":        metrics.SiteInfo.GridCode.State,
			"site_utility":      metrics.SiteInfo.GridCode.Utility,
		},
		map[string]interface{}{
			"disruptive": metrics.Powerwalls.Sync.CommissioningDiagnostic.Disruptive,
			"alert":      metrics.Powerwalls.Sync.CommissioningDiagnostic.Alert,
		},
		metrics.Powerwalls.Timestamp)

	writeAPI.WritePoint(p)

	p = influx.NewPoint(
		config.InfluxDB.MeasurementPrefix+"energy_powerwalls",
		map[string]string{
			"diagnostic":        metrics.Powerwalls.Sync.UpdateDiagnostic.Name,
			"category":          metrics.Powerwalls.Sync.UpdateDiagnostic.Category,
			"gateway_id":        metrics.Status.GatewayId,
			"firmware_version":  metrics.Status.FirmwareVersion,
			"firmware_git_hash": metrics.Status.FirmwareGitHash,
			"sync_type":         metrics.Status.SyncType,
			"site_name":         metrics.SiteInfo.SiteName,
			"site_grid_code":    metrics.SiteInfo.GridCode.GridCode,
			"site_country":      metrics.SiteInfo.GridCode.Country,
			"site_state":        metrics.SiteInfo.GridCode.State,
			"site_utility":      metrics.SiteInfo.GridCode.Utility,
		},
		map[string]interface{}{
			"disruptive": metrics.Powerwalls.Sync.UpdateDiagnostic.Disruptive,
			"alert":      metrics.Powerwalls.Sync.UpdateDiagnostic.Alert,
		},
		metrics.Powerwalls.Timestamp)

	writeAPI.WritePoint(p)

	for _, check := range metrics.Powerwalls.Sync.CommissioningDiagnostic.Checks {
		p = influx.NewPoint(
			config.InfluxDB.MeasurementPrefix+"energy_powerwalls",
			map[string]string{
				"check_name":        check.Name,
				"diagnostic":        metrics.Powerwalls.Sync.CommissioningDiagnostic.Name,
				"category":          metrics.Powerwalls.Sync.CommissioningDiagnostic.Category,
				"gateway_id":        metrics.Status.GatewayId,
				"firmware_version":  metrics.Status.FirmwareVersion,
				"firmware_git_hash": metrics.Status.FirmwareGitHash,
				"sync_type":         metrics.Status.SyncType,
				"site_name":         metrics.SiteInfo.SiteName,
				"site_grid_code":    metrics.SiteInfo.GridCode.GridCode,
				"site_country":      metrics.SiteInfo.GridCode.Country,
				"site_state":        metrics.SiteInfo.GridCode.State,
				"site_utility":      metrics.SiteInfo.GridCode.Utility,
			},
			map[string]interface{}{
				"check_status":     check.Status,
				"check_start_time": check.StartTime.UnixNano(),
				"check_end_time":   check.EndTime.UnixNano(),
				"check_message":    check.Message,
			},
			metrics.Powerwalls.Timestamp)

		writeAPI.WritePoint(p)
	}

	for _, check := range metrics.Powerwalls.Sync.UpdateDiagnostic.Checks {
		p = influx.NewPoint(
			config.InfluxDB.MeasurementPrefix+"energy_powerwalls",
			map[string]string{
				"check_name":        check.Name,
				"diagnostic":        metrics.Powerwalls.Sync.UpdateDiagnostic.Name,
				"category":          metrics.Powerwalls.Sync.UpdateDiagnostic.Category,
				"gateway_id":        metrics.Status.GatewayId,
				"firmware_version":  metrics.Status.FirmwareVersion,
				"firmware_git_hash": metrics.Status.FirmwareGitHash,
				"sync_type":         metrics.Status.SyncType,
				"site_name":         metrics.SiteInfo.SiteName,
				"site_grid_code":    metrics.SiteInfo.GridCode.GridCode,
				"site_country":      metrics.SiteInfo.GridCode.Country,
				"site_state":        metrics.SiteInfo.GridCode.State,
				"site_utility":      metrics.SiteInfo.GridCode.Utility,
			},
			map[string]interface{}{
				"check_status":     check.Status,
				"check_start_time": check.StartTime.UnixNano(),
				"check_end_time":   check.EndTime.UnixNano(),
				"check_message":    check.Message,
			},
			metrics.Powerwalls.Timestamp)

		writeAPI.WritePoint(p)
	}

	return nil
}
