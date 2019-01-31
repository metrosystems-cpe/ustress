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
	Tr         = &http.Transport{}
)

func (mkcfg *MonkeyConfig) newHTTPClient() *http.Client {
	timeout := time.Duration(2 * time.Second)

	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: timeout,
		DualStack: true,
	}
	Tr = &http.Transport{
		MaxIdleConns:        mkcfg.Threads,
		MaxIdleConnsPerHost: mkcfg.Threads,
		Dial:                (dialer).Dial,
		TLSHandshakeTimeout: timeout,
	}

	// resolve ip
	if mkcfg.Resolve != "" {
		Tr.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, mkcfg.Resolve)
		}
	}

	// insecure request
	if mkcfg.Insecure {
		Tr.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: Tr,
	}

}
