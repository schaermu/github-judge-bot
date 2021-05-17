GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

.PHONY: testcov
testcov:
	$(GOTEST) -v -coverprofile=coverage.out -coverpkg=./... ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out -o coverage.html