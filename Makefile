# Defines shell to bash when using zsh
SHELL := /bin/bash

help: ## This help message
	@echo -e "$$(grep -hE '^\S+:.*##' $(MAKEFILE_LIST) | sed -e 's/:.*##\s*/:/' -e 's/^\(.\+\):\(.*\)/\\x1b[36m\1\\x1b[m:\2/' | column -c2 -t -s :)"

unit-tests: ## Run unit tests
	go test -race ./...
	go test -v -coverprofile=coverage.out ./...

view-tests-report: ## View HTML test report on firefox
	@echo Generating HTML report...
	@go tool cover -html=coverage.out -o coverage.html
	@echo Opening file on firefox...
	@firefox coverage.html

lint: ## Lint your code
	./bin/golangci-lint run -v

setup: ## Setup binary packages for application
	wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.24.0