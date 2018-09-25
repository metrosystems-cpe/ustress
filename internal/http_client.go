package internal

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

func (mkcfg *MonkeyConfig) newHTTPClient() *http.Client {

	// resolve ip
	if mkcfg.Resolve != "" {

		dialer := &net.Dialer{
			Timeout:   2 * time.Second,
			KeepAlive: 0 * time.Second,
			DualStack: true,
		}

		tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, mkcfg.Resolve)
		}
	}

	// insecure request
	if mkcfg.Insecure {
		tr = &http.Transport{
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			MaxIdleConns:        mkcfg.Threads, // this should be set as the number of go routines
			MaxIdleConnsPerHost: mkcfg.Threads,
		}
	}

	return &http.Client{
		Timeout:   1 * time.Second,
		Transport: tr,
	}

}
