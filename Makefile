.PHONY: eternal, fmt, clean

SRC=./src
BIN=./bin

eternal:
	GOPATH=`pwd` go build -o $(BIN)/$@ $(SRC)/$@/main.go

fmt:
	find $(SRC)/eternal -name '*.go'|xargs gofmt -w

clean:
	rm -rf bin/*
