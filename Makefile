GO ?= `which go`
SCRIPTPATH ?= $(shell pwd)
POINTLANDER = /home/jl/dev/lib/go/bin/peg

default: clean gengrammar build

gengrammar:
	$(POINTLANDER) -inline -switch $(SCRIPTPATH)/src/rift/lang/rift.g

build:
	GOPATH=$(SCRIPTPATH) $(GO) build -v -o $(SCRIPTPATH)/bin/rift

clean:
	@rm -rf $(SCRIPTPATH)/bin