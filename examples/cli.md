go run cmd/ustress/main.go stress --url=https://reliability-pp.metrosystems.net/docs/ --requests=10000 --workers=20 


go run cmd/ustress/main.go stress --url=http://restmonkey.com/restmonkey/api/v1/test --requests=20 --workers=4 --resolve=127.0.0.1:8080