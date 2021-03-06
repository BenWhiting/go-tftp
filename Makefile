GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test -v
LDFLAGS=-ldflags="-s -w"
GOBUILDFLAGS=-a $(LDFLAGS)
BUILDDIR=build

all:  build

build:
	CGO_ENABLE=0 $(GOBUILD) $(GOBUILDFLAGS) -o $(BUILDDIR)/go-tftp cmd/go-tftp.go