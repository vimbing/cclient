package cclient

import (
	"crypto/tls"
	"sync"

	"github.com/tam7t/hpkp"
)

type PinnedSite struct {
	Host   string
	Pins   []string
	Failed bool
}

func GeneratePins(hosts []string, pinChannel chan PinnedSite) {
	var wait sync.WaitGroup

	for _, h := range hosts {
		wait.Add(1)

		go func(w *sync.WaitGroup, host string) {
			defer w.Done()
			if pins, err := GetSSLPins(host + ":443"); err == nil {
				pinChannel <- PinnedSite{Host: host, Pins: pins}
			} else {
				pinChannel <- PinnedSite{Failed: true}
			}

		}(&wait, h)
	}

	go func() {
		wait.Wait()
		close(pinChannel)
	}()
}

func GetSSLPins(server string) ([]string, error) {
	var pins []string

	conn, err := tls.Dial("tcp", server, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return pins, err
	}

	for _, cert := range conn.ConnectionState().PeerCertificates {
		pins = append(pins, hpkp.Fingerprint(cert))
	}

	return pins, nil
}
