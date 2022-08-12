package cclient

import (
	"time"

	"github.com/Carcraftz/fhttp/cookiejar"

	http "github.com/Carcraftz/fhttp"

	"golang.org/x/net/proxy"

	utls "github.com/Carcraftz/utls"
)

func NewClient(clientHello utls.ClientHelloID, proxyUrl string, allowRedirect bool, timeout time.Duration) (http.Client, error) {
	if len(proxyUrl) > 0 {
		dialer, err := newConnectDialer(proxyUrl)
		if err != nil {
			if allowRedirect {
				cJar, _ := cookiejar.New(nil)
				return http.Client{
					Jar:     cJar,
					Timeout: time.Second * timeout,
				}, err
			}
			cJar, _ := cookiejar.New(nil)
			return http.Client{
				Jar:     cJar,
				Timeout: time.Second * timeout,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					return http.ErrUseLastResponse
				},
			}, err
		}
		if allowRedirect {
			cJar, _ := cookiejar.New(nil)
			return http.Client{
				Jar:       cJar,
				Transport: newRoundTripper(clientHello, dialer),
				Timeout:   time.Second * timeout,
			}, nil
		}
		cJar, _ := cookiejar.New(nil)
		return http.Client{
			Jar:       cJar,
			Transport: newRoundTripper(clientHello, dialer),
			Timeout:   time.Second * timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}, nil
	} else {
		if allowRedirect {
			cJar, _ := cookiejar.New(nil)
			return http.Client{
				Jar:       cJar,
				Transport: newRoundTripper(clientHello, proxy.Direct),
				Timeout:   time.Second * timeout,
			}, nil
		}
		cJar, _ := cookiejar.New(nil)
		return http.Client{
			Jar:       cJar,
			Transport: newRoundTripper(clientHello, proxy.Direct),
			Timeout:   time.Second * timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		}, nil

	}
}
