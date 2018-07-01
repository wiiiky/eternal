.PHONY: all, eternal, eworker, fmt, clean

SRC=./src
BIN=./bin

all: eternal eworker


eternal:
	GOPATH=`pwd` go build -o $(BIN)/$@ $(SRC)/$@/main.go

eworker:
	GOPATH=`pwd` go build -o $(BIN)/$@ $(SRC)/eternal/$@/*.go

fmt:
	find $(SRC)/eternal -name '*.go'|xargs gofmt -w

clean:
	rm -rf bin/*
