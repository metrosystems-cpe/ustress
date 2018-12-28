package ustress

import (
	"net/http"
	"testing"
	"time"
)

func TestNewHTTPClientDefault(t *testing.T) {
	monkeyConfig := MonkeyConfig{
		URL:      "http://localhost:8080/test",
		Requests: 16,
		Threads:  4,
		Resolve:  "",
		Insecure: false,
	}
	expectedTimeout := time.Duration(2 * time.Second)

	t.Log("client.Timeout")
	monkeyConfig.client = monkeyConfig.newHTTPClient()
	if monkeyConfig.client.Timeout != expectedTimeout {
		t.Errorf("timeout not properly exported: expected %6v, got %6v ",
			expectedTimeout, monkeyConfig.client.Timeout)
	}

	clientConfig := monkeyConfig.client.Transport.(*http.Transport)

	t.Log("httpClient Transporter MaxIdleConns")
	if clientConfig.MaxIdleConns != monkeyConfig.Threads {
		t.Errorf("timeout not properly exported: expected %6v, got %6v ",
			clientConfig.MaxIdleConns, monkeyConfig.Threads)
	}
	t.Log("httpClient Transporter MaxIdleConnsPerHost")
	if clientConfig.MaxIdleConnsPerHost != monkeyConfig.Threads {
		t.Errorf("timeout not properly exported: expected %6v, got %6v ",
			clientConfig.MaxIdleConnsPerHost, monkeyConfig.Threads)
	}
}

func TestNewHTTPClientDefaultInsecure(t *testing.T) {
	monkeyConfig := MonkeyConfig{
		URL:      "http://localhost:8080/test",
		Requests: 16,
		Threads:  4,
		Resolve:  "",
		Insecure: true,
	}

	expectedTimeout := time.Duration(2 * time.Second)

	t.Log("client.Timeout")
	monkeyConfig.client = monkeyConfig.newHTTPClient()
	if monkeyConfig.client.Timeout != expectedTimeout {
		t.Errorf("timeout not properly exported: expected %6v, got %6v ",
			expectedTimeout, monkeyConfig.client.Timeout)
	}

	clientConfig := monkeyConfig.client.Transport.(*http.Transport)

	t.Log("httpClient Transporter MaxIdleConns")
	if clientConfig.MaxIdleConns != monkeyConfig.Threads {
		t.Errorf("timeout not properly exported: expected %6v, got %6v ",
			clientConfig.MaxIdleConns, monkeyConfig.Threads)
	}
	t.Log("httpClient Transporter MaxIdleConnsPerHost")
	if clientConfig.MaxIdleConnsPerHost != monkeyConfig.Threads {
		t.Errorf("timeout not properly exported: expected %6v, got %6v ",
			clientConfig.MaxIdleConnsPerHost, monkeyConfig.Threads)
	}
}

func TestNewHTTPClientDefaultResolve(t *testing.T) {
	monkeyConfig := MonkeyConfig{
		URL:      "http://localhost:8080/test",
		Requests: 16,
		Threads:  4,
		Resolve:  "10.29.80.28:443",
		Insecure: true,
	}

	expectedTimeout := time.Duration(2 * time.Second)

	t.Log("client.Timeout")
	monkeyConfig.client = monkeyConfig.newHTTPClient()
	if monkeyConfig.client.Timeout != expectedTimeout {
		t.Errorf("timeout not properly exported: expected %6v, got %6v ",
			expectedTimeout, monkeyConfig.client.Timeout)
	}

	clientConfig := monkeyConfig.client.Transport.(*http.Transport)

	t.Log("httpClient Transporter MaxIdleConns")
	if clientConfig.MaxIdleConns != monkeyConfig.Threads {
		t.Errorf("timeout not properly exported: expected %6v, got %6v ",
			clientConfig.MaxIdleConns, monkeyConfig.Threads)
	}
	t.Log("httpClient Transporter MaxIdleConnsPerHost")
	if clientConfig.MaxIdleConnsPerHost != monkeyConfig.Threads {
		t.Errorf("timeout not properly exported: expected %6v, got %6v ",
			clientConfig.MaxIdleConnsPerHost, monkeyConfig.Threads)
	}
}
