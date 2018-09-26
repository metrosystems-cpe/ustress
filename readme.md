# REST Monkey

Performs high load / rest tests for an endpoint using one or more concurrent requests ( only GET for now ).

## Usage

Arguments:
 - Required:
   - url      - the endpoint to probe
   - requests - the number of requests to be sent to an endpoint
   - workers  - is the total number concurrent requests

 Optional:
  - insecure  - skip certificate validation in case of https  requests
  - resolve   - similar cu curl resolve (resolves a domain to an ip:port)

### UI:
   - new probe request, calls monkey via a websocket connection
   - view stored reports

### Probe handler:
Warning: if you want to probe a URL with arguments you have to URL encode it

```
# simple request
http://localhost:8080/restmonkey/api/v1/probe?url=http://localhost:8080/restmonkey/api/v1/test&requests=10&workers=4

# insecure request
http://localhost:8080/restmonkey/api/v1/probe?url=http://localhost:8080/restmonkey/api/v1/test&requests=10&workers=4&insecure=true


# will resolve restmonkey.com to 127.0.0.1:8080
# useful in case you want to probe a single dc
http://localhost:8080/restmonkey/api/v1/probe?url=http://restmonkey.com/restmonkey/api/v1/test&requests=10&workers=4&insecure=true&resolve=127.0.0.1:8080
```

## Reports

```json
{
    "uuid": "c2319ed0-ef22-4198-83e6-5323b0315c88",
    "timestamp": "2018-09-26T14:21:33.473406+03:00",
    "config": {
        "URL": "http://restmonkey.com/restmonkey/api/v1/test",
        "Requests": 16,
        "Threads": 4,
        "Resolve": "127.0.0.1:8080",
        "Insecure": false
    },
    "stats": {
        "median": 0.2542787085,
        "50_percentile": 0.254276985,
        "75_percentile": 0.25509104,
        "95_percentile": 0.255408255,
        "99_percentile": 0.255408255,
        "error_percentage": 0
    },
    "durationTotal": 1.016421351,
    "data": [{
        "request": 3,
        "status": 200,
        "thread": 2,
        "duration": 0.254235007,
        "error": ""
    }, {
        "request": 2,
        "status": 200,
        "thread": 1,
        "duration": 0.254280432,
        "error": ""
    }, {
        "request": 1,
        "status": 200,
        "thread": 3,
        "duration": 0.254276985,
        "error": ""
    }, {
        "request": 4,
        "status": 200,
        "thread": 4,
        "duration": 0.2543515,
        "error": ""
    }, {
        "request": 7,
        "status": 200,
        "thread": 3,
        "duration": 0.253150453,
        "error": ""
    }, {
        "request": 8,
        "status": 200,
        "thread": 4,
        "duration": 0.255286629,
        "error": ""
    }, {
        "request": 5,
        "status": 200,
        "thread": 2,
        "duration": 0.25539477,
        "error": ""
    }, {
        "request": 6,
        "status": 200,
        "thread": 1,
        "duration": 0.25542174,
        "error": ""
    }, {
        "request": 9,
        "status": 200,
        "thread": 3,
        "duration": 0.253625896,
        "error": ""
    }, {
        "request": 12,
        "status": 200,
        "thread": 1,
        "duration": 0.255065439,
        "error": ""
    }, {
        "request": 11,
        "status": 200,
        "thread": 2,
        "duration": 0.25509104,
        "error": ""
    }, {
        "request": 10,
        "status": 200,
        "thread": 4,
        "duration": 0.255145366,
        "error": ""
    }, {
        "request": 13,
        "status": 200,
        "thread": 3,
        "duration": 0.25427686,
        "error": ""
    }, {
        "request": 16,
        "status": 200,
        "thread": 4,
        "duration": 0.251557443,
        "error": ""
    }, {
        "request": 14,
        "status": 200,
        "thread": 1,
        "duration": 0.251566642,
        "error": ""
    }, {
        "request": 15,
        "status": 200,
        "thread": 2,
        "duration": 0.251576823,
        "error": ""
    }]
}
--- abort in case of http error
{
    "uuid": "bb2f72fe-f371-43cd-bbda-abf6c66fd0d3",
    "timestamp": "2018-09-26T14:27:08.988795+03:00",
    "config": {
        "URL": "http://restmonkey.com/restmonkey/api/v1/test",
        "Requests": 16,
        "Threads": 4,
        "Resolve": "127.0.0.1:1234",
        "Insecure": false
    },
    "stats": {
        "median": 0,
        "50_percentile": 0,
        "75_percentile": 0,
        "95_percentile": 0,
        "99_percentile": 0,
        "error_percentage": 25
    },
    "durationTotal": 0.002375223,
    "data": [{
        "request": 2,
        "status": 0,
        "thread": 2,
        "duration": 0.001668874,
        "error": "Get http://restmonkey.com/restmonkey/api/v1/test: dial tcp 127.0.0.1:1234: connect: connection refused"
    }, {
        "request": 3,
        "status": 0,
        "thread": 1,
        "duration": 0.001562397,
        "error": "Get http://restmonkey.com/restmonkey/api/v1/test: dial tcp 127.0.0.1:1234: connect: connection refused"
    }, {
        "request": 1,
        "status": 0,
        "thread": 3,
        "duration": 0.001494121,
        "error": "Get http://restmonkey.com/restmonkey/api/v1/test: dial tcp 127.0.0.1:1234: connect: connection refused"
    }, {
        "request": 4,
        "status": 0,
        "thread": 4,
        "duration": 0.001680404,
        "error": "Get http://restmonkey.com/restmonkey/api/v1/test: dial tcp 127.0.0.1:1234: connect: connection refused"
    }]
}

```



## HTTP Handlers

@todo - swager :)

| Handlers                    | Foo            | Description    |
| --------------------------- |:-------------- |:---------------------------- |
| /                           |                | redirects to ui              |
| /restmonkey                 |                | redirects to ui              |
| /restmonkey/ui              |                | user interface               |
| /restmonkey/data/           |                | http filesystem              |
| /restmonkey/api/v1/ws       |                | websocket command endpoint   |
| /restmonkey/api/v1/reports  |                | exposes local data folder    |
| /restmonkey/api/v1/probe    |                | http command endpoint        |
| /restmonkey/api/v1/test     |                | created to test concurrency  |
| /.well-known/ready          |                | -                            |
| /.well-known/live           |                | -                            |
| /.well-known/metrics        |                | -                            |


## Improvements - contributions are welcomed

- [x] slack notification
- [ ] swagger docu
- [ ] rebuild ui in node
- [x] abort a test if user send bad url (if firs n% of request represents 100% error rate, abort) - 30%
- [ ] you know when you write bad code when you cannot define go tests ( #me )
- [ ] define test in a .yaml file
 - [ ] ability to run test in a defined period of a day / year
 - [ ] ability to re-run tests automatically  
- [ ] ability to make auth tests | needs advanced http client config  
- [ ] save reports in a database, cache with TTL


## Pprof

```
go tool pprof http://localhost:8080/debug/pprof/heap
go tool pprof -top http://localhost:8080/debug/pprof/heap
go tool pprof http://localhost:8080/debug/pprof/goroutine
go tool pprof -png http://localhost:8080/debug/pprof/goroutine > out.png
```