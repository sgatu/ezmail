SHELL := /bin/bash
_build-api:
	@rm -rf ./dist/api
	@mkdir ./dist/api
	go build -o ./dist/api/ezmail ./cmd/api/main.go
	@chmod +x ./dist/api/ezmail
_build-executor:
	@rm -rf ./dist/executor
	@mkdir ./dist/executor
	go build -o ./dist/executor/executor ./cmd/consumer/main.go
	@chmod +x ./dist/executor/executor
build-api:
	@$(MAKE) -s _build-api
build-executor:
	@$(MAKE) -s _build-executor
run:
	@$(MAKE) -s _build-api
	@cp ./.env.local ./dist/api/.env
	@cd ./dist/api/ && export $$(grep -v '^#' .env | xargs) && ./ezmail
run-exec:
	@$(MAKE) -s _build-executor
	@cp ./.env.local ./dist/executor/.env
	@cd ./dist/executor && export $$(grep -v '^#' .env | xargs) && ./executor
test:
	go test ./.../test -count=1 -coverpkg=./... -coverprofile=/tmp/cover.out
cover:
	@$(MAKE) -s test
	go tool cover -func=/tmp/cover.out 
cover-short:
	@$(MAKE) -s test > /dev/null
	@echo "Total coverage: " $$(go tool cover -func=/tmp/cover.out | grep "total:"  | awk '{print $$NF}')

