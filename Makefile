BIN = rexdep

all: clean test build

build: deps
	go build -o build/$(BIN) .

install: deps
	go install

cross: deps
	goxc -build-ldflags="" -os="linux darwin freebsd windows" -arch="386 amd64 arm" -d . \
	    -resources-include='README*' -n $(BIN)

deps:
	go get -d -v .

test: testdeps build
	go test -v ./...

testdeps:
	go get -d -v -t .

clean:
	rm -rf build snapshot debian
	go clean

.PHONY: build cross deps test testdeps clean
