GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test

.PHONY: cover
cover:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out -o coverage.html