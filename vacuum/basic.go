package vacuum

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

func (v *Vacuum) StartCharge() error {
	return v.doSimple("app_charge")
}

// FindMe plays the vacuums "Find me" phrase.
func (v *Vacuum) FindMe() error {
	return v.doSimple("find_me")
}
