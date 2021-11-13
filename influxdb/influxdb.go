package influxdb

import (
	"crypto/tls"
	"fmt"
	influx "github.com/influxdata/influxdb-client-go/v2"
	influxAPI "github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/iwvelando/tesla-energy-stats-collector/config"
	"github.com/iwvelando/tesla-energy-stats-collector/model"
	"strconv"
	"strings"
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

	if config.InfluxDB.FlushInterval == 0 {
		config.InfluxDB.FlushInterval = 30
	}

	options := influx.DefaultOptions().
		SetFlushInterval(1000 * config.InfluxDB.FlushInterval).
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

	// Powerwall diagnostic check results
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

	// Overall powerwall usage information
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
			config.InfluxDB.MeasurementPrefix+"energy_powerwalls",
			map[string]string{
				"powerwall_part_number":   block.PackagePartNumber,
				"powerwall_serial_number": block.PackageSerialNumber,
				"gateway_id":              metrics.Status.GatewayId,
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
		config.InfluxDB.MeasurementPrefix+"energy_configuration",
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
		config.InfluxDB.MeasurementPrefix+"energy_network",
		map[string]string{
			"diagnostic":        metrics.NetworkConnectionTests.Name,
			"category":          metrics.NetworkConnectionTests.Category,
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
			"disruptive": metrics.NetworkConnectionTests.Disruptive,
			"alert":      metrics.NetworkConnectionTests.Alert,
		},
		metrics.NetworkConnectionTests.Timestamp)

	writeAPI.WritePoint(p)

	// Network connectivity tests
	for _, check := range metrics.NetworkConnectionTests.Checks {
		p = influx.NewPoint(
			config.InfluxDB.MeasurementPrefix+"energy_network",
			map[string]string{
				"check_name":        check.Name,
				"diagnostic":        metrics.NetworkConnectionTests.Name,
				"category":          metrics.NetworkConnectionTests.Category,
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
			},
			metrics.NetworkConnectionTests.Timestamp)

		writeAPI.WritePoint(p)
	}

	// Overall solar information
	p = influx.NewPoint(
		config.InfluxDB.MeasurementPrefix+"energy_pv",
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
			"pv_power_limit":                                 metrics.SolarPowerwall.PvPowerLimit,
			"power_status_setpoint":                          metrics.SolarPowerwall.PowerStatusSetpoint,
			"pvac_state":                                     metrics.SolarPowerwall.PvacStatus.State,
			"pvac_disabled":                                  metrics.SolarPowerwall.PvacStatus.Disabled,
			"pvac_disabled_reasons":                          strings.Join(metrics.SolarPowerwall.PvacStatus.DisabledReasons[:], ","),
			"pvac_grid_state":                                metrics.SolarPowerwall.PvacStatus.GridState,
			"pvac_inv_state":                                 metrics.SolarPowerwall.PvacStatus.InvState,
			"pvac_v_out":                                     metrics.SolarPowerwall.PvacStatus.VOut,
			"pvac_f_out":                                     metrics.SolarPowerwall.PvacStatus.FOut,
			"pvac_p_out":                                     metrics.SolarPowerwall.PvacStatus.POut,
			"pvac_q_out":                                     metrics.SolarPowerwall.PvacStatus.QOut,
			"pvac_i_out":                                     metrics.SolarPowerwall.PvacStatus.IOut,
			"pvac_alerts_last_rx_time":                       metrics.SolarPowerwall.PvacAlerts.LastRxTime.UnixNano(),
			"pvac_alerts_receive_mux_bitmask":                metrics.SolarPowerwall.PvacAlerts.ReceivedMuxBitmask,
			"pvac_alerts_a001_inv_l1_hw_overcurrent":         metrics.SolarPowerwall.PvacAlerts.PVACA001InvL1HWOvercurrent,
			"pvac_alerts_a002_inv_l2_hw_overcurrent":         metrics.SolarPowerwall.PvacAlerts.PVACA002InvL2HWOvercurrent,
			"pvac_alerts_a003_inv_hvbus_hw_overvoltage":      metrics.SolarPowerwall.PvacAlerts.PVACA003InvHVBusHWOvervoltage,
			"pvac_alerts_a004_pv_hw_cmpss_oc_stga":           metrics.SolarPowerwall.PvacAlerts.PVACA004PvHWCMPSSOCSTGA,
			"pvac_alerts_a005_pv_hw_cmpss_oc_stgb":           metrics.SolarPowerwall.PvacAlerts.PVACA005PvHWCMPSSOCSTGB,
			"pvac_alerts_a006_pv_hw_cmpss_oc_stgc":           metrics.SolarPowerwall.PvacAlerts.PVACA006PvHWCMPSSOCSTGC,
			"pvac_alerts_a007_pv_hw_cmpss_oc_stgd":           metrics.SolarPowerwall.PvacAlerts.PVACA007PvHWCMPSSOCSTGD,
			"pvac_alerts_a008_inv_hvbus_undervoltage":        metrics.SolarPowerwall.PvacAlerts.PVACA008InvHVBusUndervoltage,
			"pvac_alerts_a009_swappboot":                     metrics.SolarPowerwall.PvacAlerts.PVACA009SwAppBoot,
			"pvac_alerts_a010_inv_ac_overvoltage":            metrics.SolarPowerwall.PvacAlerts.PVACA010InvACOvervoltage,
			"pvac_alerts_a011_inv_ac_undervoltage":           metrics.SolarPowerwall.PvacAlerts.PVACA011InvACUndervoltage,
			"pvac_alerts_a012_inv_ac_overfrequency":          metrics.SolarPowerwall.PvacAlerts.PVACA012InvACOverfrequency,
			"pvac_alerts_a013_inv_ac_underfrequency":         metrics.SolarPowerwall.PvacAlerts.PVACA013InvACUnderfrequency,
			"pvac_alerts_a014_pvs_disabled_relay":            metrics.SolarPowerwall.PvacAlerts.PVACA014PVSDisabledRelay,
			"pvac_alerts_a015_pv_hw_allegro_oc_stga":         metrics.SolarPowerwall.PvacAlerts.PVACA015PvHWAllegroOCSTGA,
			"pvac_alerts_a016_pv_hw_allegro_oc_stgb":         metrics.SolarPowerwall.PvacAlerts.PVACA016PvHWAllegroOCSTGB,
			"pvac_alerts_a017_pv_hw_allegro_oc_stgc":         metrics.SolarPowerwall.PvacAlerts.PVACA017PvHWAllegroOCSTGC,
			"pvac_alerts_a018_pv_hw_allegro_oc_stgd":         metrics.SolarPowerwall.PvacAlerts.PVACA018PvHWAllegroOCSTGD,
			"pvac_alerts_a019_ambient_overtemperature":       metrics.SolarPowerwall.PvacAlerts.PVACA019AmbientOvertemperature,
			"pvac_alerts_a020_dsp_overtemperature":           metrics.SolarPowerwall.PvacAlerts.PVACA020DspOvertemperature,
			"pvac_alerts_a021_dcac_heatsink_overtemperature": metrics.SolarPowerwall.PvacAlerts.PVACA021DcacHeatsinkOvertemperature,
			"pvac_alerts_a022_mppt_heatsink_overtemperature": metrics.SolarPowerwall.PvacAlerts.PVACA022MpptHeatsinkOvertemperature,
			"pvac_alerts_a023_unused":                        metrics.SolarPowerwall.PvacAlerts.PVACA023Unused,
			"pvac_alerts_a024_pvacrx_command_mia":            metrics.SolarPowerwall.PvacAlerts.PVACA024PVACrxCommandMia,
			"pvac_alerts_a025_pvs_status_mia":                metrics.SolarPowerwall.PvacAlerts.PVACA025PVSStatusMia,
			"pvac_alerts_a026_inv_ac_peak_overvoltage":       metrics.SolarPowerwall.PvacAlerts.PVACA026InvACPeakOvervoltage,
			"pvac_alerts_a027_inv_k1_relay_welded":           metrics.SolarPowerwall.PvacAlerts.PVACA027InvK1RelayWelded,
			"pvac_alerts_a028_inv_k2_relay_welded":           metrics.SolarPowerwall.PvacAlerts.PVACA028InvK2RelayWelded,
			"pvac_alerts_a029_pump_faulted":                  metrics.SolarPowerwall.PvacAlerts.PVACA029PumpFaulted,
			"pvac_alerts_a030_fan_faulted":                   metrics.SolarPowerwall.PvacAlerts.PVACA030FanFaulted,
			"pvac_alerts_a031_vfcheck_ov":                    metrics.SolarPowerwall.PvacAlerts.PVACA031VFCheckOV,
			"pvac_alerts_a032_vfcheck_uv":                    metrics.SolarPowerwall.PvacAlerts.PVACA032VFCheckUV,
			"pvac_alerts_a033_vfcheck_of":                    metrics.SolarPowerwall.PvacAlerts.PVACA033VFCheckOF,
			"pvac_alerts_a034_vfcheck_uf":                    metrics.SolarPowerwall.PvacAlerts.PVACA034VFCheckUF,
			"pvac_alerts_a035_vfcheck_rocof":                 metrics.SolarPowerwall.PvacAlerts.PVACA035VFCheckRoCoF,
			"pvac_alerts_a036_inv_lost_il_control":           metrics.SolarPowerwall.PvacAlerts.PVACA036InvLostILControl,
			"pvac_alerts_a037_pvs_processor_nerror":          metrics.SolarPowerwall.PvacAlerts.PVACA037PVSProcessorNERROR,
			"pvac_alerts_a038_inv_failed_xcap_precharge":     metrics.SolarPowerwall.PvacAlerts.PVACA038InvFailedXcapPrecharge,
			"pvac_alerts_a039_inv_hvbus_sw_overvoltage":      metrics.SolarPowerwall.PvacAlerts.PVACA039InvHVBusSWOvervoltage,
			"pvac_alerts_a040_pump_correction_saturated":     metrics.SolarPowerwall.PvacAlerts.PVACA040PumpCorrectionSaturated,
			"pvac_alerts_a041_excess_pv_clamp_triggered":     metrics.SolarPowerwall.PvacAlerts.PVACA041ExcessPVClampTriggered,
			"pvs_state":                                  metrics.SolarPowerwall.PvsStatus.State,
			"pvs_disabled":                               metrics.SolarPowerwall.PvsStatus.Disabled,
			"pvs_enable_output":                          metrics.SolarPowerwall.PvsStatus.EnableOutput,
			"pvs_v_ll":                                   metrics.SolarPowerwall.PvsStatus.Vll,
			"pvs_self_test_state":                        metrics.SolarPowerwall.PvsStatus.SelfTestState,
			"pvs_alerts_last_rx_time":                    metrics.SolarPowerwall.PvsAlerts.LastRxTime.UnixNano(),
			"pvs_alerts_receive_mux_bitmask":             metrics.SolarPowerwall.PvsAlerts.ReceivedMuxBitmask,
			"pvs_alerts_a001_watchdogreset":              metrics.SolarPowerwall.PvsAlerts.PVSA001WatchdogReset,
			"pvs_alerts_a002_sw_app_boot":                metrics.SolarPowerwall.PvsAlerts.PVSA002SWAppBoot,
			"pvs_alerts_a003_v12voutofbounds":            metrics.SolarPowerwall.PvsAlerts.PVSA003V12vOutOfBounds,
			"pvs_alerts_a004_v1v5outofbounds":            metrics.SolarPowerwall.PvsAlerts.PVSA004V1v5OutOfBounds,
			"pvs_alerts_a005_vafdrefoutofbounds":         metrics.SolarPowerwall.PvsAlerts.PVSA005VAfdRefOutOfBounds,
			"pvs_alerts_a006_gfovercurrent300":           metrics.SolarPowerwall.PvsAlerts.PVSA006GfOvercurrent300,
			"pvs_alerts_a007_unused_7":                   metrics.SolarPowerwall.PvsAlerts.PVSA007UNUSED7,
			"pvs_alerts_a008_unused_8":                   metrics.SolarPowerwall.PvsAlerts.PVSA008UNUSED8,
			"pvs_alerts_a009_gfovercurrent030":           metrics.SolarPowerwall.PvsAlerts.PVSA009GfOvercurrent030,
			"pvs_alerts_a010_pvisolationtotal":           metrics.SolarPowerwall.PvsAlerts.PVSA010PvIsolationTotal,
			"pvs_alerts_a011_pvisolationstringa":         metrics.SolarPowerwall.PvsAlerts.PVSA011PvIsolationStringA,
			"pvs_alerts_a012_pvisolationstringb":         metrics.SolarPowerwall.PvsAlerts.PVSA012PvIsolationStringB,
			"pvs_alerts_a013_pvisolationstringc":         metrics.SolarPowerwall.PvsAlerts.PVSA013PvIsolationStringC,
			"pvs_alerts_a014_pvisolationstringd":         metrics.SolarPowerwall.PvsAlerts.PVSA014PvIsolationStringD,
			"pvs_alerts_a015_selftestgroundfault":        metrics.SolarPowerwall.PvsAlerts.PVSA015SelfTestGroundFault,
			"pvs_alerts_a016_esmfault":                   metrics.SolarPowerwall.PvsAlerts.PVSA016ESMFault,
			"pvs_alerts_a017_mcistringa":                 metrics.SolarPowerwall.PvsAlerts.PVSA017MciStringA,
			"pvs_alerts_a018_mcistringb":                 metrics.SolarPowerwall.PvsAlerts.PVSA018MciStringB,
			"pvs_alerts_a019_mcistringc":                 metrics.SolarPowerwall.PvsAlerts.PVSA019MciStringC,
			"pvs_alerts_a020_mcistringd":                 metrics.SolarPowerwall.PvsAlerts.PVSA020MciStringD,
			"pvs_alerts_a021_rapidshutdown":              metrics.SolarPowerwall.PvsAlerts.PVSA021RapidShutdown,
			"pvs_alerts_a022_mci1signallevel":            metrics.SolarPowerwall.PvsAlerts.PVSA022Mci1SignalLevel,
			"pvs_alerts_a023_mci2signallevel":            metrics.SolarPowerwall.PvsAlerts.PVSA023Mci2SignalLevel,
			"pvs_alerts_a024_mci3signallevel":            metrics.SolarPowerwall.PvsAlerts.PVSA024Mci3SignalLevel,
			"pvs_alerts_a025_mci4signallevel":            metrics.SolarPowerwall.PvsAlerts.PVSA025Mci4SignalLevel,
			"pvs_alerts_a026_mci1pvvoltage":              metrics.SolarPowerwall.PvsAlerts.PVSA026Mci1PvVoltage,
			"pvs_alerts_a027_mci2pvvoltage":              metrics.SolarPowerwall.PvsAlerts.PVSA027Mci2PvVoltage,
			"pvs_alerts_a028_systeminitfailed":           metrics.SolarPowerwall.PvsAlerts.PVSA028SystemInitFailed,
			"pvs_alerts_a029_pvarcfault":                 metrics.SolarPowerwall.PvsAlerts.PVSA029PvArcFault,
			"pvs_alerts_a030_vdcov":                      metrics.SolarPowerwall.PvsAlerts.PVSA030VDcOv,
			"pvs_alerts_a031_mci3pvvoltage":              metrics.SolarPowerwall.PvsAlerts.PVSA031Mci3PvVoltage,
			"pvs_alerts_a032_mci4pvvoltage":              metrics.SolarPowerwall.PvsAlerts.PVSA032Mci4PvVoltage,
			"pvs_alerts_a033_dataexception":              metrics.SolarPowerwall.PvsAlerts.PVSA033DataException,
			"pvs_alerts_a034_peimpedance":                metrics.SolarPowerwall.PvsAlerts.PVSA034PeImpedance,
			"pvs_alerts_a035_pvarcdetected":              metrics.SolarPowerwall.PvsAlerts.PVSA035PvArcDetected,
			"pvs_alerts_a036_pvarclockout":               metrics.SolarPowerwall.PvsAlerts.PVSA036PvArcLockout,
			"pvs_alerts_a037_pvarcfault2":                metrics.SolarPowerwall.PvsAlerts.PVSA037PvArcFault2,
			"pvs_alerts_a038_pvarcfault_selftest":        metrics.SolarPowerwall.PvsAlerts.PVSA038PvArcFaultSelfTest,
			"pvs_alerts_a039_selftestrelayfault":         metrics.SolarPowerwall.PvsAlerts.PVSA039SelfTestRelayFault,
			"pvs_alerts_a040_ledirrationalfault":         metrics.SolarPowerwall.PvsAlerts.PVSA040LEDIrrationalFault,
			"pvs_alerts_a041_mcipowerswitch":             metrics.SolarPowerwall.PvsAlerts.PVSA041MciPowerSwitch,
			"pvs_alerts_a042_mcipowerfault":              metrics.SolarPowerwall.PvsAlerts.PVSA042MciPowerFault,
			"pvs_alerts_a043_lockedpvstringsafety":       metrics.SolarPowerwall.PvsAlerts.PVSA043LockedPvStringSafety,
			"pvs_alerts_a044_faultstatepvstringsafety":   metrics.SolarPowerwall.PvsAlerts.PVSA044FaultStatePvStringSafety,
			"pvs_alerts_a045_relaycoilirrationalfault":   metrics.SolarPowerwall.PvsAlerts.PVSA045RelayCoilIrrationalFault,
			"pvs_alerts_a046_relaycoilirrationallockout": metrics.SolarPowerwall.PvsAlerts.PVSA046RelayCoilIrrationalLockout,
			"pvs_alerts_a047_acsensorirrationalfault":    metrics.SolarPowerwall.PvsAlerts.PVSA047AcSensorIrrationalFault,
			"pvs_alerts_a048_dcsensorirrationalfault":    metrics.SolarPowerwall.PvsAlerts.PVSA048DcSensorIrrationalFault,
			"pvs_alerts_a049_arcsignalmibspihealth":      metrics.SolarPowerwall.PvsAlerts.PVSA049ArcSignalMibspiHealth,
			"pvs_alerts_a050_relaycoilirrationalwarning": metrics.SolarPowerwall.PvsAlerts.PVSA050RelayCoilIrrationalWarning,
			"pvs_alerts_a051_dcbusshortcircuitdetected":  metrics.SolarPowerwall.PvsAlerts.PVSA051DcBusShortCircuitDetected,
			"pvs_alerts_a052_pvarcfault_preselftest":     metrics.SolarPowerwall.PvsAlerts.PVSA052PvArcFaultPreSelfTest,
		},
		metrics.Operation.Timestamp)

	writeAPI.WritePoint(p)

	// Solar string information
	for _, stringInverter := range metrics.SolarPowerwall.PvacStatus.StringVitals {
		p = influx.NewPoint(
			config.InfluxDB.MeasurementPrefix+"energy_pv",
			map[string]string{
				"string_id":         strconv.Itoa(stringInverter.StringId),
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
				"string_connected":        stringInverter.Connected,
				"string_measured_voltage": stringInverter.MeasuredVoltage,
				"string_current":          stringInverter.Current,
				"string_measured_power":   stringInverter.MeasuredPower,
			},
			metrics.Operation.Timestamp)

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
				config.InfluxDB.MeasurementPrefix+"energy_faults",
				map[string]string{
					"fault_name":        fault.AlertName,
					"fault_subname":     decodedAlert.Name,
					"fault_units":       decodedAlert.Units,
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
