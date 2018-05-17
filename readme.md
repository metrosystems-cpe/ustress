# Traffic Monkey

it is Web application designed to be deployed in various kubernetes clusters from where to start send traffic to an endpoint.

urlHandlers:
 - /probe          - > send a get request and the monkey will start the atack 
 - /data           - > will expose 
 - /.well-known/*  - > live / ready / metrics (no metrics)



## Current deployment process 

in order to make it work, you need 2 terminal windows and run a lot of commands :) 

Warning: if --rm flag (default) is set and all tty's to the kubernetes pod are interupted the pod will be destroyed and with him all the report data


```
Terminal 1
step 1 squarectl kubectl-shell -c mcc -v ccms && kubectl config get-contexts && kubectl config use-context <context>
step 2 kubectl run sre-shell --rm -i --tty --image ubuntu:18.04 -- bash
step 3 apt update && apt install -y curl
step 7 ./traefikmonkey

Terminal 2
step 4 squarectl kubectl-shell -c mcc -v ccms && kubectl config get-contexts && kubectl config use-context <context> && kubectl get po
step 5 kubectl cp ./trafficmonkey-linux-amd64 sre-shell-300457757-xk08h:/trafficmonkey
step 6 kubectl exec -i -t sre-shell-300457757-xk08h bash
step 8 curl 'http://localhost:9090/probe?url=https://idam.metrosystems.net/.well-known/openid-configuration&requests=1000&workers=10'
step 8 curl 'http://localhost:9090/probe?url=http://proxy.identity-prod:80/.well-known/openid-configuration&requests=1000&workers=10'
step 8 curl 'http://localhost:9090/probe?url=http://proxy-k8s-001-live1-mcc-gb-lon1.metroscales.io:30021/.well-known/openid-configuration&requests=1000&workers=10'
step 9 kubectl cp sre-shell-300457757-xk08h:/data/ ./data

```

## analyze generated data. 

it was designed to generate percentile statistics and not detect 503 ( application code for timout error ) :) 

### analyze-data.sh
```
find . -type f -name "*.json" | while read file; do \
printf " >>> file: %s\n" $file; \
jq '."url"' $file; \
jq '."timestamp"' $file ; \
jq '."stats"' $file; \
jq '.["data"]' $file | grep status | sort | uniq -c; done
```