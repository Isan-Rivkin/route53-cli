RUN_CMD=go run main.go
GOCMD=go
GOFMT=$(GOCMD)fmt

run: 
	${RUN_CMD}
debug: 
	${RUN_CMD} -debug
test:
	go test -count=1 ./...
test-verbose:
	go test -v -count=1 ./...

fmt: ## Validate go format
	@echo checking gofmt...
	@res=$$($(GOFMT) -d -e -s $$(find . -type d \( -path ./src/vendor \) -prune -o -name '*.go' -print)); \
	if [ -n "$${res}" ]; then \
		echo checking gofmt fail... ; \
		echo "$${res}"; \
		exit 1; \
	else \
		echo Your code formating is according gofmt standards; \
	fi
lint:
	golangci-lint run
