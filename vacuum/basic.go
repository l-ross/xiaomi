package vacuum

import "encoding/json"

// Start cleaning
func (v *Vacuum) Start() error {
	return v.doSimple("app_start")
}

// Stop cleaning
func (v *Vacuum) Stop() error {
	return v.doSimple("app_stop")
}

// StartSpot will start spot cleaning
func (v *Vacuum) StartSpot() error {
	return v.doSimple("app_spot")
}

// Pause cleaning
func (v *Vacuum) Pause() error {
	return v.doSimple("app_pause")
}

// StartCharge will start charging
func (v *Vacuum) StartCharge() error {
	return v.doSimple("app_charge")
}

// FindMe plays the vacuums "Find me" phrase.
func (v *Vacuum) FindMe() error {
	return v.doSimple("find_me")
}

func (v *Vacuum) WakeupRobot() error {
	return v.doSimple("app_wakeup_robot")
}

type DNDTimer struct {
	Enabled     bool
	StartHour   int
	StartMinute int
	EndHour     int
	EndMinute   int
}

// GetDNDTimer retrieves the Do Not Disturb timers for the Vacuum.
func (v *Vacuum) GetDNDTimer() ([]DNDTimer, error) {
	type dndTimer struct {
		Enabled     int `json:"enabled"`
		StartHour   int `json:"start_hour"`
		StartMinute int `json:"start_minute"`
		EndHour     int `json:"end_hour"`
		EndMinute   int `json:"end_minute"`
	}

	rsp := make([]dndTimer, 0)

	err := v.do("get_dnd_timer", nil, &rsp)
	if err != nil {
		return nil, err
	}

	timers := make([]DNDTimer, len(rsp))

	for i, d := range rsp {
		timers[i] = DNDTimer{
			Enabled:     d.Enabled == 1,
			StartHour:   d.StartHour,
			StartMinute: d.StartMinute,
			EndHour:     d.EndHour,
			EndMinute:   d.EndMinute,
		}
	}

	return timers, nil
}

type SetDNDTimerParams struct {
	StartHour   int
	StartMinute int
	EndHour     int
	EndMinutes  int
}

// SetDNDTimer sets the Do Not Disturb timers for the Vacuum.
func (v *Vacuum) SetDNDTimer(params SetDNDTimerParams) error {
	p, err := json.Marshal([]int{params.StartHour, params.StartMinute, params.EndHour, params.EndMinutes})
	if err != nil {
		return err
	}

	return v.do("set_dnd_timer", p, nil)
}
