PHONY: gen-mocks
gen-mocks:
	go generate ./...

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: test-report
test-report:
	go test ./... -coverprofile=c.out
	go tool cover -html=c.out -o test_coverage.html