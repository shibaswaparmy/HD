
Installation
1. go get github.com/rakyll/statik


Steps to follow
1. cd shibaswaparmy/heimdall/server
2. Update swagger.yaml file inside swagger-ui directory
3. cd shibaswaparmy/heimdall/server && statik -src=./swagger-ui
4. cd shibaswaparmy/heimdall && make build
5. cd shibaswaparmy/heimdall && make run-server

Visit http://localhost:1317/swagger-ui/ 


Reference
- https://github.com/rakyll/statik