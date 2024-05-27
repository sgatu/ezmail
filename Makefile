_build:
	@rm -rf ./dist
	@mkdir ./dist
	go build -o ./dist/ezmail
	@chmod +x ./dist/ezmail
build:
	@$(MAKE) -s _build
run:
	@$(MAKE) -s _build
	./dist/ezmail
