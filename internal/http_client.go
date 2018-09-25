package internal

import (
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

func (mkcfg *MonkeyConfig) newHTTPClient() *http.Client {

	tr = &http.Transport{
		MaxIdleConns:        mkcfg.Threads,
		MaxIdleConnsPerHost: mkcfg.Threads,
	}

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
		tr.TLSClientConfig.InsecureSkipVerify = true
	}

	return &http.Client{
		Timeout:   1 * time.Second,
		Transport: tr,
	}

}
