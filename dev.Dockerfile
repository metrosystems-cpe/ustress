FROM              golang:1.11.0 

ARG               PROJECT_DIR=
WORKDIR           ${PROJECT_DIR} 
COPY              . .
RUN               go get ./...
RUN               ls -al
