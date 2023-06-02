package timeout

import "time"

var Timeouts = map[string]bool{}

func StartTimeout(name string, duration time.Duration) {
	Timeouts[name] = true
	select {
	case <-time.After(duration):
		Timeouts[name] = false
	}
}

func GetTimeout(name string) bool {
	if timeout, ok := Timeouts[name]; ok {
		return timeout
	}
	return false
}
