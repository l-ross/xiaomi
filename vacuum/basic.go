package vacuum

func (v *Vacuum) Start() error {
	return v.doSimple("app_start")
}

func (v *Vacuum) Stop() error {
	return v.doSimple("app_stop")
}

func (v *Vacuum) StartSpot() error {
	return v.doSimple("app_spot")
}

func (v *Vacuum) Pause() error {
	return v.doSimple("app_pause")
}

func (v *Vacuum) StartCharge() error {
	return v.doSimple("app_charge")
}

func (v *Vacuum) FindMe() error {
	return v.doSimple("find_me")
}
