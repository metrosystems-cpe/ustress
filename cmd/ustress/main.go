package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"git.metrosystems.net/reliability-engineering/ustress/log"
	"git.metrosystems.net/reliability-engineering/ustress/ustress"
	"git.metrosystems.net/reliability-engineering/ustress/web"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	verbose = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	app     = kingpin.New("stress", "A URL stress application.")

	stress   = app.Command("stress", "stress a URL")
	url      = stress.Flag("stress.url", "URL to probe.").Required().String()
	requests = stress.Flag("stress.requests", "Number of request to be sent.").Required().Int()
	workers  = stress.Flag("stress.workers", "Number of concurent workers").Required().Int()
	resolve  = stress.Flag("stress.resolve", "Force resolve of HOST:PORT to ADDRESS").String()
	insecure = stress.Flag("stress.insecure", "Ignore invalid certificate").Bool()
	method   = stress.Flag("stress.method", "HTTP Method to use").String()
	payload  = stress.Flag("stress.payload", "Payload to send").String()
	headers  = stress.Flag("stress.headers", "Headers to set for request").String()

	webServer     = app.Command("web", "start the http server")
	startWeb      = webServer.Flag("web.start", "Start http server.").Required().Bool()
	listenAddress = webServer.Flag("web.listen-address", "Address on which to start the web server").Default(":8080").String()
)

func loadHeaders(headers string) map[string]string {
	headersList := strings.Split(headers, ";")
	headersMap := make(map[string]string)
	for _, h := range headersList {
		header := strings.Split(h, ":")
		if len(header) < 2 {
			return nil
		}
		headersMap[header[0]] = header[1]
	}
	return headersMap
}

func main() {
	app.Version("0.0.1")
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Manual stress
	case stress.FullCommand():
		restMK := ustress.NewConfig(*url, *requests, *workers, *resolve, *insecure, *method, *payload, loadHeaders(*headers))
		fmt.Printf("%#v", restMK)
		messages, err := ustress.NewReport(restMK)
		if err != nil {
			log.LogWithFields.Error(err.Error())
		}
		fmt.Println(string(messages))
	// start the web server
	case webServer.FullCommand():
		if *startWeb {
			mux := web.MuxHandlers()
			log.LogWithFields.Infof("Starting proxy server on: %v", *listenAddress)
			if err := http.ListenAndServe(*listenAddress, mux); err != nil {
				log.LogWithFields.Fatalf("ListenAndServe: %v", err.Error())
			}
		}
	}
}
