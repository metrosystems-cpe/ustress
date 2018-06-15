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
/restmonkey/api/v1/probe??url=http://localhost:8080&requests=10&workers=4

# will resolve idam to be-gcw datacenter's LB
/restmonkey/api/v1/probe?insecure=true&resolve=10.29.80.28:443&url=https://idam-pp.metrosystems.net/.well-known/openid-configuration&requests=10&workers=2
```

## Reports



```json
{
    "uuid": "f2b260af-305f-4c6a-9063-577593c1c882",
    "timestamp": "2018-06-12T15:13:02.907475+03:00",
    "config": {
        "URL": "http://localhost:8080/restmonkey/api/v1/test",
        "Requests": 4,
        "Threads": 4,
        "Resolve": "",
        "Insecure": false
    },
    "stats": {
        "median": 0.25212597,
        "50_percentile": 0.2521247,
        "75_percentile": 0.25212724,
        "95_percentile": 0.2521405025,
        "99_percentile": 0.2521405025,
        "error_percentage": 0
    },
    "durationTotal": 0.502027589,
    "data": [{
        "request": 2,
        "status": 200,
        "thread": 3,
        "duration": 0.25212724,
        "error": ""
    }, {
        "request": 3,
        "status": 200,
        "thread": 1,
        "duration": 0.252122993,
        "error": ""
    }, {
        "request": 4,
        "status": 200,
        "thread": 2,
        "duration": 0.2521247,
        "error": ""
    }, {
        "request": 1,
        "status": 200,
        "thread": 4,
        "duration": 0.252153765,
        "error": ""
    }]
}

---
{
    "uuid": "3a94b039-a321-4bdf-be98-813f869dd20f",
    "timestamp": "2018-06-12T15:10:53.532581+03:00",
    "config": {
        "URL": "/restmonkey/api/v1/test",
        "Requests": 4,
        "Threads": 4,
        "Resolve": "",
        "Insecure": false
    },
    "stats": {
        "median": 0,
        "50_percentile": 0,
        "75_percentile": 0,
        "95_percentile": 0,
        "99_percentile": 0,
        "error_percentage": 100
    },
    "durationTotal": 0.507984186,
    "data": [{
        "request": 2,
        "status": 0,
        "thread": 4,
        "duration": 0.004012851,
        "error": "Get /restmonkey/api/v1/test: unsupported protocol scheme \"\""
    }, {
        "request": 3,
        "status": 0,
        "thread": 1,
        "duration": 0.004014902,
        "error": "Get /restmonkey/api/v1/test: unsupported protocol scheme \"\""
    }, {
        "request": 4,
        "status": 0,
        "thread": 3,
        "duration": 0.004008321,
        "error": "Get /restmonkey/api/v1/test: unsupported protocol scheme \"\""
    }, {
        "request": 1,
        "status": 0,
        "thread": 2,
        "duration": 0.004039681,
        "error": "Get /restmonkey/api/v1/test: unsupported protocol scheme \"\""
    }]
}

```



## HTTP Handlers

@todo - swager :)

| Handlers                    | Foo            | Bar            |
| --------------------------- |:-------------- | --------------:|
| /                           |                |                |
| /restmonkey                 |                |                |
| /restmonkey/ui              |                |                |
| /restmonkey/data/           |                |                |
| /restmonkey/api/v1/ws       |                |                |
| /restmonkey/api/v1/reports  |                |                |
| /restmonkey/api/v1/probe    |                |                |
| /restmonkey/api/v1/test     |                |                |
| /.well-known/ready          |                |                |
| /.well-known/live           |                |                |
| /.well-known/metrics        |                |                |


## Improvements - contributions are welcomed

- [x] slack notification
- [ ] swager
- [ ] abort a test if user send bad url (if firs n% of request represents 100% error rate, abort)
- [ ] you know when you write bad code when you cannot define go tests ( #me )
- [ ] define test in a .yaml file
 - [ ] ability to run test in a defined period of a day / year
 - [ ] ability to re-run tests automatically  
- [ ] ability to make auth tests | needs advanced http client config  
- [ ] save reports in a database with TTL
