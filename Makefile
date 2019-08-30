GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
LDFLAGS=-ldflags="-s -w"
GOBUILDFLAGS=-a $(LDFLAGS)
BUILDDIR=build

all: test build

test: 
	$(GOTEST) ./...

build:
	CGO_ENABLE=0 $(GOBUILD) $(GOBUILDFLAGS) -o $(BUILDDIR)/go-tftp cmd/go-tftp.go