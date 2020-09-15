.PHONY: dependency unit-test integration-test docker-up-all docker-down-all docker-up-db docker-down-db clear


dev-setup:
	@mkdir -p devenv/var/log && cp -r build/etc devenv/ && touch devenv/var/log/trasa.log

dev-run:
	@cd build/docker && docker-compose up --build

clear:
	@cd build/docker && docker-compose down
