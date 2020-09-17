.PHONY: dependency unit-test integration-test docker-up-all docker-down-all docker-up-db docker-down-db clear

dependency:
	@cd ../server && go get -v ./...

end-to-end-test: docker-up-all dependency
	@cd e2e && yarn install && yarn jest
	@make clear

integration-test: docker-up-db dependency
	@go test -v  ./server
	@make clear

unit-test: dependency
	@cd ../server && go test -v -short ./...

docker-up-db:
	@cd build/integration && docker-compose up -d

docker-down-db:
	@cd build/integration && docker-compose down

docker-up-all:
	@cd build/e2e && docker-compose up --build -d

docker-down-all:
	@cd build/e2e && docker-compose down

clear: docker-down-all docker-down-db
