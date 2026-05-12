.PHONY: *

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Builds and tags the app Docker container
	docker build --tag koreader-sync .
run: ## Runs the build, then runs the Docker container on port :8080
	$(MAKE) build
	docker run -e LOG_LEVEL=DEBUG -p 8080:8080 -v ./data:/app/data koreader-sync
test: ## Runs the tests for the application
	go test ./... -coverprofile testcoverage.txt
	go tool cover -html testcoverage.txt -o testcoverage.html