install-lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.34.1
lint-check:
	@echo "Performing lint check"
	@golangci-lint --out-format line-number -D gosimple -E gofumpt,goimports,gosec run ./...

test-unit:
	@echo "Run unit tests"
	go test -v ./...
