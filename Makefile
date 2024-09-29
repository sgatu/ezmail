_build-api:
	@rm -rf ./dist/api
	@mkdir ./dist/api
	@cp ./.env.local ./dist/api/.env
	go build -o ./dist/api/ezmail ./cmd/api/main.go
	@chmod +x ./dist/api/ezmail
_build-executor:
	@rm -rf ./dist/executor
	@mkdir ./dist/executor
	@cp ./.env.local ./dist/executor/.env
	go build -o ./dist/executor/executor ./cmd/consumer/main.go
	@chmod +x ./dist/executor/executor
build-api:
	@$(MAKE) -s _build-api
build-executor:
	@$(MAKE) -s _build-executor
run:
	@$(MAKE) -s _build-api
	@cd ./dist/api/ && ./ezmail
run-exec:
	@$(MAKE) -s _build-executor
	@cd ./dist/executor && ./executor
