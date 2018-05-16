$ kubectl run sre-shell --rm -i --tty --image ubuntu:18.04 -- bash



```
$ squarectl kubectl-shell -c platform -v reliability
$ kubectl config get-contexts
CURRENT   NAME                                      CLUSTER                                   AUTHINFO                               NAMESPACE
*         context-platform-reliability-pp-gb-lon1   cluster-platform-reliability-pp-gb-lon1   user-platform-reliability-pp-gb-lon1   reliability-pp


http://identity-pp-proxy.apps-api-k8s-001-test1-mcc-be-gcw1.metroscales.io/

curl 'http://localhost:9090/probe?url=https://idam.metrosystems.net/.well-known/openid-configuration&requests=1000&workers=20'