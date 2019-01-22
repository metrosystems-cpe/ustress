FROM              golang:1.11.0 AS server-build
WORKDIR           /go/src/git.metrosystems.net/reliability-engineering/ustress
COPY              . .
RUN               go get github.com/golang/dep/cmd/dep
RUN               dep ensure -vendor-only
RUN               make linux
RUN               ls -al

FROM              quay.io/prometheus/busybox:latest
COPY              --from=server-build /go/src/git.metrosystems.net/reliability-engineering/ustress/ustress-linux-amd64 /ustress
COPY              web/ui                                 /web/ui
RUN               ls -al
EXPOSE            8080 
ENTRYPOINT        [ "/ustress", "web", "--web.start" ]