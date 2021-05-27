.DEFAULT_GOAL := build
BIN_FILE=github-judge

GO=go
GOCOVER=$(GO) tool cover
GOTEST=$(GO) test


run:
	$(GO) run .

build:
	$(GO) build -o "dist/${BIN_FILE}"

clean:
	$(GO) clean
	rm --force cp.out nohup.out

test:
	$(GOTEST) ./...

cover:
	$(GOTEST) -v -coverprofile=coverage.out -coverpkg=./... ./...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out -o coverage.html