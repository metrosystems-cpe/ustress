package main

import (
	"fmt"
	"git.metrosystems.net/reliability-engineering/ustress/ustress"
	"net/http"
	"os"
	"strings"

	"git.metrosystems.net/reliability-engineering/ustress/log"
	"git.metrosystems.net/reliability-engineering/ustress/web"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	appVersion = "0.1.0"
	verbose    = kingpin.Flag("verbose", "Verbose mode.").Short('v').Bool()
	app        = kingpin.New("ustress", "A URL stress application.")

	stress = app.Command("stress", "stress a URL")
	webServer = app.Command("web", "start the http server")
	startWeb      = webServer.Flag("start", "Start http server.").Required().Bool()
	listenAddress = webServer.Flag("listen-address", "Address on which to start the web server").Default(":8080").String()
	url	= stress.Flag("url", "URL to probe.").Required().String()
	requests = stress.Flag("requests", "Number of request to be sent.").Int()
	workers = stress.Flag("workers", "Number of concurent workers").Default("1").Int()
	payload = stress.Flag("payload", "Payload to send").String()
	headers = stress.Flag("headers", "Headers to set for request").String()
	method = stress.Flag("method", "HTTP Method to use").Default("GET").String()
	withResponse = stress.Flag("with-response", "To return response or not").Default("false").Bool()
	streamOut = stress.Flag("stream-output", "Stream output").Default("false").Bool()
	insecure = stress.Flag("insecure", "Ignore invalid certificate").Bool()
	resolve = stress.Flag("resolve", "Force resolve of HOST:PORT to ADDRESS").String()
	duration = stress.Flag("duration", "Stress duration").Int()
	frequency = stress.Flag("frequency", "Requests hit frequency").Int()
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

func streamOutput(r *ustress.Report,q  chan bool) {
	fmt.Print(string(r.JSON()))
}

func main() {
	app.Version(appVersion)

	var (
		report *ustress.Report
		err error
	)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Manual stress
	case stress.FullCommand():
		cfg,_ := ustress.NewStressConfig(
			ustress.NewOption("Duration", *duration),
			ustress.NewOption("URL", *url),
			ustress.NewOption("Frequency", *frequency),
			ustress.NewOption("WithResponse", *withResponse),
			ustress.NewOption("Threads", *workers),
			ustress.NewOption("Resolve", *resolve),
			ustress.NewOption("Insecure", *insecure),
			ustress.NewOption("Method", *method),
			ustress.NewOption("Headers", loadHeaders(*headers)),
			ustress.NewOption("Requests", *requests),
		)

		if *streamOut {
			report, err = ustress.NewReport(cfg, streamOutput, 500)

		} else {
			report, err = ustress.NewReport(cfg, nil, 0)

		}

		if err != nil {
			
			log.LogWithFields.Error(err.Error())
		}
		jsonReport := report.JSON()
		fmt.Println(string(jsonReport))

	// start the web server
	case webServer.FullCommand():
		if *startWeb {
			a, err := web.NewAppFromEnv()
			log.LogError(err)
			if a == nil {
				a = web.NewApp(appVersion, web.LocalCassandraConfig())
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
