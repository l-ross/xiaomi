package vacuum

type Info struct {
	HardwareVersion string `json:"hw_ver"`
	FirmwareVersion string `json:"fw_ver"`
	AccessPoint     struct {
		SSID  string `json:"ssid"`
		BSSID string `json:"bssid"`
		RSSI  int    `json:"rssi"`
	} `json:"ap"`
	Network struct {
		LocalIP string `json:"localIp"`
		Mask    string `json:"mask"`
		Gateway string `json:"gw"`
	} `json:"netif"`
	Model string `json:"model"`
	MAC   string `json:"mac"`
	Token string `json:"token"`
	Life  int    `json:"life"`
}

// Info retrieves the hardware and network info of the Vacuum.
func (v *Vacuum) Info() (*Info, error) {
	rsp := &Info{}

	err := v.do("miIO.info", nil, rsp)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

type WIFIStatus struct {
	State         string `json:"state"`
	AuthFailCount int    `json:"auth_fail_count"`
	ConnFailCount int    `json:"conn_fail_count"`
	DHCPFailCount int    `json:"dhcp_fail_count"`
}

func (v *Vacuum) WIFIStatus() (*WIFIStatus, error) {
	rsp := &WIFIStatus{}

	err := v.do("miIO.wifi_assoc_state", nil, rsp)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

type Locale struct {
	Name       string `json:"name"`
	BOM        string `json:"bom"`
	Location   string `json:"location"`
	Language   string `json:"language"`
	WIFIPlan   string `json:"wifiplan"`
	TimeZone   string `json:"timezone"`
	LogServer  string `json:"logserver"`
	FeatureSet int    `json:"featureset"`
}

func (v *Vacuum) Locale() (*Locale, error) {
	rsp := make([]*Locale, 0)

	err := v.do("app_get_locale", nil, &rsp)
	if err != nil {
		return nil, err
	}

	if len(rsp) != 1 {
		return nil, ErrUnexpectedResponse
	}

	return rsp[0], nil
}

// FirmwareFeatures retrieves the enabled features of this vacuum.
//
// The known features are available as FeatureCode constants.
func (v *Vacuum) FirmwareFeatures() ([]int, error) {
	ret := make([]int, 0)

	err := v.do("get_fw_features", nil, &ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type Status struct {
	MessageVersion  int `json:"msg_ver"`
	MessageSequence int `json:"msg_seq"`
	// See StatusCode for the known vacuum states.
	State int `json:"state"`
	// See ErrorCode for the known vacuum errors.
	ErrorCode int `json:"error_code"`

	Battery                int `json:"battery"`
	CleanArea              int `json:"clean_area"`
	CleanTime              int `json:"clean_time"`
	DNDEnabled             int `json:"dnd_enabled"`
	FanPower               int `json:"fan_power"`
	InCleaning             int `json:"in_cleaning"`
	InFreshState           int `json:"in_fresh_state"`
	InReturning            int `json:"in_returning"`
	IsLocating             int `json:"is_locating"`
	LabStatus              int `json:"lab_status"`
	LockStatus             int `json:"lock_status"`
	MapPresent             int `json:"map_present"`
	MapStatus              int `json:"map_status"`
	MopForbiddenEnable     int `json:"mop_forbidden_enable"`
	WaterBoxCarriageStatus int `json:"water_box_carriage_status"`
	WaterBoxMode           int `json:"water_box_mode"`
}

// Status retrieve the over state of the Vacuum.
func (v *Vacuum) Status() (*Status, error) {
	rsp := make([]*Status, 0)

	err := v.do("get_status", nil, &rsp)
	if err != nil {
		return nil, err
	}

	if len(rsp) != 1 {
		return nil, ErrUnexpectedResponse
	}

	return rsp[0], nil
}

type InitialStatus struct {
	Locale  Locale `json:"local_info"`
	Feature []int  `json:"feature_info"`
	Status  Status `json:"status_info"`
}

func (v *Vacuum) InitialStatus() (*InitialStatus, error) {
	rsp := make([]*InitialStatus, 0)

	err := v.do("app_get_init_status", nil, &rsp)
	if err != nil {
		return nil, err
	}

	if len(rsp) != 1 {
		return nil, ErrUnexpectedResponse
	}

	return rsp[0], nil
}

type NetworkInfo struct {
	SSID  string `json:"ssid"`
	IP    string `json:"ip"`
	MAC   string `json:"mac"`
	BSSID string `json:"bssid"`
	RSSI  int    `json:"rssi"`
}

func (v *Vacuum) NetworkInfo() (*NetworkInfo, error) {
	rsp := &NetworkInfo{}

	err := v.do("get_network_info", nil, rsp)
	if err != nil {
		return nil, err
	}

	return rsp, nil
}

type serialNumber struct {
	SerialNumber string `json:"serial_number"`
}

func (v *Vacuum) SerialNumber() (string, error) {
	rsp := make([]*serialNumber, 0)

	err := v.do("get_serial_number", nil, &rsp)
	if err != nil {
		return "", err
	}

	if len(rsp) != 1 {
		return "", ErrUnexpectedResponse
	}

	return rsp[0].SerialNumber, nil
}
