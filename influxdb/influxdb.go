package influxdb

import (
	"crypto/tls"
	"fmt"
	influx "github.com/influxdata/influxdb-client-go/v2"
	influxAPI "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/iwvelando/tesla-energy-stats-collector/config"
	"github.com/iwvelando/tesla-energy-stats-collector/model"
	"strings"
)

// Connect authenticates to InfluxDB and returns a client
func Connect(conf *config.Configuration) (influx.Client, influxAPI.WriteAPI, error) {
	var auth string
	if conf.InfluxDB.Token != "" {
		auth = conf.InfluxDB.Token
	} else if conf.InfluxDB.Username != "" && conf.InfluxDB.Password != "" {
		auth = fmt.Sprintf("%s:%s", conf.InfluxDB.Username, conf.InfluxDB.Password)
	} else {
		auth = ""
	}

	var writeDest string
	if conf.InfluxDB.Bucket != "" {
		writeDest = conf.InfluxDB.Bucket
	} else if conf.InfluxDB.Database != "" && conf.InfluxDB.RetentionPolicy != "" {
		writeDest = fmt.Sprintf("%s/%s", conf.InfluxDB.Database, conf.InfluxDB.RetentionPolicy)
	} else {
		return nil, nil, fmt.Errorf("must configure at least one of bucket or database/retention policy")
	}

	if conf.InfluxDB.FlushInterval == 0 {
		conf.InfluxDB.FlushInterval = 30
	}

	options := influx.DefaultOptions().
		SetFlushInterval(1000 * conf.InfluxDB.FlushInterval).
		SetTLSConfig(&tls.Config{
			InsecureSkipVerify: conf.InfluxDB.SkipVerifySsl,
		})
	client := influx.NewClientWithOptions(conf.InfluxDB.Address, auth, options)

	writeAPI := client.WriteAPI(conf.InfluxDB.Organization, writeDest)

	return client, writeAPI, nil
}

// WriteAll writes the Teg data structure into InfluxDB
func WriteAll(conf *config.Configuration, writeAPI influxAPI.WriteAPI, metrics model.Teg) error {

	// Meters data
	p := influx.NewPoint(
		conf.InfluxDB.MeasurementPrefix+"energy_meters",
		map[string]string{
			"gateway_id":        metrics.Status.GatewayID,
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
			"measured_frequency":              metrics.SiteInfo.MeasuredFrequency,
			"max_system_energy_kwh":           metrics.SiteInfo.MaxSystemEnergyKwh,
			"max_system_power_kw":             metrics.SiteInfo.MaxSystemPowerKw,
			"max_site_meter_power_kw":         metrics.SiteInfo.MaxSiteMeterPowerKw,
			"min_site_meter_power_kw":         metrics.SiteInfo.MinSiteMeterPowerKw,
			"nominal_system_energy_kwh":       metrics.SiteInfo.NominalSystemEnergyKwh,
			"nominal_system_power_kw":         metrics.SiteInfo.NominalSystemPowerKw,
			"panel_max_current":               metrics.SiteInfo.PanelMaxCurrent,
			"grid_voltage_setting":            metrics.SiteInfo.GridCode.GridVoltageSetting,
			"grid_frequency_setting":          metrics.SiteInfo.GridCode.GridFreqSetting,
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
		conf.InfluxDB.MeasurementPrefix+"energy_powerwalls",
		map[string]string{
			"gateway_id":        metrics.Status.GatewayID,
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
			"sync_updating":                 metrics.Powerwalls.Sync.Updating,
			"charge_percent":                metrics.SystemStateOfEnergy.Percentage,
		},
		metrics.Powerwalls.Timestamp)

	writeAPI.WritePoint(p)

	// Overall powerwall sync diagnostics
	p = influx.NewPoint(
		conf.InfluxDB.MeasurementPrefix+"energy_powerwalls",
		map[string]string{
			"diagnostic":        metrics.Powerwalls.Sync.CommissioningDiagnostic.Name,
			"category":          metrics.Powerwalls.Sync.CommissioningDiagnostic.Category,
			"gateway_id":        metrics.Status.GatewayID,
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
		conf.InfluxDB.MeasurementPrefix+"energy_powerwalls",
		map[string]string{
			"diagnostic":        metrics.Powerwalls.Sync.UpdateDiagnostic.Name,
			"category":          metrics.Powerwalls.Sync.UpdateDiagnostic.Category,
			"gateway_id":        metrics.Status.GatewayID,
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

	// Powerwall diagnostic check results
	for _, check := range metrics.Powerwalls.Sync.CommissioningDiagnostic.Checks {
		p = influx.NewPoint(
			conf.InfluxDB.MeasurementPrefix+"energy_powerwalls",
			map[string]string{
				"check_name":        check.Name,
				"diagnostic":        metrics.Powerwalls.Sync.CommissioningDiagnostic.Name,
				"category":          metrics.Powerwalls.Sync.CommissioningDiagnostic.Category,
				"gateway_id":        metrics.Status.GatewayID,
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
			conf.InfluxDB.MeasurementPrefix+"energy_powerwalls",
			map[string]string{
				"check_name":        check.Name,
				"diagnostic":        metrics.Powerwalls.Sync.UpdateDiagnostic.Name,
				"category":          metrics.Powerwalls.Sync.UpdateDiagnostic.Category,
				"gateway_id":        metrics.Status.GatewayID,
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

	// Overall powerwall usage information
	p = influx.NewPoint(
		conf.InfluxDB.MeasurementPrefix+"energy_powerwalls",
		map[string]string{
			"gateway_id":        metrics.Status.GatewayID,
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
			"battery_target_power":                metrics.SystemStatus.BatteryTargetPower,
			"battery_target_reactive_power":       metrics.SystemStatus.BatteryTargetReactivePower,
			"nominal_full_pack_energy":            metrics.SystemStatus.NominalFullPackEnergyWattHours,
			"nominal_energy_remaining_watt_hours": metrics.SystemStatus.NominalEnergyRemainingWattHours,
			"max_power_energy_remaining":          metrics.SystemStatus.MaxPowerEnergyRemaining,
			"max_power_energy_to_be_charged":      metrics.SystemStatus.MaxPowerEnergyToBeCharged,
			"max_charge_power":                    metrics.SystemStatus.MaxChargePowerWatts,
			"max_discharge_power":                 metrics.SystemStatus.MaxDischargePowerWatts,
			"max_apparent_power":                  metrics.SystemStatus.MaxApparentPower,
			"instantaneous_max_discharge_power":   metrics.SystemStatus.InstantaneousMaxDischargePower,
			"instantaneous_max_charge_power":      metrics.SystemStatus.InstantaneousMaxChargePower,
			"grid_services_power":                 metrics.SystemStatus.GridServicesPower,
			"system_island_state":                 metrics.SystemStatus.SystemIslandState,
			"available_blocks":                    metrics.SystemStatus.AvailableBlocks,
			"ffr_power_availability_high":         metrics.SystemStatus.FfrPowerAvailabilityHigh,
			"ffr_power_availability_low":          metrics.SystemStatus.FfrPowerAvailabilityLow,
			"load_charge_constraint":              metrics.SystemStatus.LoadChargeConstraint,
			"max_sustained_ramp_rate":             metrics.SystemStatus.MaxSustainedRampRate,
			"can_reboot":                          metrics.SystemStatus.CanReboot,
			"smart_inv_delta_p":                   metrics.SystemStatus.SmartInvDeltaP,
			"smart_inv_delta_q":                   metrics.SystemStatus.SmartInvDeltaQ,
			"system_status_updating":              metrics.SystemStatus.Updating,
			"last_toggle_timestamp":               metrics.SystemStatus.LastToggleTimestamp.UnixNano(),
			"solar_real_power_limit":              metrics.SystemStatus.SolarRealPowerLimit,
			"score":                               metrics.SystemStatus.Score,
			"blocks_controlled":                   metrics.SystemStatus.BlocksControlled,
			"primary":                             metrics.SystemStatus.Primary,
			"auxiliary_load":                      metrics.SystemStatus.AuxiliaryLoad,
			"all_enable_lines_high":               metrics.SystemStatus.AllEnableLinesHigh,
			"inverter_nominal_usable_power":       metrics.SystemStatus.InverterNominalUsablePowerWatts,
			"expected_energy_remaining":           metrics.SystemStatus.ExpectedEnergyRemaining,
		},
		metrics.SystemStatus.Timestamp)

	writeAPI.WritePoint(p)

	// Individual powerwall usage information
	for _, block := range metrics.SystemStatus.BatteryBlocks {
		powerwallChargePercent := 0.0
		if block.NominalFullPackEnergy > 0 {
			powerwallChargePercent = float64(block.NominalEnergyRemainingWattHours) / float64(block.NominalFullPackEnergy) * 100.0
		}
		p = influx.NewPoint(
			conf.InfluxDB.MeasurementPrefix+"energy_powerwalls",
			map[string]string{
				"powerwall_part_number":   block.PackagePartNumber,
				"powerwall_serial_number": block.PackageSerialNumber,
				"gateway_id":              metrics.Status.GatewayID,
				"firmware_version":        metrics.Status.FirmwareVersion,
				"firmware_git_hash":       metrics.Status.FirmwareGitHash,
				"sync_type":               metrics.Status.SyncType,
				"site_name":               metrics.SiteInfo.SiteName,
				"site_grid_code":          metrics.SiteInfo.GridCode.GridCode,
				"site_country":            metrics.SiteInfo.GridCode.Country,
				"site_state":              metrics.SiteInfo.GridCode.State,
				"site_utility":            metrics.SiteInfo.GridCode.Utility,
			},
			map[string]interface{}{
				"powerwall_pinv_state":               block.PinvState,
				"powerwall_pinv_grid_state":          block.PinvGridState,
				"powerwall_nominal_energy_remaining": block.NominalEnergyRemainingWattHours,
				"powerwall_nominal_full_pack_energy": block.NominalFullPackEnergy,
				"powerwall_charge_percent":           powerwallChargePercent,
				"powerwall_p_out":                    block.POut,
				"qowerwall_q_out":                    block.QOut,
				"powerwall_v_out":                    block.VOut,
				"powerwall_f_out":                    block.FOut,
				"powerwall_i_out":                    block.IOut,
				"powerwall_energy_charged":           block.EnergyCharged,
				"powerwall_energy_discharged":        block.EnergyDischarged,
				"powerwall_off_grid":                 block.OffGrid,
				"powerwall_vf_mode":                  block.VfMode,
				"powerwall_wobble_detected":          block.WobbleDetected,
				"powerwall_charge_power_clamped":     block.ChargePowerClamped,
				"powerwall_backup_ready":             block.BackupReady,
				"powerwall_op_seq_state":             block.OpSeqState,
				"powerwall_disabled_reasons":         strings.Join(block.DisabledReasons[:], ","),
			},
			metrics.SystemStatus.Timestamp)

		writeAPI.WritePoint(p)
	}

	// Overall site information and configuration
	p = influx.NewPoint(
		conf.InfluxDB.MeasurementPrefix+"energy_configuration",
		map[string]string{
			"gateway_id":        metrics.Status.GatewayID,
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
			"mode":                          metrics.Operation.RealMode,
			"backup_reserve_percent":        metrics.Operation.BackupReservePercent,
			"freq_shift_load_shed_soe":      metrics.Operation.FreqShiftLoadShedSoe,
			"freq_shift_load_shed_delta_f":  metrics.Operation.FreqShiftLoadShedDeltaF,
			"net_meter_mode":                metrics.SiteInfo.NetMeterMode,
			"sitemaster_status":             metrics.Sitemaster.Status,
			"sitemaster_running":            metrics.Sitemaster.Running,
			"sitemaster_connected_to_tesla": metrics.Sitemaster.ConnectedToTesla,
			"sitemaster_power_supply_mode":  metrics.Sitemaster.PowerSupplyMode,
			"sitemaster_can_reboot":         metrics.Sitemaster.CanReboot,
			"grid_status":                   metrics.SystemGridStatus.GridStatus,
			"grid_services_active":          metrics.SystemGridStatus.GridServicesActive,
		},
		metrics.Operation.Timestamp)

	writeAPI.WritePoint(p)

	// Overall network diagnostics
	p = influx.NewPoint(
		conf.InfluxDB.MeasurementPrefix+"energy_network",
		map[string]string{
			"diagnostic":        metrics.NetworkConnectionTests.Name,
			"category":          metrics.NetworkConnectionTests.Category,
			"gateway_id":        metrics.Status.GatewayID,
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
			"disruptive": metrics.NetworkConnectionTests.Disruptive,
			"alert":      metrics.NetworkConnectionTests.Alert,
		},
		metrics.NetworkConnectionTests.Timestamp)

	writeAPI.WritePoint(p)

	// Network connectivity tests
	for _, check := range metrics.NetworkConnectionTests.Checks {
		p = influx.NewPoint(
			conf.InfluxDB.MeasurementPrefix+"energy_network",
			map[string]string{
				"check_name":        check.Name,
				"diagnostic":        metrics.NetworkConnectionTests.Name,
				"category":          metrics.NetworkConnectionTests.Category,
				"gateway_id":        metrics.Status.GatewayID,
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
			},
			metrics.NetworkConnectionTests.Timestamp)

		writeAPI.WritePoint(p)
	}

	// System status grid fault readings
	var valueString string
	for _, fault := range metrics.SystemStatus.GridFaults {
		for _, decodedAlert := range fault.DecodedAlert {
			switch decodedAlert.Value.(type) {
			case float64:
				valueString = fmt.Sprintf("%f", decodedAlert.Value.(float64))
			default:
				valueString = decodedAlert.Value.(string)
			}
			p = influx.NewPoint(
				conf.InfluxDB.MeasurementPrefix+"energy_faults",
				map[string]string{
					"fault_name":        fault.AlertName,
					"fault_subname":     decodedAlert.Name,
					"fault_units":       decodedAlert.Units,
					"gateway_id":        metrics.Status.GatewayID,
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
					"grid_fault_ts":                  fault.Timestamp,
					"grid_fault_isfault":             fault.AlertIsFault,
					"grid_fault_alert_raw":           fault.AlertRaw,
					"grid_fault_ecu_type":            fault.EcuType,
					"grid_fault_ecu_part_number":     fault.EcuPackagePartNumber,
					"grid_fault_ecu_serial_number":   fault.EcuPackageSerialNumber,
					"grid_fault_decoded_alert_value": valueString,
				},
				metrics.SystemStatus.Timestamp)

			writeAPI.WritePoint(p)
		}
	}

	return nil
}
