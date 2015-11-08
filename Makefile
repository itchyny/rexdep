BIN = rexdep
VERSION = v0.0.0
DESCRIPTION = Roughly extract dependency from source code
AUTHOR = itchyny

BUILD_FLAGS = "\
	    -X main.Author \"$(AUTHOR)\" \
	    -X main.Description \"$(DESCRIPTION)\" \
	    -X main.Name \"$(BIN)\" \
	    -X main.Version \"$(VERSION)\" \
	    "

all: clean test build

test: build
	go test -v ./...

build: deps
	go build -ldflags=$(BUILD_FLAGS) -o build/$(BIN) .

install: deps
	go install -ldflags=$(BUILD_FLAGS)

cross: deps
	goxc -build-ldflags=$(BUILD_FLAGS) \
	    -os="linux darwin freebsd windows" -arch="386 amd64 arm" -d . \
	    -resources-include='README*' -n $(BIN)

deps:
	go get -d -v .

testdeps:
	go get -d -v -t .

clean:
	rm -rf build snapshot debian
	go clean

.PHONY: test build cross deps testdeps clean
