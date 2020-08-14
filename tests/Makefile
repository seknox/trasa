.PHONY: dependency unit-test integration-test docker-up docker-down clear

dependency:
	@cd ../server && go get -v ./...

integration-test: docker-up dependency
	@go test -v  ./integration-tests

unit-test: dependency
	@cd ../server && go test -v -short ./...

docker-up:
	@cd ../build/integration-test && docker-compose up -d

docker-down:
	@cd ../build/integration-test && docker-compose down

clear: docker-down
