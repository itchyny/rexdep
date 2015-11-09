BIN = rexdep

all: clean test build

test: build
	go test -v ./...

build: deps
	go build -o build/$(BIN) .

install: deps
	go install

cross: deps
	goxc -os="linux darwin freebsd windows" -arch="386 amd64 arm" -d . \
	    -resources-include='README*' -n $(BIN)

deps:
	go get -d -v .

testdeps:
	go get -d -v -t .

clean:
	rm -rf build snapshot debian
	go clean

.PHONY: test build cross deps testdeps clean
