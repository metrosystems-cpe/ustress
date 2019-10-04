# uStress

Performs high load / rest tests for an endpoint using one or more concurrent requests.


## Installation 

To use the UI 
```console
git clone https://github.com/metrosystems-cpe/ustress
cd ./ustress
make build_docker
docker run -p "8080:8080" ustress
```

To use the CLI
```console
git clone https://github.com/metrosystems-cpe/ustress
cd ./ustress
make linux
ln ustress-linux-amd64 /usr/local/bin/ustress
```

For Development
```console
docker-compose up
```

## Usage

```console
usage: ustress [<flags>] <command> [<args> ...]

A URL stress application.

Flags:
  --help     Show context-sensitive help (also try --help-long and --help-man).
  --version  Show application version.

Commands:
  help [<command>...]
    Show help.


  stress --url=URL [<flags>]
    stress a URL

    --url=URL              URL to probe.
    --requests=REQUESTS    Number of request to be sent.
    --workers=1            Number of concurent workers
    --payload=PAYLOAD      Payload to send
    --headers=HEADERS      Headers to set for request
    --method="GET"         HTTP Method to use
    --with-response        To return response or not
    --stream-output        Stream output
    --insecure             Ignore invalid certificate
    --resolve=RESOLVE      Force resolve of HOST:PORT to ADDRESS
    --duration=DURATION    Stress duration
    --frequency=FREQUENCY  Requests hit frequency

  web [<flags>]
    start the http server

    --listen-address=":8080"  Address on which to start the web server
    --cassandra-envvar="CASS_CREDS"
                              Env var where cassandra creds are found
```


### Probe handler:
```console
$ curl -s "http://localhost:8080/ustress/api/v1/probe?url=http://localhost:8080/ustress/api/v1/test&requests=4&workers=4" | jq
{
  "entries": {
    "uuid": "81f97756-2eb0-42f8-a928-79158c6b6103",
    "timestamp": "2019-10-04T14:47:01.005038+01:00",
    "config": {
      "url": "http://localhost:8080/ustress/api/v1/test",
      "method": "",
      "requests": 4,
      "threads": 4,
      "resolve": "",
      "insecure": false,
      "payload": "",
      "headers": null,
      "duration": 0,
      "frequency": 0,
      "withResponse": false
    },
    "stats": {
      "median": 0.256141254,
      "50_percentile": 0.256108721,
      "75_percentile": 0.256173787,
      "95_percentile": 0.256180062,
      "99_percentile": 0.256180062,
      "error_percentage": 0,
      "codes_count": {
        "200": 4
      }
    },
    "durationTotal": 0.256393686,
    "data": [
      {
        "request": 1,
        "status": 200,
        "thread": 1,
        "duration": 0.256039672,
        "error": "",
        "response": ""
      },
      {
        "request": 4,
        "status": 200,
        "thread": 4,
        "duration": 0.256186337,
        "error": "",
        "response": ""
      },
      {
        "request": 2,
        "status": 200,
        "thread": 2,
        "duration": 0.256108721,
        "error": "",
        "response": ""
      },
      {
        "request": 3,
        "status": 200,
        "thread": 3,
        "duration": 0.256173787,
        "error": "",
        "response": ""
      }
    ],
    "completed": true
  },
  "error": ""
}
```

## UI

![](https://media.giphy.com/media/H8KIwTlNAu1k13Xr2p/giphy.gif)

## HTTP Handlers

| Handlers                    | Foo            | Description                  |
| --------------------------- |:-------------- |:---------------------------- |
| /                           |                | redirects to ui              |
| /ustress                    |                | redirects to ui              |
| /ustress/ui                 |                | user interface               |
| /ustress/data/              |                | http filesystem              |
| /ustress/api/v1/ws          |                | websocket command endpoint   |
| /ustress/api/v1/reports     |                | exposes local data folder    |
| /ustress/api/v1/probe       |                | http command endpoint        |
| /ustress/api/v1/test        |                | created to test concurrency  |
| /.well-known/ready          |                | -                            |
| /.well-known/live           |                | -                            |
| /.well-known/metrics        |                | -                            |


## Improvements - contributions are welcomed

- [ ] define test in a .yaml file
- [ ] ability to run test in a defined period of a day / year
- [ ] ability to re-run tests automatically  

## Pprof

```
go tool pprof http://localhost:8080/debug/pprof/heap
go tool pprof -top http://localhost:8080/debug/pprof/heap
go tool pprof http://localhost:8080/debug/pprof/goroutine
go tool pprof -png http://localhost:8080/debug/pprof/goroutine > out.png
```
