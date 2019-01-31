FROM              golang:1.11.0 AS server-build

WORKDIR           /go/src/git.metrosystems.net/reliability-engineering/ustress
COPY              . .
RUN               go get github.com/golang/dep/cmd/dep
RUN               dep ensure -vendor-only
RUN               ls -al
EXPOSE 8080
ENTRYPOINT        [ "go", "run", "cmd/ustress/main.go", \
    "web", "--start", "--config ./configuration.yaml"]