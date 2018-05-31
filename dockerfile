FROM quay.io/prometheus/busybox:latest


COPY              restmonkey-linux-amd64             /bin/restmonkey
EXPOSE            8080 
ENTRYPOINT        [  "/bin/restmonkey"  ]