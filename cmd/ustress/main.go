package main

import (
	"fmt"
	"net/http"
	"os"

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

	webServer     = app.Command("web", "start the http server")
	startWeb      = webServer.Flag("web.start", "Start http server.").Required().Bool()
	listenAddress = webServer.Flag("web.listen-address", "Address on which to start the web server").Default(":8080").String()
)

func main() {
	app.Version("0.0.1")
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	// Manual stress
	case stress.FullCommand():
		restMK := ustress.NewConfig(*url, *requests, *workers, *resolve, *insecure)
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
