# Build react app into static files
FROM node:10-alpine as ui-build 
ENV               PROJECT_DIR /go/src/git.metrosystems.net/reliability-engineering/ustress/
RUN               mkdir -p ${PROJECT_DIR}/web/ui 
WORKDIR           ${PROJECT_DIR}/web/ui
COPY              web/ui .
RUN               npm install
RUN               npm run build
RUN               pwd && ls -al

# Building app and injecting static files
FROM              golang:1.11.0 AS server-build
ENV               PROJECT_DIR /go/src/git.metrosystems.net/reliability-engineering/ustress/
RUN               mkdir -p ${PROJECT_DIR}
WORKDIR           ${PROJECT_DIR}
COPY              . .
COPY              --from=ui-build ${PROJECT_DIR}/web/ui ${PROJECT_DIR}/web/ui/.
RUN               go get github.com/golang/dep/cmd/dep
RUN               dep ensure -vendor-only
RUN               make linux
RUN               pwd && ls -al

FROM              quay.io/prometheus/busybox:latest
COPY              --from=server-build /go/src/git.metrosystems.net/reliability-engineering/ustress/ustress-linux-amd64 /ustress
COPY              --from=server-build /go/src/git.metrosystems.net/reliability-engineering/ustress/web/ui /web/ui
COPY              configuration.yaml .
RUN               pwd && ls -al
EXPOSE            8080 
ENTRYPOINT        [ "/ustress", "web", "--start" ]