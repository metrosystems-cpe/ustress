# uStress

Performs high load / rest tests for an endpoint using one or more concurrent requests ( only GET for now ).


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
Warning: if you want to probe a URL with arguments you have to URL encode it


```
# simple request
http://localhost:8080/ustress/api/v1/probe?url=http://localhost:8080/ustress/api/v1/test&requests=10&workers=4

# insecure request
http://localhost:8080/ustress/api/v1/probe?url=http://localhost:8080/ustress/api/v1/test&requests=10&workers=4&insecure=true


# will resolve ustress.com to 127.0.0.1:8080
# useful in case you want to probe a single dc
http://localhost:8080/ustress/api/v1/probe?url=http://ustress.com/ustress/api/v1/test&requests=10&workers=4&insecure=true&resolve=127.0.0.1:8080
```

## HTTP Handlers

@todo - swager :)

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

- [ ] swagger doc
- [ ] abort a test if user send bad url (if firs n% of request represents 100% error rate, abort) - 30%
- [ ] you know when you write bad code when you cannot define go tests ( #me )
- [ ] define test in a .yaml file
- [ ] ability to run test in a defined period of a day / year
- [ ] ability to re-run tests automatically  
- [x] ability to make auth tests | needs advanced http client config  
- [x] save reports in a database, cache with TTL
- [ ] ability run tests over a time period with a specified request frequency 



## Pprof

```
go tool pprof http://localhost:8080/debug/pprof/heap
go tool pprof -top http://localhost:8080/debug/pprof/heap
go tool pprof http://localhost:8080/debug/pprof/goroutine
go tool pprof -png http://localhost:8080/debug/pprof/goroutine > out.png
```
