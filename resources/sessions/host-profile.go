package resources

import (
	"time"
)

type HostProfile struct {
	LoginAttempts        map[string]*LoginAttempt // Key is used email address
	LastLoginAttemptTime time.Time
	AuthorizedTime       time.Time
	UserId               uint
	Authorized           uint
}

func newHostProfile() HostProfile {
	h := HostProfile{}
	h.LoginAttempts = make(map[string]*LoginAttempt)
	h.Authorized = 0
	return h
}
