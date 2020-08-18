.PHONY: dependency unit-test integration-test docker-up-all docker-down-all docker-up-db docker-down-db clear

dependency:
	@cd ../server && go get -v ./...

end-to-end-test: docker-up-all dependency
	@yarn jest
	@make clear

integration-test: docker-up-db dependency
	@go test -v  ./integration-tests
	@make clear

unit-test: dependency
	@cd ../server && go test -v -short ./...

docker-up-db:
	@cd build/integration-test && docker-compose up -d

docker-down-db:
	@cd build/integration-test && docker-compose down

docker-up-all:
	@cd build/end-to-end && docker-compose up -d

docker-down-all:
	@cd build/end-to-end && docker-compose down

clear: docker-down-all docker-down-db
