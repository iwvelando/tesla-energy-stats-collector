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

// Response for /api/meters/status
type TegMetersStatus struct {
	Status string      `json:"status"`
	Errors interface{} `json:"errors"`
	Serial string      `json:"serial"`
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

// Response for /api/sitemaster
type TegSitemaster struct {
	Status           string `json:"status"`
	Running          bool   `json:"running"`
	ConnectedToTesla bool   `json:"connected_to_tesla"`
	PowerSupplyMode  bool   `json:"power_supply_mode"`
	CanReboot        string `json:"can_reboot"`
}

// Response for /api/solar_powerwall
type TegSolarPowerwall struct {
	PvacStatus          TegPvacStatus `json:"pvac_status"`
	PvsStatus           TegPvsStatus  `json:"pvs_status"`
	PvPowerLimit        int           `json:"pv_power_limit"`
	PowerStatusSetpoint string        `json:"power_status_setpoint"`
	PvacAlerts          TegPvacAlerts `json:"pvac_alerts"`
	PvsAlerts           TegPvsAlerts  `json:"pvs_alerts"`
}

type TegPvacStatus struct {
	State           string            `json:"state"`
	Disabled        bool              `json:"disabled"`
	DisabledReasons []string          `json:"disabled_reasons"`
	GridState       string            `json:"grid_state"`
	InvState        string            `json:"inv_state"`
	VOut            float64           `json:"v_out"`
	FOut            float64           `json:"f_out"`
	POut            int               `json:"p_out"`
	QOut            int               `json:"q_out"`
	IOut            float64           `json:"i_out"`
	StringVitals    []TegStringVitals `json:"string_vitals"`
}

type TegStringVitals struct {
	StringId        int     `json:"string_id"`
	Connected       bool    `json:"connected"`
	MeasuredVoltage float64 `json:"measured_voltage"`
	Current         float64 `json:"current"`
	MeasuredPower   int     `json:"measured_power"`
}

type TegPvsStatus struct {
	State         string  `json:"state"`
	Disabled      bool    `json:"disabled"`
	EnableOutput  bool    `json:"enable_output"`
	Vll           float64 `json:"v_ll"`
	SelfTestState string  `json:"self_test_state"`
}

type TegPvacAlerts struct {
	LastRxTime                          string `json:"LastRxTime"`
	ReceivedMuxBitmask                  int    `json:"ReceivedMuxBitmask"`
	PVACAlertMatrixIndex                int    `json:"PVAC_alertMatrixIndex"`
	PVACA001InvL1HWOvercurrent          bool   `json:"PVAC_a001_inv_L1_HW_overcurrent"`
	PVACA002InvL2HWOvercurrent          bool   `json:"PVAC_a002_inv_L2_HW_overcurrent"`
	PVACA003InvHVBusHWOvervoltage       bool   `json:"PVAC_a003_inv_HVBus_HW_overvoltage"`
	PVACA004PvHWCMPSSOCSTGA             bool   `json:"PVAC_a004_pv_HW_CMPSS_OC_STGA"`
	PVACA005PvHWCMPSSOCSTGB             bool   `json:"PVAC_a005_pv_HW_CMPSS_OC_STGB"`
	PVACA006PvHWCMPSSOCSTGC             bool   `json:"PVAC_a006_pv_HW_CMPSS_OC_STGC"`
	PVACA007PvHWCMPSSOCSTGD             bool   `json:"PVAC_a007_pv_HW_CMPSS_OC_STGD"`
	PVACA008InvHVBusUndervoltage        bool   `json:"PVAC_a008_inv_HVBus_undervoltage"`
	PVACA009SwAppBoot                   bool   `json:"PVAC_a009_SwAppBoot"`
	PVACA010InvACOvervoltage            bool   `json:"PVAC_a010_inv_AC_overvoltage"`
	PVACA011InvACUndervoltage           bool   `json:"PVAC_a011_inv_AC_undervoltage"`
	PVACA012InvACOverfrequency          bool   `json:"PVAC_a012_inv_AC_overfrequency"`
	PVACA013InvACUnderfrequency         bool   `json:"PVAC_a013_inv_AC_underfrequency"`
	PVACA014PVSDisabledRelay            bool   `json:"PVAC_a014_PVS_disabled_relay"`
	PVACA015PvHWAllegroOCSTGA           bool   `json:"PVAC_a015_pv_HW_Allegro_OC_STGA"`
	PVACA016PvHWAllegroOCSTGB           bool   `json:"PVAC_a016_pv_HW_Allegro_OC_STGB"`
	PVACA017PvHWAllegroOCSTGC           bool   `json:"PVAC_a017_pv_HW_Allegro_OC_STGC"`
	PVACA018PvHWAllegroOCSTGD           bool   `json:"PVAC_a018_pv_HW_Allegro_OC_STGD"`
	PVACA019AmbientOvertemperature      bool   `json:"PVAC_a019_ambient_overtemperature"`
	PVACA020DspOvertemperature          bool   `json:"PVAC_a020_dsp_overtemperature"`
	PVACA021DcacHeatsinkOvertemperature bool   `json:"PVAC_a021_dcac_heatsink_overtemperature"`
	PVACA022MpptHeatsinkOvertemperature bool   `json:"PVAC_a022_mppt_heatsink_overtemperature"`
	PVACA023Unused                      bool   `json:"PVAC_a023_unused"`
	PVACA024PVACrxCommandMia            bool   `json:"PVAC_a024_PVACrx_Command_mia"`
	PVACA025PVSStatusMia                bool   `json:"PVAC_a025_PVS_Status_mia"`
	PVACA026InvACPeakOvervoltage        bool   `json:"PVAC_a026_inv_AC_peak_overvoltage"`
	PVACA027InvK1RelayWelded            bool   `json:"PVAC_a027_inv_K1_relay_welded"`
	PVACA028InvK2RelayWelded            bool   `json:"PVAC_a028_inv_K2_relay_welded"`
	PVACA029PumpFaulted                 bool   `json:"PVAC_a029_pump_faulted"`
	PVACA030FanFaulted                  bool   `json:"PVAC_a030_fan_faulted"`
	PVACA031VFCheckOV                   bool   `json:"PVAC_a031_VFCheck_OV"`
	PVACA032VFCheckUV                   bool   `json:"PVAC_a032_VFCheck_UV"`
	PVACA033VFCheckOF                   bool   `json:"PVAC_a033_VFCheck_OF"`
	PVACA034VFCheckUF                   bool   `json:"PVAC_a034_VFCheck_UF"`
	PVACA035VFCheckRoCoF                bool   `json:"PVAC_a035_VFCheck_RoCoF"`
	PVACA036InvLostILControl            bool   `json:"PVAC_a036_inv_lost_iL_control"`
	PVACA037PVSProcessorNERROR          bool   `json:"PVAC_a037_PVS_processor_nERROR"`
	PVACA038InvFailedXcapPrecharge      bool   `json:"PVAC_a038_inv_failed_xcap_precharge"`
	PVACA039InvHVBusSWOvervoltage       bool   `json:"PVAC_a039_inv_HVBus_SW_overvoltage"`
	PVACA040PumpCorrectionSaturated     bool   `json:"PVAC_a040_pump_correction_saturated"`
	PVACA041ExcessPVClampTriggered      bool   `json:"PVAC_a041_excess_PV_clamp_triggered"`
}

type TegPvsAlerts struct {
	LastRxTime                        string `json:"LastRxTime"`
	ReceivedMuxBitmask                int    `json:"ReceivedMuxBitmask"`
	PVSA001WatchdogReset              bool   `json:"PVS_a001_WatchdogReset"`
	PVSA002SWAppBoot                  bool   `json:"PVS_a002_SW_App_Boot"`
	PVSA003V12vOutOfBounds            bool   `json:"PVS_a003_V12vOutOfBounds"`
	PVSA004V1v5OutOfBounds            bool   `json:"PVS_a004_V1v5OutOfBounds"`
	PVSA005VAfdRefOutOfBounds         bool   `json:"PVS_a005_VAfdRefOutOfBounds"`
	PVSA006GfOvercurrent300           bool   `json:"PVS_a006_GfOvercurrent300"`
	PVSA007UNUSED7                    bool   `json:"PVS_a007_UNUSED_7"`
	PVSA008UNUSED8                    bool   `json:"PVS_a008_UNUSED_8"`
	PVSA009GfOvercurrent030           bool   `json:"PVS_a009_GfOvercurrent030"`
	PVSA010PvIsolationTotal           bool   `json:"PVS_a010_PvIsolationTotal"`
	PVSA011PvIsolationStringA         bool   `json:"PVS_a011_PvIsolationStringA"`
	PVSA012PvIsolationStringB         bool   `json:"PVS_a012_PvIsolationStringB"`
	PVSA013PvIsolationStringC         bool   `json:"PVS_a013_PvIsolationStringC"`
	PVSA014PvIsolationStringD         bool   `json:"PVS_a014_PvIsolationStringD"`
	PVSA015SelfTestGroundFault        bool   `json:"PVS_a015_SelfTestGroundFault"`
	PVSA016ESMFault                   bool   `json:"PVS_a016_ESMFault"`
	PVSA017MciStringA                 bool   `json:"PVS_a017_MciStringA"`
	PVSA018MciStringB                 bool   `json:"PVS_a018_MciStringB"`
	PVSA019MciStringC                 bool   `json:"PVS_a019_MciStringC"`
	PVSA020MciStringD                 bool   `json:"PVS_a020_MciStringD"`
	PVSA021RapidShutdown              bool   `json:"PVS_a021_RapidShutdown"`
	PVSA022Mci1SignalLevel            bool   `json:"PVS_a022_Mci1SignalLevel"`
	PVSA023Mci2SignalLevel            bool   `json:"PVS_a023_Mci2SignalLevel"`
	PVSA024Mci3SignalLevel            bool   `json:"PVS_a024_Mci3SignalLevel"`
	PVSA025Mci4SignalLevel            bool   `json:"PVS_a025_Mci4SignalLevel"`
	PVSA026Mci1PvVoltage              bool   `json:"PVS_a026_Mci1PvVoltage"`
	PVSA027Mci2PvVoltage              bool   `json:"PVS_a027_Mci2PvVoltage"`
	PVSA028SystemInitFailed           bool   `json:"PVS_a028_systemInitFailed"`
	PVSA029PvArcFault                 bool   `json:"PVS_a029_PvArcFault"`
	PVSA030VDcOv                      bool   `json:"PVS_a030_VDcOv"`
	PVSA031Mci3PvVoltage              bool   `json:"PVS_a031_Mci3PvVoltage"`
	PVSA032Mci4PvVoltage              bool   `json:"PVS_a032_Mci4PvVoltage"`
	PVSA033DataException              bool   `json:"PVS_a033_dataException"`
	PVSA034PeImpedance                bool   `json:"PVS_a034_PeImpedance"`
	PVSA035PvArcDetected              bool   `json:"PVS_a035_PvArcDetected"`
	PVSA036PvArcLockout               bool   `json:"PVS_a036_PvArcLockout"`
	PVSA037PvArcFault2                bool   `json:"PVS_a037_PvArcFault2"`
	PVSA038PvArcFaultSelfTest         bool   `json:"PVS_a038_PvArcFault_SelfTest"`
	PVSA039SelfTestRelayFault         bool   `json:"PVS_a039_SelfTestRelayFault"`
	PVSA040LEDIrrationalFault         bool   `json:"PVS_a040_LEDIrrationalFault"`
	PVSA041MciPowerSwitch             bool   `json:"PVS_a041_MciPowerSwitch"`
	PVSA042MciPowerFault              bool   `json:"PVS_a042_MciPowerFault"`
	PVSA043LockedPvStringSafety       bool   `json:"PVS_a043_LockedPvStringSafety"`
	PVSA044FaultStatePvStringSafety   bool   `json:"PVS_a044_FaultStatePvStringSafety"`
	PVSA045RelayCoilIrrationalFault   bool   `json:"PVS_a045_RelayCoilIrrationalFault"`
	PVSA046RelayCoilIrrationalLockout bool   `json:"PVS_a046_RelayCoilIrrationalLockout"`
	PVSA047AcSensorIrrationalFault    bool   `json:"PVS_a047_AcSensorIrrationalFault"`
	PVSA048DcSensorIrrationalFault    bool   `json:"PVS_a048_DcSensorIrrationalFault"`
	PVSA049ArcSignalMibspiHealth      bool   `json:"PVS_a049_arcSignalMibspiHealth"`
	PVSA050RelayCoilIrrationalWarning bool   `json:"PVS_a050_RelayCoilIrrationalWarning"`
	PVSA051DcBusShortCircuitDetected  bool   `json:"PVS_a051_DcBusShortCircuitDetected"`
	PVSA052PvArcFaultPreSelfTest      bool   `json:"PVS_a052_PvArcFault_PreSelfTest"`
}

// Response for /api/solars
type TegSolars struct {
	Brand            string `json:"brand"`
	Model            string `json:"model"`
	PowerRatingWatts int    `json:"power_rating_watts"`
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
