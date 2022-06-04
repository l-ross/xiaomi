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

func (v *Vacuum) Info() (*Info, error) {
	ret := &Info{}

	err := v.do("miIO.info", nil, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type WIFIStatus struct {
	State         string `json:"state"`
	AuthFailCount int    `json:"auth_fail_count"`
	ConnFailCount int    `json:"conn_fail_count"`
	DHCPFailCount int    `json:"dhcp_fail_count"`
}

func (v *Vacuum) WIFIStatus() (*WIFIStatus, error) {
	ret := &WIFIStatus{}

	err := v.do("miIO.wifi_assoc_state", nil, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
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
	ret := make([]*Locale, 0)

	err := v.do("app_get_locale", nil, &ret)
	if err != nil {
		return nil, err
	}

	if len(ret) != 1 {
		return nil, ErrUnexpectedResponse
	}

	return ret[0], nil
}

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
	State           int `json:"state"`
	ErrorCode       int `json:"error_code"`

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

func (v *Vacuum) Status() (*Status, error) {
	ret := make([]*Status, 0)

	err := v.do("get_status", nil, &ret)
	if err != nil {
		return nil, err
	}

	if len(ret) != 1 {
		return nil, ErrUnexpectedResponse
	}

	return ret[0], nil
}

type InitialStatus struct {
	Locale  Locale `json:"local_info"`
	Feature []int  `json:"feature_info"`
	Status  Status `json:"status_info"`
}

func (v *Vacuum) InitialStatus() (*InitialStatus, error) {
	ret := make([]*InitialStatus, 0)

	err := v.do("app_get_init_status", nil, &ret)
	if err != nil {
		return nil, err
	}

	if len(ret) != 1 {
		return nil, ErrUnexpectedResponse
	}

	return ret[0], nil
}

type NetworkInfo struct {
	SSID  string `json:"ssid"`
	IP    string `json:"ip"`
	MAC   string `json:"mac"`
	BSSID string `json:"bssid"`
	RSSI  int    `json:"rssi"`
}

func (v *Vacuum) NetworkInfo() (*NetworkInfo, error) {
	ret := &NetworkInfo{}

	err := v.do("get_network_info", nil, ret)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

type serialNumber struct {
	SerialNumber string `json:"serial_number"`
}

func (v *Vacuum) SerialNumber() (string, error) {
	ret := make([]*serialNumber, 0)

	err := v.do("get_serial_number", nil, &ret)
	if err != nil {
		return "", err
	}

	if len(ret) != 1 {
		return "", ErrUnexpectedResponse
	}

	return ret[0].SerialNumber, nil
}
