FROM              golang:1.11.0 

ARG               PROJECT_DIR=/go/src/git.metrosystems.net/reliability-engineering/ustress/
WORKDIR           ${PROJECT_DIR} 
COPY              . .
RUN               go get github.com/golang/dep/cmd/dep
RUN               dep ensure --vendor-only 
RUN               ls -al
