# uStress

Performs high load / rest tests for an endpoint using one or more concurrent requests ( only GET for now ).

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


  stress --url=URL --requests=REQUESTS --workers=WORKERS [<flags>]
    stress a URL

    --url=URL            URL to probe.
    --requests=REQUESTS  Number of request to be sent.
    --workers=WORKERS    Number of concurent workers
    --resolve=RESOLVE    Force resolve of HOST:PORT to ADDRESS
    --insecure           Ignore invalid certificate
    --method=METHOD      HTTP Method to use
    --payload=PAYLOAD    Payload to send
    --headers=HEADERS    Headers to set for request
    --with-response      To return response or not

  web --[no-]start --config=CONFIG [<flags>]
    start the http server

    --start                   Start http server.
    --listen-address=":8080"  Address on which to start the web server
    --config=CONFIG           Path to configuration
```


### UI:
   - new probe request, calls monkey via a websocket connection
   - view stored reports

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

- [x] slack notification
- [ ] swagger docu
- [x] rebuild ui in node
- [x] abort a test if user send bad url (if firs n% of request represents 100% error rate, abort) - 30%
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
