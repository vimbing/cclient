package cclient

import "github.com/tam7t/hpkp"

type SSLPinningSecurityError struct{}

func (s SSLPinningSecurityError) Error() string {
	return "ssl_pinning_security_error"
}

type ClientOptions struct {
	SSLPinningOptions SSLPinningOptions
}

type SSLPinningOptions struct {
	Required         bool
	Hosts            []string
	AutoGeneratePins bool
	Storage          hpkp.MemStorage
	Notifier         func()
}
