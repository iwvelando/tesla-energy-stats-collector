package model

import (
	"encoding/json"
	"time"
)

const dateTimeNano = "2006-01-02T15:04:05.999999999-07:00"
const dateTimeStatus = "2006-01-02 15:04:05 -0700"

// AuthPayload defines the payload for /api/login/Basic
type AuthPayload struct {
	Username   string
	Password   string
	Email      string
	ForceSmOff bool
}

// AuthResponse defines the response for /api/login/Basic
type AuthResponse struct {
	Email        string   `json:"email"`
	Firstname    string   `json:"firstname"`
	Lastname     string   `json:"lastname"`
	Roles        []string `json:"roles"`
	Token        string   `json:"token"`
	Provider     string   `json:"provider"`
	LoginTimeRaw string   `json:"loginTime"` // dateTimeNano
	LoginTime    time.Time
}

// ParseTime converts string times in AuthResponse to time.Time
func (r *AuthResponse) ParseTime() error {
	if r.LoginTimeRaw != "" {
		t, err := time.Parse(dateTimeNano, r.LoginTimeRaw)
		if err != nil {
			return err
		}
		r.LoginTime = t
	}

	return nil
}

// Teg defines the aggregated metrics
type Teg struct {
	Meters                 TegMeters
	MetersStatus           TegMetersStatus
	Operation              TegOperation
	Powerwalls             TegPowerwalls
	SiteInfo               TegSiteInfo
	Sitemaster             TegSitemaster
	Solars                 []TegSolars
	Status                 TegStatus
	NetworkConnectionTests TegNetworkConnectionTests
	SystemTesting          TegSystemTesting
	UpdateStatus           TegUpdateStatus
	SystemStatus           TegSystemStatus
	SystemGridStatus       TegSystemGridStatus
	SystemStateOfEnergy    TegSystemStateOfEnergy
}

// TegMeters defines the response for /api/meters/aggregates
type TegMeters struct {
	Timestamp time.Time
	Site      TegMetersAggregate `json:"site"`
	Battery   TegMetersAggregate `json:"battery"`
	Load      TegMetersAggregate `json:"load"`
	Solar     TegMetersAggregate `json:"solar"`
}

// TegMetersAggregate defines meters data underneath TegMeters
type TegMetersAggregate struct {
	LastCommunicationTimeRaw          string `json:"last_communication_time"` // dateTimeNano
	LastCommunicationTime             time.Time
	InstantPowerWatts                 float64 `json:"instant_power"`
	InstantReactivePowerWatts         float64 `json:"instant_reactive_power"`
	InstantApparentPowerWatts         float64 `json:"instant_apparent_power"`
	Frequency                         float64 `json:"frequency"`
	EnergyExportedWatts               float64 `json:"energy_exported"`
	EnergyImportedWatts               float64 `json:"energy_imported"`
	InstantAverageVoltage             float64 `json:"instant_average_voltage"`
	InstantAverageCurrent             float64 `json:"instant_average_current"`
	IACurrent                         float64 `json:"i_a_current"`
	IBCurrent                         float64 `json:"i_b_current"`
	ICCurrent                         float64 `json:"i_c_current"`
	LastPhaseVoltageCommunicationTime string  `json:"last_phase_voltage_communication_time"`
	LastPhasePowerCommunicationTime   string  `json:"last_phase_power_communication_time"`
	Timeout                           int     `json:"timeout"`
	NumMetersAggregated               int     `json:"num_meters_aggregated"`
	InstantTotalCurrent               float64 `json:"instant_total_current"`
}

// ParseTime converts string times in TegMeters to time.Time
func (r *TegMeters) ParseTime() error {
	if r.Site.LastCommunicationTimeRaw != "" {
		t, err := time.Parse(dateTimeNano, r.Site.LastCommunicationTimeRaw)
		if err == nil {
			r.Site.LastCommunicationTime = t
		}
	}

	if r.Battery.LastCommunicationTimeRaw != "" {
		t, err := time.Parse(dateTimeNano, r.Battery.LastCommunicationTimeRaw)
		if err == nil {
			r.Battery.LastCommunicationTime = t
		}
	}

	if r.Load.LastCommunicationTimeRaw != "" {
		t, err := time.Parse(dateTimeNano, r.Load.LastCommunicationTimeRaw)
		if err == nil {
			r.Load.LastCommunicationTime = t
		}
	}

	if r.Solar.LastCommunicationTimeRaw != "" {
		t, err := time.Parse(dateTimeNano, r.Solar.LastCommunicationTimeRaw)
		if err == nil {
			r.Solar.LastCommunicationTime = t
		}
	}

	return nil
}

// TegMetersStatus defines the response for /api/meters/status
type TegMetersStatus struct {
	Timestamp time.Time
	Status    string      `json:"status"`
	Errors    interface{} `json:"errors"`
	Serial    string      `json:"serial"`
}

// TegOperation defines the response for /api/operation
type TegOperation struct {
	Timestamp               time.Time
	RealMode                string  `json:"real_mode"`
	BackupReservePercent    float64 `json:"backup_reserve_percent"`
	FreqShiftLoadShedSoe    int     `json:"freq_shift_load_shed_soe"`
	FreqShiftLoadShedDeltaF float64 `json:"freq_shift_load_shed_delta_f"`
}

// TegPowerwalls defines the response for /api/powerwalls
type TegPowerwalls struct {
	Sync                       TegPowerwallsSync `json:"sync"`
	Timestamp                  time.Time
	Powerwalls                 []TegPowerwall `json:"powerwalls"`
	Msa                        interface{}    `json:"msa"`
	GatewayID                  string         `json:"gateway_din"`
	PhaseDetectionLastError    string         `json:"phase_detection_last_error"`
	OnGridCheckError           string         `json:"on_grid_check_error"`
	States                     interface{}    `json:"states"`
	BubbleShedding             bool           `json:"bubble_shedding"`
	GridCodeValidating         bool           `json:"grid_code_validating"`
	PhaseDetectionNotAvailable bool           `json:"phase_detection_not_available"`
	RunningPhaseDetection      bool           `json:"running_phase_detection"`
	CheckingIfOffgrid          bool           `json:"checking_if_offgrid"`
	Updating                   bool           `json:"updating"`
	Enumerating                bool           `json:"enumerating"`
	GridQualifying             bool           `json:"grid_qualifying"`
}

// TegPowerwall defines powerwall data underneath TegPowerwalls
type TegPowerwall struct {
	CommissioningDiagnostic     TegPowerwallDiagnostic `json:"commissioning_diagnostic"`
	UpdateDiagnostic            TegPowerwallDiagnostic `json:"update_diagnostic"`
	Type                        string                 `json:"Type"`
	PackagePartNumber           string                 `json:"PackagePartNumber"`
	GridState                   string                 `json:"grid_state"`
	PackageSerialNumber         string                 `json:"PackageSerialNumber"`
	Subtype                     string                 `json:"type"`
	BcType                      interface{}            `json:"bc_type"`
	GridReconnectionTimeSeconds int                    `json:"grid_reconnection_time_seconds"`
	UnderPhaseDetection         bool                   `json:"under_phase_detection"`
	Updating                    bool                   `json:"updating"`
	InConfig                    bool                   `json:"in_config"`
}

// TegPowerwallsSync defines high-level diagnostic data underneath TegPowerwalls
type TegPowerwallsSync struct {
	Updating                bool                   `json:"updating"`
	CommissioningDiagnostic TegPowerwallDiagnostic `json:"commissioning_diagnostic"`
	UpdateDiagnostic        TegPowerwallDiagnostic `json:"update_diagnostic"`
}

// TegPowerwallDiagnostic defines low-level diagnostic data underneath TegPowerwallsSync
type TegPowerwallDiagnostic struct {
	Checks     []TegPowerwallsCheck `json:"checks"`
	Name       string               `json:"name"`
	Category   string               `json:"category"`
	Inputs     interface{}          `json:"inputs"`
	Disruptive bool                 `json:"disruptive"`
	Alert      bool                 `json:"alert"`
}

// TegPowerwallsCheck defines checks data underneath TegPowerwallDiagnostic
type TegPowerwallsCheck struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	StartTimeRaw string `json:"start_time"` // dateTimeNano
	StartTime    time.Time
	EndTimeRaw   string `json:"end_time"` // dateTimeNano
	EndTime      time.Time
	Message      string      `json:"message"`
	Progress     int         `json:"progress"`
	Results      interface{} `json:"results"`
	Debug        interface{} `json:"debug"`
	Checks       interface{} `json:"checks"`
}

// ParseTime converts string times in TegPowerwalls to time.Time
func (r TegPowerwalls) ParseTime() error {
	for i := range r.Powerwalls {
		for j, check := range r.Powerwalls[i].CommissioningDiagnostic.Checks {
			if check.StartTimeRaw != "" {
				t, err := time.Parse(dateTimeNano, check.StartTimeRaw)
				if err != nil {
					return err
				}
				r.Powerwalls[i].CommissioningDiagnostic.Checks[j].StartTime = t
			}

			if check.EndTimeRaw != "" {
				t, err := time.Parse(dateTimeNano, check.EndTimeRaw)
				if err != nil {
					return err
				}
				r.Powerwalls[i].CommissioningDiagnostic.Checks[j].EndTime = t
			}
		}

		for j, check := range r.Powerwalls[i].UpdateDiagnostic.Checks {
			if check.StartTimeRaw != "" {
				t, err := time.Parse(dateTimeNano, check.StartTimeRaw)
				if err != nil {
					return err
				}
				r.Powerwalls[i].UpdateDiagnostic.Checks[j].StartTime = t
			}

			if check.EndTimeRaw != "" {
				t, err := time.Parse(dateTimeNano, check.EndTimeRaw)
				if err != nil {
					return err
				}
				r.Powerwalls[i].UpdateDiagnostic.Checks[j].EndTime = t
			}
		}
	}

	for i, check := range r.Sync.CommissioningDiagnostic.Checks {
		if check.StartTimeRaw != "" {
			t, err := time.Parse(dateTimeNano, check.StartTimeRaw)
			if err != nil {
				return err
			}
			r.Sync.CommissioningDiagnostic.Checks[i].StartTime = t
		}

		if check.EndTimeRaw != "" {
			t, err := time.Parse(dateTimeNano, check.EndTimeRaw)
			if err != nil {
				return err
			}
			r.Sync.CommissioningDiagnostic.Checks[i].EndTime = t
		}
	}

	for i, check := range r.Sync.UpdateDiagnostic.Checks {
		if check.StartTimeRaw != "" {
			t, err := time.Parse(dateTimeNano, check.StartTimeRaw)
			if err != nil {
				return err
			}
			r.Sync.UpdateDiagnostic.Checks[i].StartTime = t
		}

		if check.EndTimeRaw != "" {
			t, err := time.Parse(dateTimeNano, check.EndTimeRaw)
			if err != nil {
				return err
			}
			r.Sync.UpdateDiagnostic.Checks[i].EndTime = t
		}
	}

	return nil
}

// TegSiteInfo defines the response for /api/site_info
type TegSiteInfo struct {
	Timestamp              time.Time
	MeasuredFrequency      float64             `json:"measured_frequency"`
	MaxSystemEnergyKwh     float64             `json:"max_system_energy_kWh"`
	MaxSystemPowerKw       float64             `json:"max_system_power_kW"`
	SiteName               string              `json:"site_name"`
	Timezone               string              `json:"timezone"`
	NetMeterMode           string              `json:"net_meter_mode"`
	MaxSiteMeterPowerKw    int                 `json:"max_site_meter_power_kW"`
	MinSiteMeterPowerKw    int                 `json:"min_site_meter_power_kW"`
	NominalSystemEnergyKwh float64             `json:"nominal_system_energy_kWh"`
	NominalSystemPowerKw   float64             `json:"nominal_system_power_kW"`
	PanelMaxCurrent        int                 `json:"panel_max_current"`
	GridCode               TegSiteInfoGridCode `json:"grid_code"`
}

// TegSiteInfoGridCode defines local grid metadata underneath TegSiteInfo
type TegSiteInfoGridCode struct {
	GridCode           string `json:"grid_code"`
	GridVoltageSetting int    `json:"grid_voltage_setting"`
	GridFreqSetting    int    `json:"grid_freq_setting"`
	GridPhaseSetting   string `json:"grid_phase_setting"`
	Country            string `json:"country"`
	State              string `json:"state"`
	Utility            string `json:"utility"`
}

// TegSitemaster defines the response for /api/sitemaster
type TegSitemaster struct {
	Timestamp        time.Time
	Status           string `json:"status"`
	Running          bool   `json:"running"`
	ConnectedToTesla bool   `json:"connected_to_tesla"`
	PowerSupplyMode  bool   `json:"power_supply_mode"`
	CanReboot        string `json:"can_reboot"`
}

// TegSolars defines the response for /api/solars
type TegSolars struct {
	Timestamp        time.Time
	Brand            string `json:"brand"`
	Model            string `json:"model"`
	PowerRatingWatts int    `json:"power_rating_watts"`
}

// TegStatus defines the response for /api/status
type TegStatus struct {
	Timestamp       time.Time
	GatewayID       string `json:"din"`
	StartTimeRaw    string `json:"start_time"` // 2021-10-26 16:01:02 +0800
	StartTime       time.Time
	UptimeRaw       string `json:"up_time_seconds"` // 89h51m33.77086138s
	Uptime          time.Duration
	IsNew           bool        `json:"is_new"`
	FirmwareVersion string      `json:"version"`
	FirmwareGitHash string      `json:"git_hash"`
	CommissionCount int         `json:"commission_count"`
	DeviceType      string      `json:"device_type"`
	SyncType        string      `json:"sync_type"`
	Leader          interface{} `json:"leader"`
	Followers       interface{} `json:"followers"`
}

// ParseTime converts string times in TegStatus to time.Time
func (r *TegStatus) ParseTime() error {
	if r.StartTimeRaw != "" {
		t, err := time.Parse(dateTimeStatus, r.StartTimeRaw)
		if err != nil {
			return err
		}
		r.StartTime = t
	}

	if r.UptimeRaw != "" {
		d, err := time.ParseDuration(r.UptimeRaw)
		if err != nil {
			return err
		}
		r.Uptime = d
	}

	return nil
}

// TegNetworkConnectionTests defines the response for /api/system/networks/conn_tests
type TegNetworkConnectionTests struct {
	Timestamp  time.Time
	Name       string                      `json:"name"`
	Category   string                      `json:"category"`
	Disruptive bool                        `json:"disruptive"`
	Inputs     interface{}                 `json:"inputs"`
	Checks     []TegNetworkConnectionCheck `json:"checks"`
	Alert      bool                        `json:"alert"`
}

// TegNetworkConnectionCheck defines check data underneath TegNetworkConnectionTests
type TegNetworkConnectionCheck struct {
	Name         string `json:"name"`
	Status       string `json:"status"`
	StartTimeRaw string `json:"start_time"` // dateTimeNano
	StartTime    time.Time
	EndTimeRaw   string `json:"end_time"` // dateTimeNano
	EndTime      time.Time
	Results      interface{} `json:"results"`
	Debug        interface{} `json:"debug"`
	Checks       interface{} `json:"checks"`
}

// ParseTime converts string times in TegNetworkConnectionTests to time.Time
func (r *TegNetworkConnectionTests) ParseTime() error {
	for i, check := range r.Checks {
		if check.StartTimeRaw != "" {
			t, err := time.Parse(dateTimeNano, check.StartTimeRaw)
			if err != nil {
				return err
			}
			r.Checks[i].StartTime = t
		}

		if check.EndTimeRaw != "" {
			t, err := time.Parse(dateTimeNano, check.EndTimeRaw)
			if err != nil {
				return err
			}
			r.Checks[i].EndTime = t
		}
	}

	return nil
}

// TegSystemTesting defines the response for /api/system/testing
type TegSystemTesting struct {
	Timestamp       time.Time
	Running         bool        `json:"running"`
	Status          string      `json:"status"`
	ChargeTests     interface{} `json:"charge_tests"`
	MeterResults    interface{} `json:"meter_results"`
	InverterResults interface{} `json:"inverter_results"`
	Hysteresis      int         `json:"hysteresis"`
	Error           string      `json:"error"`
	Errors          interface{} `json:"errors"`
	Tests           interface{} `json:"tests"`
}

// TegUpdateStatus defines the response for /api/system/update/status
type TegUpdateStatus struct {
	Timestamp               time.Time
	State                   string        `json:"state"`
	Info                    TegUpdateInfo `json:"info"`
	CurrentTime             int           `json:"current_time"`
	LastStatusTime          int           `json:"last_status_time"`
	FirmwareVersion         string        `json:"version"`
	OfflineUpdating         bool          `json:"offline_updating"`
	OfflineUpdateError      string        `json:"offline_update_error"`
	EstimatedBytesPerSecond interface{}   `json:"estimated_bytes_per_second"`
}

// TegUpdateInfo defines update metadata underneath TegUpdateStatus
type TegUpdateInfo struct {
	Status []string `json:"status"`
}

// TegSystemStatus defines the response for /api/system_status
type TegSystemStatus struct {
	Timestamp                       time.Time
	CommandSource                   string            `json:"command_source"`
	BatteryTargetPower              float64           `json:"battery_target_power"`
	BatteryTargetReactivePower      int               `json:"battery_target_reactive_power"`
	NominalFullPackEnergyWattHours  int               `json:"nominal_full_pack_energy"`
	NominalEnergyRemainingWattHours int               `json:"nominal_energy_remaining"`
	MaxPowerEnergyRemaining         int               `json:"max_power_energy_remaining"`
	MaxPowerEnergyToBeCharged       int               `json:"max_power_energy_to_be_charged"`
	MaxChargePowerWatts             int               `json:"max_charge_power"`
	MaxDischargePowerWatts          float64           `json:"max_discharge_power"`
	MaxApparentPower                int               `json:"max_apparent_power"`
	InstantaneousMaxDischargePower  int               `json:"instantaneous_max_discharge_power"`
	InstantaneousMaxChargePower     int               `json:"instantaneous_max_charge_power"`
	GridServicesPower               float64           `json:"grid_services_power"`
	SystemIslandState               string            `json:"system_island_state"`
	AvailableBlocks                 int               `json:"available_blocks"`
	BatteryBlocks                   []TegBatteryBlock `json:"battery_blocks"`
	FfrPowerAvailabilityHigh        float64           `json:"ffr_power_availability_high"`
	FfrPowerAvailabilityLow         float64           `json:"ffr_power_availability_low"`
	LoadChargeConstraint            int               `json:"load_charge_constraint"`
	MaxSustainedRampRate            int               `json:"max_sustained_ramp_rate"`
	GridFaults                      []TegGridFault    `json:"grid_faults"`
	CanReboot                       string            `json:"can_reboot"`
	SmartInvDeltaP                  int               `json:"smart_inv_delta_p"`
	SmartInvDeltaQ                  int               `json:"smart_inv_delta_q"`
	Updating                        bool              `json:"updating"`
	LastToggleTimestampRaw          string            `json:"last_toggle_timestamp"` // dateTimeNano
	LastToggleTimestamp             time.Time
	SolarRealPowerLimit             float64 `json:"solar_real_power_limit"`
	Score                           int     `json:"score"`
	BlocksControlled                int     `json:"blocks_controlled"`
	Primary                         bool    `json:"primary"`
	AuxiliaryLoad                   int     `json:"auxiliary_load"`
	AllEnableLinesHigh              bool    `json:"all_enable_lines_high"`
	InverterNominalUsablePowerWatts int     `json:"inverter_nominal_usable_power"`
	ExpectedEnergyRemaining         int     `json:"expected_energy_remaining"`
}

// TegBatteryBlock defines individual powerwall metadata underneath TegSystemStatus
type TegBatteryBlock struct {
	Type                            string   `json:"Type"`
	PackagePartNumber               string   `json:"PackagePartNumber"`
	PackageSerialNumber             string   `json:"PackageSerialNumber"`
	DisabledReasons                 []string `json:"disabled_reasons"`
	PinvState                       string   `json:"pinv_state"`
	PinvGridState                   string   `json:"pinv_grid_state"`
	NominalEnergyRemainingWattHours int      `json:"nominal_energy_remaining"`
	NominalFullPackEnergy           int      `json:"nominal_full_pack_energy"`
	POut                            float64  `json:"p_out"`
	QOut                            float64  `json:"q_out"`
	VOut                            float64  `json:"v_out"`
	FOut                            float64  `json:"f_out"`
	IOut                            float64  `json:"i_out"`
	EnergyCharged                   int      `json:"energy_charged"`
	EnergyDischarged                int      `json:"energy_discharged"`
	OffGrid                         bool     `json:"off_grid"`
	VfMode                          bool     `json:"vf_mode"`
	WobbleDetected                  bool     `json:"wobble_detected"`
	ChargePowerClamped              bool     `json:"charge_power_clamped"`
	BackupReady                     bool     `json:"backup_ready"`
	OpSeqState                      string   `json:"OpSeqState"`
	Version                         string   `json:"version"`
}

// TegGridFault defines grid fault data underneath TegSystemStatus
type TegGridFault struct {
	Timestamp              int    `json:"timestamp"`
	AlertName              string `json:"alert_name"`
	AlertIsFault           bool   `json:"alert_is_fault"`
	DecodedAlertRaw        string `json:"decoded_alert"`
	DecodedAlert           []TegGridAlert
	AlertRaw               int    `json:"alert_raw"`
	FirmwareGitHash        string `json:"git_hash"`
	SiteUID                string `json:"site_uid"`
	EcuType                string `json:"eco_type"`
	EcuPackagePartNumber   string `json:"ecu_package_part_number"`
	EcuPackageSerialNumber string `json:"ecu_package_serial_number"`
}

// TegGridAlert defines grid alert data underneath TegGridFault
type TegGridAlert struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"` // Observed as string or float64
	Units string      `json:"units"`
}

// ParseTime converts string times in TegSystemStatus to time.Time
func (r *TegSystemStatus) ParseTime() error {
	if r.LastToggleTimestampRaw != "" {
		t, err := time.Parse(dateTimeNano, r.LastToggleTimestampRaw)
		if err != nil {
			return err
		}
		r.LastToggleTimestamp = t
	}

	return nil
}

// ParseFaults unwraps JSON fault data
func (r *TegSystemStatus) ParseFaults() error {
	alert := &[]TegGridAlert{}
	for i, fault := range r.GridFaults {
		if fault.DecodedAlertRaw != "" {
			err := json.Unmarshal([]byte(fault.DecodedAlertRaw), alert)
			if err != nil {
				return err
			}
			r.GridFaults[i].DecodedAlert = *alert
		}
	}
	return nil
}

// TegSystemGridStatus defines the response for /api/system_status/grid_status
type TegSystemGridStatus struct {
	Timestamp          time.Time
	GridStatus         string `json:"grid_status"`
	GridServicesActive bool   `json:"grid_services_active"`
}

// TegSystemStateOfEnergy defines the response for /api/system_status/soe
type TegSystemStateOfEnergy struct {
	Timestamp  time.Time
	Percentage float64 `json:"percentage"`
}

