package cclient

import (
	"errors"
	"reflect"
	"time"

	"github.com/tam7t/hpkp"
	http "github.com/vimbing/fhttp"

	"golang.org/x/net/proxy"

	utls "github.com/vimbing/utls"
)

func handleSslPinningOption(opts SSLPinningOptions) (SSLPinningOptions, error) {
	if reflect.DeepEqual(opts, SSLPinningOptions{}) || !opts.Required {
		return SSLPinningOptions{}, nil
	}

	if opts.AutoGeneratePins {
		pinStorage := hpkp.MemStorage{}

		if len(opts.Hosts) == 0 {
			return SSLPinningOptions{}, errors.New("no hosts to pin")
		}

		pinChannel := make(chan PinnedSite, 1000)
		GeneratePins(opts.Hosts, pinChannel)

		for pinned := range pinChannel {
			pinStorage.Add(pinned.Host, &hpkp.Header{
				Permanent:  true,
				Sha256Pins: pinned.Pins,
			})
		}

		opts.Storage = pinStorage

		return opts, nil
	}

	return opts, nil
}

func NewClient(clientHello utls.ClientHelloID, proxyUrl string, allowRedirect bool, timeout time.Duration, opts ...ClientOptions) (http.Client, error) {
	var sslPinningOptions SSLPinningOptions
	var err error

	if len(opts) > 0 {
		sslPinningOptions, err = handleSslPinningOption(opts[0].SSLPinningOptions)

		if err != nil {
			return http.Client{}, err
		}

	}

	if len(proxyUrl) > 0 {
		dialer, err := newConnectDialer(proxyUrl)
		if err != nil {
			if allowRedirect {
				return http.Client{
					Timeout: time.Second * timeout,
				}, err
			}
			return http.Client{
				Timeout: time.Second * timeout,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}, err
		}
		if allowRedirect {
			return http.Client{
				Transport: newRoundTripper(clientHello, sslPinningOptions, dialer),
				Timeout:   time.Second * timeout,
			}, nil
		}
		return http.Client{
			Transport: newRoundTripper(clientHello, sslPinningOptions, dialer),
			Timeout:   time.Second * timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}, nil
	} else {
		if allowRedirect {
			return http.Client{
				Transport: newRoundTripper(clientHello, sslPinningOptions, proxy.Direct),
				Timeout:   time.Second * timeout,
			}, nil
		}

		return http.Client{
			Transport: newRoundTripper(clientHello, sslPinningOptions, proxy.Direct),
			Timeout:   time.Second * timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}, nil

	}
}
