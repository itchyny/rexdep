BIN = rexdep

all: clean test build

build: deps
	go build -o build/$(BIN) .

install: deps
	go install

cross: deps
	goxc -max-processors=8 -build-ldflags="" \
	    -os="linux darwin freebsd netbsd windows" -arch="386 amd64 arm" -d . \
	    -resources-include='README*' -n $(BIN)

deps:
	go get -d -v .

test: testdeps build
	go test -v ./...

testdeps:
	go get -d -v -t .

LINT_RET = .golint.txt
lint: lintdeps build
	go vet
	rm -f $(LINT_RET)
	golint ./... | tee .golint.txt
	test ! -s $(LINT_RET)

lintdeps:
	go get -d -v -t .
	go get -u github.com/golang/lint/golint

clean:
	rm -rf build snapshot debian
	go clean

.PHONY: build cross deps test testdeps lint lintdeps clean
