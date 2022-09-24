package vacuum

import "encoding/json"

func (v *Vacuum) RemoteStart() error {
	return v.doSimple("app_rc_start")
}

func (v *Vacuum) RemoteEnd() error {
	return v.doSimple("app_rc_end")
}

type RemoteMoveParams struct {
	// Angle must be between -3.1 and 3.1
	Angle float32 `json:"omega"`
	// Velocity must be between -0.3 and 0.3
	Velocity float32 `json:"velocity"`
	// Sequence needs to increase with each instance of calling RemoteMove
	Sequence int `json:"seqnum"`
	// Duration in ms
	Duration int `json:"duration"`
}

// RemoteMove sends the provided move commands.
// RemoteStart must be called first and the vacuums should be in the ManualMode state.
// The vacuums state can be checked via the Status method.
func (v *Vacuum) RemoteMove(params []RemoteMoveParams) error {
	p, err := json.Marshal(params)
	if err != nil {
		return err
	}

	return v.do("app_rc_move", p, nil)
}
