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
	i := &Info{}

	err := v.do("miIO.info", nil, i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

type WIFIStatus struct {
	State         string `json:"state"`
	AuthFailCount int    `json:"auth_fail_count"`
	ConnFailCount int    `json:"conn_fail_count"`
	DHCPFailCount int    `json:"dhcp_fail_count"`
}

func (v *Vacuum) WIFIStatus() (*WIFIStatus, error) {
	w := &WIFIStatus{}

	err := v.do("miIO.wifi_assoc_state", nil, w)
	if err != nil {
		return nil, err
	}

	return w, nil
}
