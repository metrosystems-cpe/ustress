package ustress

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
)

var (
	httpClient = &http.Client{}
	tr         = &http.Transport{}
)

func (mkcfg *MonkeyConfig) newHTTPClient() *http.Client {
	timeout := time.Duration(2 * time.Second)

	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: timeout,
		DualStack: true,
	}
	tr = &http.Transport{
		MaxIdleConns:        mkcfg.Threads,
		MaxIdleConnsPerHost: mkcfg.Threads,
		Dial:                (dialer).Dial,
		TLSHandshakeTimeout: timeout,
	}

	// resolve ip
	if mkcfg.Resolve != "" {
		tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, mkcfg.Resolve)
		}
	}

	// insecure request
	if mkcfg.Insecure {
		tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

}
