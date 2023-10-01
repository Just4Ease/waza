local:
	go run main.go
schema:
	go generate ./...
	make alignment
fmt:
	go fmt
alignment:
	go run golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment -fix ./models > /dev/null 2>&1 || :

dev:
	nodemon --exec "go run" main.go --signal SIGTERM -e .js,.css,.go,.html,.env

build:
	docker build -t waza-justice-test .
docker: build
	docker run -e ENVIRONMENT=staging -e PORT=4000 waza-justice-test
