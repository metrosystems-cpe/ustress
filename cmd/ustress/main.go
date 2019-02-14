package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"git.metrosystems.net/reliability-engineering/ustress/log"
	"git.metrosystems.net/reliability-engineering/ustress/ustress"
	"git.metrosystems.net/reliability-engineering/ustress/web"
	"git.metrosystems.net/reliability-engineering/ustress/web/core"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	appVersion = "0.0.1"
	verbose    = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	app        = kingpin.New("ustress", "A URL stress application.")

	stress       = app.Command("stress", "stress a URL")
	url          = stress.Flag("url", "URL to probe.").Required().String()
	requests     = stress.Flag("requests", "Number of request to be sent.").Required().Int()
	workers      = stress.Flag("workers", "Number of concurent workers").Default("1").Int()
	payload      = stress.Flag("payload", "Payload to send").String()
	headers      = stress.Flag("headers", "Headers to set for request").String()
	method       = stress.Flag("method", "HTTP Method to use").Default("GET").String()
	withResponse = stress.Flag("with-response", "To return response or not").Default("false").Bool()
	insecure     = stress.Flag("insecure", "Ignore invalid certificate").Bool()
	resolve      = stress.Flag("resolve", "Force resolve of HOST:PORT to ADDRESS").String()

	webServer     = app.Command("web", "start the http server")
	startWeb      = webServer.Flag("start", "Start http server.").Required().Bool()
	listenAddress = webServer.Flag("listen-address", "Address on which to start the web server").Default(":8080").String()
	configPath    = webServer.Flag("config", "Path to configuration").String()
)

func loadHeaders(headers string) map[string]string {
	headersList := strings.Split(headers, ";")
	headersMap := make(map[string]string)
	for _, h := range headersList {
		header := strings.Split(h, ":")
		if len(header) < 2 {
			return nil
		}
		headersMap[strings.TrimSpace(header[0])] = strings.TrimSpace(header[1])
	}
	return headersMap
}

func main() {
	app.Version(appVersion)
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Manual stress
	case stress.FullCommand():
		restMK := ustress.NewConfig(
			*url, *requests, *workers,
			*resolve, *insecure, *method,
			*payload, loadHeaders(*headers), *withResponse)
		fmt.Printf("%#v", restMK)
		report, err := ustress.NewReport(restMK)
		if err != nil {
			log.LogWithFields.Error(err.Error())
		}
		jsonReport := report.JSON()
		fmt.Println(string(jsonReport))
	// start the web server
	case webServer.FullCommand():
		var a *core.App
		if *startWeb {

			cpath := *configPath
			if cpath != "" {
				a = web.NewAppFromYAML(*configPath)
			} else {
				a = web.NewAppFromYAML("./configuration.yaml")
			}
			mux := web.MuxHandlers(a)
			defer a.Session.Close()
			log.LogWithFields.Infof("Starting proxy server on: %v", *listenAddress)
			if err := http.ListenAndServe(*listenAddress, mux); err != nil {
				log.LogWithFields.Fatalf("ListenAndServe: %v", err.Error())
			}
		}
	}
}
