package model

// Payload for /api/login/Basic
type AuthPayload struct {
	Username     string
	Password     string
	Email        string
	Force_Sm_Off bool
}

// Response for /api/login/Basic
type AuthResponse struct {
	Email     string   `json:"email"`
	Firstname string   `json:"firstname"`
	Lastname  string   `json:"lastname"`
	Roles     []string `json:"roles"`
	Token     string   `json:"token"`
	Provider  string   `json:"provider"`
	loginTime string   `json:"loginTime"`
}

// Response for /api/status
type TegStatus struct {
	GatewayId       string      `json:"din"`
	StartTime       string      `json:"start_time"`
	Uptime          string      `json:"up_time_seconds"`
	IsNew           bool        `json:"is_new"`
	FirmwareVersion string      `json:"version"`
	FirmwareGitHash string      `json:"git_hash"`
	CommissionCount int         `json:"commission_count"`
	DeviceType      string      `json:"device_type"`
	SyncType        string      `json:"sync_type"`
	Leader          interface{} `json:"leader"`
	Followers       interface{} `json:"followers"`
}

// Response for /api/meters/aggregates
type TegMeters struct {
	Site    TegMetersAggregate `json:"site"`
	Battery TegMetersAggregate `json:"battery"`
	Load    TegMetersAggregate `json:"load"`
	Solar   TegMetersAggregate `json:"solar"`
}

type TegMetersAggregate struct {
	LastCommunicationTime             string  `json:"last_communication_time"`
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

// Response for /api/operation
type TegOperation struct {
	RealMode             string  `json:"real_mode"`
	BackupReservePercent float64 `json:"backup_reserve_percent"`
}

// Response for /api/powerwalls
type TegPowerwalls struct {
	Enumerating                bool              `json:"enumerating"`
	Updating                   bool              `json:"updating"`
	CheckingIfOffgrid          bool              `json:"checking_if_offgrid"`
	RunningPhaseDetection      bool              `json:"running_phase_detection"`
	PhaseDetectionLastError    string            `json:"phase_detection_last_error"`
	BubbleShedding             bool              `json:"bubble_shedding"`
	OnGridCheckError           string            `json:"on_grid_check_error"`
	GridQualifying             bool              `json:"grid_qualifying"`
	GridCodeValidating         bool              `json:"grid_code_validating"`
	PhaseDetectionNotAvailable bool              `json:"phase_detection_not_available"`
	Powerwalls                 []TegPowerwall    `json:"powerwalls"`
	GatewayId                  string            `json:"gateway_din"`
	Sync                       TegPowerwallsSync `json:"sync"`
	Msa                        interface{}       `json:"msa"`
	States                     interface{}       `json:"states"`
}

type TegPowerwall struct {
	Type                        string                  `json:"Type"`
	PackagePartNumber           string                  `json:"PackagePartNumber"`
	PackageSerialNumber         string                  `json:"PackageSerialNumber"`
	Subtype                     string                  `json:"type"`
	GridState                   string                  `json:"grid_state"`
	GridReconnectionTimeSeconds int                     `json:"grid_reconnection_time_seconds"`
	UnderPhaseDetection         bool                    `json:"under_phase_detection"`
	Updating                    bool                    `json:"updating"`
	CommissioningDiagnostic     TegPowerwallsDiagnostic `json:"commissioning_diagnostic"`
	UpdateDiagnostic            TegPowerwallsDiagnostic `json:"update_diagnostic"`
	BcType                      interface{}             `json:"bc_type"`
	InConfig                    bool                    `json:"in_config"`
}

type TegPowerwallsSync struct {
	Updating                bool                    `json:"updating"`
	CommissioningDiagnostic TegPowerwallsDiagnostic `json:"commissioning_diagnostic"`
	UpdateDiagnostic        TegPowerwallsDiagnostic `json:"update_diagnostic"`
}

type TegPowerwallsDiagnostic struct {
	Name       string               `json:"name"`
	Category   string               `json:"category"`
	Disruptive bool                 `json:"disruptive"`
	Inputs     interface{}          `json:"inputs"`
	Checks     []TegPowerwallsCheck `json:"checks"`
	Alert      bool                 `json:"alert"`
}

type TegPowerwallsCheck struct {
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	StartTime string      `json:"start_time"`
	EndTime   string      `json:"end_time"`
	Message   string      `json:"message"`
	Progress  int         `json:"progress"`
	Results   interface{} `json:"results"`
	Debug     interface{} `json:"debug"`
	Checks    interface{} `json:"checks"`
}

// Response for /api/site_info
type TegSiteInfo struct {
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

type TegSiteInfoGridCode struct {
	GridCode           string `json:"grid_code"`
	GridVoltageSetting int    `json:"grid_voltage_setting"`
	GridFreqSetting    int    `json:"grid_freq_setting"`
	GridPhaseSetting   string `json:"grid_phase_setting"`
	Country            string `json:"country"`
	State              string `json:"state"`
	Utility            string `json:"utility"`
}

// Response for /api/solars
type TegSolars struct {
	Brand            string `json:"brand"`
	Model            string `json:"model"`
	PowerRatingWatts int    `json:"power_rating_watts"`
}

// Response for /api/system/networks/conn_tests
type TegNetworkConnectionTests struct {
	Name       string                      `json:"name"`
	Category   string                      `json:"category"`
	Disruptive bool                        `json:"disruptive"`
	Inputs     interface{}                 `json:"inputs"`
	Checks     []TegNetworkConnectionCheck `json:"checks"`
	Alert      bool                        `json:"alert"`
}

type TegNetworkConnectionCheck struct {
	Name      string      `json:"name"`
	Status    string      `json:"status"`
	StartTime string      `json:"start_time"`
	EndTime   string      `json:"end_time"`
	Results   interface{} `json:"results"`
	Debug     interface{} `json:"debug"`
	Checks    interface{} `json:"checks"`
}

// Response for /api/system/testing
type TegSystemTesting struct {
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

// Response for /api/system/update/status
type TegUpdateStatus struct {
	State                   string        `json:"state"`
	Info                    TegUpdateInfo `json:"info"`
	CurrentTime             int           `json:"current_time"`
	LastStatusTime          int           `json:"last_status_time"`
	FirmwareVersion         string        `json:"version"`
	OfflineUpdating         bool          `json:"offline_updating"`
	OfflineUpdateError      string        `json:"offline_update_error"`
	EstimatedBytesPerSecond interface{}   `json:"estimated_bytes_per_second"`
}

type TegUpdateInfo struct {
	Status []string `json:"status"`
}

// Response for /api/system_status
type TegSystemStatus struct {
	CommandSource                   string            `json:"command_source"`
	BatteryTargetPower              float64           `json:"battery_target_power"`
	BatteryTargetReactivePower      int               `json:"battery_target_reactive_power"`
	NominalFullPackEnergyWattHours  int               `json:"nominal_full_pack_energy"`
	NominalEnergyRemainingWattHours int               `json:"nominal_energy_remaining"`
	MaxPowerEnergyRemaining         int               `json:"max_power_energy_remaining"`
	MaxPowerEnergyToBeCharged       int               `json:"max_power_energy_to_be_charged"`
	MaxChargePowerWatts             int               `json:"max_charge_power"`
	MaxDischargePowerWatts          int               `json:"max_discharge_power"`
	MaxApparentPower                int               `json:"max_apparent_power"`
	InstantaneousMaxDischargePower  int               `json:"instantaneous_max_discharge_power"`
	InstantaneousMaxChargePower     int               `json:"instantaneous_max_charge_power"`
	GridServicesPower               int               `json:"grid_services_power"`
	SystemIslandState               string            `json:"system_island_state"`
	AvailableBlocks                 int               `json:"available_blocks"`
	BatteryBlocks                   []TegBatteryBlock `json:"battery_blocks"`
	FfrPowerAvailabilityHigh        int               `json:"ffr_power_availability_high"`
	FfrPowerAvailabilityLow         int               `json:"ffr_power_availability_low"`
	LoadChargeConstraint            int               `json:"load_charge_constraint"`
	MaxSustainedRampRate            int               `json:"max_sustained_ramp_rate"`
	GridFaults                      []TegGridFault    `json:"grid_faults"`
	CanReboot                       string            `json:"can_reboot"`
	SmartInvDeltaP                  int               `json:"smart_inv_delta_p"`
	SmartInvDeltaQ                  int               `json:"smart_inv_delta_q"`
	LastToggleTimestamp             string            `json:"last_toggle_timestamp"`
	SolarRealPowerLimit             float64           `json:"solar_real_power_limit"`
	Score                           int               `json:"score"`
	BlocksControlled                int               `json:"blocks_controlled"`
	Primary                         bool              `json:"primary"`
	AuxiliaryLoad                   int               `json:"auxiliary_load"`
	AllEnableLinesHigh              bool              `json:"all_enable_lines_high"`
	InverterNominalUsablePowerWatts int               `json:"inverter_nominal_usable_power"`
	ExpectedEnergyRemaining         int               `json:"expected_energy_remaining"`
}

type TegBatteryBlock struct {
	Type                            string      `json:"Type"`
	PackagePartNumber               string      `json:"PackagePartNumber"`
	PackageSerialNumber             string      `json:"PackageSerialNumber"`
	DisabledReasons                 interface{} `json:"disabled_reasons"`
	PinvState                       string      `json:"pinv_state"`
	PinvGridState                   string      `json:"pinv_grid_state"`
	NominalEnergyRemainingWattHours int         `json:"nominal_energy_remaining"`
	NominalFullPackEnergy           int         `json:"nominal_full_pack_energy"`
	POut                            int         `json:"p_out"`
	QOut                            int         `json:"q_out"`
	VOut                            float64     `json:"v_out"`
	FOut                            float64     `json:"f_out"`
	IOut                            float64     `json:"i_out"`
	EnergyCharged                   int         `json:"energy_charged"`
	EnergyDischarged                int         `json:"energy_discharged"`
	OffGrid                         bool        `json:"off_grid"`
	VfMode                          bool        `json:"vf_mode"`
	WobbleDetected                  bool        `json:"wobble_detected"`
	ChargePowerClamped              bool        `json:"charge_power_clamped"`
	BackupReady                     bool        `json:"backup_ready"`
	OpSeqState                      string      `json:"OpSeqState"`
	Version                         string      `json:"version"`
}

type TegGridFault struct {
	Timestamp              int        `json:"timestamp"`
	AlertName              string     `json:"alert_name"`
	AlertIsFault           bool       `json:"alert_is_fault"`
	DecodedAlert           []TegAlert `json:"decoded_alert"`
	AlertRaw               int        `json:"alert_raw"`
	FirmwareGitHash        string     `json:"git_hash"`
	SiteUid                string     `json:"site_uid"`
	EcuType                string     `json:"eco_type"`
	EcuPackagePartNumber   string     `json:"ecu_package_part_number"`
	EcuPackageSerialNumber string     `json:"ecu_package_serial_number"`
}

type TegAlert struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"` // Observed as string or float64
	Units string      `json:"units"`
}

// Response for /api/system_status/grid_status
type TegSystemGridStatus struct {
	GridStatus         string `json:"grid_status"`
	GridServicesActive bool   `json:"grid_services_active"`
}

// Response for /api/system_status/soe
type TegSystemStateOfEnergy struct {
	Percentage float64 `json:"percentage"`
}
