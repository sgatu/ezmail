_build-api:
	@rm -rf ./dist
	@mkdir ./dist
	@cp ./.env.local ./dist/.env
	go build -o ./dist/ezmail ./cmd/api/main.go
	@chmod +x ./dist/ezmail
build-api:
	@$(MAKE) -s _build-api
run:
	@$(MAKE) -s _build-api
	@cd ./dist/ && ./ezmail
