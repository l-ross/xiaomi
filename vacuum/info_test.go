package vacuum

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVacuum_Info(t *testing.T) {
	t.Parallel()

	v := &Vacuum{
		client: &mockClient{
			rsp: []byte(`
				{
					"partner_id": "",
					"id": 7840,
					"code": 0,
					"message": "ok",
					"result": {
						"hw_ver": "Linux",
						"fw_ver": "3.3.6_003061",
						"ap": {
							"ssid": "mylocalwifi",
							"bssid": "F2:9F:C2:AA:AA:AA",
							"rssi": -63
						},
						"netif": {
							"localIp": "192.168.9.199",
							"mask": "255.255.255.0",
							"gw": "192.168.9.1"
						},
						"model": "rockrobo.vacuum.v1",
						"mac": "28:6C:07:AA:AA:AA",
						"token": "4474746a612345678905612345678900",
						"life": 62848
					}
				}`,
			),
		},
	}

	got, err := v.Info()
	require.NoError(t, err)
	assert.Equal(t, "Linux", got.HardwareVersion)
}

func TestVacuum_WIFIStatus(t *testing.T) {
	t.Parallel()

	v := &Vacuum{
		client: &mockClient{
			rsp: []byte(`
				{
					"id": 37,
					"code": 0,
					"message": "ok",
					"result": {
						"state": "ONLINE",
						"auth_fail_count": 0,
						"conn_success_count": 1,
						"conn_fail_count": 0,
						"dhcp_fail_count": 0
					}
				}`,
			),
		},
	}

	got, err := v.WIFIStatus()
	require.NoError(t, err)
	assert.Equal(t, "ONLINE", got.State)
}

func TestVacuum_Locale(t *testing.T) {
	t.Parallel()

	v := &Vacuum{
		client: &mockClient{
			rsp: []byte(`
				{
					"result": [{
							"name": "custom_A.03.0070_CE",
							"bom": "A.03.0070",
							"location": "de",
							"language": "en",
							"wifiplan": "",
							"timezone": "Europe/Berlin",
							"logserver": "awsde0.fds.api.xiaomi.com",
							"featureset": 1
						}
					],
					"id": 5879
				}`,
			),
		},
	}

	got, err := v.Locale()
	require.NoError(t, err)
	assert.Equal(t, "custom_A.03.0070_CE", got.Name)
}

func TestVacuum_FirmwareFeatures(t *testing.T) {
	t.Parallel()

	v := &Vacuum{
		client: &mockClient{
			rsp: []byte(`
				{
					"result": [111, 112, 113, 114, 115, 116, 117, 118, 119, 122, 125],
					"id": 7177
				}`,
			),
		},
	}

	got, err := v.FirmwareFeatures()
	require.NoError(t, err)
	require.Len(t, got, 11)
	assert.Equal(t, got[0], 111)
}

func TestVacuum_Status(t *testing.T) {
	t.Parallel()

	v := &Vacuum{
		client: &mockClient{
			rsp: []byte(`
				{
					"result": [{
							"msg_ver": 2,
							"msg_seq": 52,
							"state": 8,
							"battery": 100,
							"clean_time": 15,
							"clean_area": 140000,
							"error_code": 0,
							"map_present": 1,
							"in_cleaning": 0,
							"in_returning": 0,
							"in_fresh_state": 1,
							"lab_status": 1,
							"water_box_status": 1,
							"fan_power": 102,
							"dnd_enabled": 0,
							"map_status": 3,
							"is_locating": 0,
							"lock_status": 0,
							"water_box_mode": 204,
							"water_box_carriage_status": 0,
							"mop_forbidden_enable": 0
						}
					],
					"id": 96
				}`,
			),
		},
	}

	got, err := v.Status()
	require.NoError(t, err)
	assert.Equal(t, 2, got.MessageVersion)
}

func TestVacuum_InitialStatus(t *testing.T) {
	t.Parallel()

	v := &Vacuum{
		client: &mockClient{
			rsp: []byte(`
					{
						"result": [{
								"local_info": {
									"name": "custom_A.03.0070_CE",
									"bom": "A.03.0070",
									"location": "de",
									"language": "en",
									"wifiplan": "",
									"timezone": "Europe/Berlin",
									"logserver": "awsde0.fds.api.xiaomi.com",
									"featureset": 1
								},
								"feature_info": [111, 112, 113, 114, 115, 116, 117, 118, 119, 122, 125],
								"status_info": {
									"state": 8,
									"battery": 100,
									"clean_time": 2496,
									"clean_area": 34912500,
									"error_code": 0,
									"in_cleaning": 0,
									"in_returning": 0,
									"in_fresh_state": 1,
									"lab_status": 1,
									"water_box_status": 1,
									"map_status": 3,
									"is_locating": 0,
									"lock_status": 0,
									"water_box_mode": 204,
									"water_box_carriage_status": 1,
									"mop_forbidden_enable": 1
								}
							}
						],
						"id": 6652
					}`,
			),
		},
	}

	got, err := v.InitialStatus()
	require.NoError(t, err)
	assert.Equal(t, "custom_A.03.0070_CE", got.Locale.Name)
}

func TestVacuum_NetworkInfo(t *testing.T) {
	t.Parallel()

	v := &Vacuum{
		client: &mockClient{
			rsp: []byte(`
				{
					"id": 7840,
					"result": {
						"ssid": "iot",
						"ip": "192.168.178.99",
						"mac": "b0:4a:33:04:f6:5f",
						"bssid": "2a:b8:29:e4:cb:45",
						"rssi": -52
					}
				}`,
			),
		},
	}

	got, err := v.NetworkInfo()
	require.NoError(t, err)
	assert.Equal(t, "iot", got.SSID)
}

func TestVacuum_SerialNumber(t *testing.T) {
	t.Parallel()

	v := &Vacuum{
		client: &mockClient{
			rsp: []byte(`
				{
					"result": [{
							"serial_number": "1387100330000"
						}
					],
					"id": 1
				}`,
			),
		},
	}

	got, err := v.SerialNumber()
	require.NoError(t, err)
	assert.Equal(t, "1387100330000", got)
}
