# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
BINARY_NAME=gogrep

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

build-static:
	CGO_ENABLED=0 $(GOBUILD) -o $(BINARY_NAME) -a -ldflags '-extldflags "-static"' -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	./$(BINARY_NAME)

.PHONY: all build build-static test clean run
