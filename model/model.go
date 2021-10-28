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
	Id              string   `json:"din"`
	StartTime       string   `json:"start_time"`
	Uptime          string   `json:"up_time_seconds"`
	IsNew           bool     `json:"is_new"`
	FirmwareVersion string   `json:"version"`
	FirmwareGitHash string   `json:"git_hash"`
	CommissionCount int      `json:"commission_count"`
	DeviceType      string   `json:"device_type"`
	SyncType        string   `json:"sync_type"`
	Leader          string   `json:"leader"`
	Followers       []string `json:"followers"`
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
