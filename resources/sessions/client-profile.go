package resources

type ClientProfile struct {
	LoginAttempts    map[string]map[string]*LoginAttempt // Key is email address
	AuthorizationMap map[string]uint                     // Key is tenant identifier
}

func newClientProfile() ClientProfile {
	c := ClientProfile{}
	c.LoginAttempts = make(map[string]map[string]*LoginAttempt)
	c.AuthorizationMap = make(map[string]uint)
	return c
}
