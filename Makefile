.PHONY: all, eternal, eventworker, fmt, clean

SRC=./src
BIN=./bin

all: eternal eventworker


eternal:
	GOPATH=`pwd` go build -o $(BIN)/$@ $(SRC)/$@/main.go

eventworker:
	GOPATH=`pwd` go build -o $(BIN)/$@ $(SRC)/eternal/$@/*.go

fmt:
	find $(SRC)/eternal -name '*.go'|xargs gofmt -w

clean:
	rm -rf bin/*
