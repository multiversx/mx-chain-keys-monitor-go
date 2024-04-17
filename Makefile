test:
	@echo "  >  Running unit tests"
	go test -cover -race -coverprofile=coverage.txt -covermode=atomic -v ./...

lint-install:
ifeq (,$(wildcard test -f bin/golangci-lint))
	@echo "Installing golint"
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s
endif

run-lint:
	@echo "Running golint"
	bin/golangci-lint run --max-issues-per-linter 0 --max-same-issues 0 --timeout=2m

lint: lint-install run-lint

cli-docs:
	cd ./cmd && bash ./CLI.md.sh

check-cli-md:
	cd ./cmd/monitor && go build
	cd ./cmd && bash ./CLI.md.sh
	@status=$$(git status --porcelain | grep CLI); \
    	if [ ! -z "$${status}" ]; \
    	then \
    		echo "Error - please update all CLI.md files by running the 'cli-docs' or 'check-cli-md' from Makefile!"; \
    		exit 1; \
    	fi