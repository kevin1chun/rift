GO ?= `which go`
SCRIPTPATH ?= $(shell pwd)
POINTLANDER = $(SCRIPTPATH)/bin/peg

default: gengrammar build

gengrammar:
	$(POINTLANDER) -inline -switch $(SCRIPTPATH)/src/rift/lang/rift.g

getdeps:
	GOPATH=$(SCRIPTPATH) $(GO) get github.com/pointlander/peg

build:
	GOPATH=$(SCRIPTPATH) $(GO) build -v -o $(SCRIPTPATH)/bin/rift

clean:
	@rm -rf $(SCRIPTPATH)/bin