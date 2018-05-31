FROM quay.io/prometheus/busybox:latest


COPY              restmonkey-linux-amd64             /restmonkey
COPY              ui                                 /ui
EXPOSE            8080 
ENTRYPOINT        [  "/restmonkey"  ]