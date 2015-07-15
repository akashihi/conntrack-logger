VERSION=0.1.3

# Where to install to.
PREFIX?=/usr/bin

default: build-all
build-all: deb

.PHONY: go-check
go-check:
	@go version > /dev/null || (echo "Go not found. You need to install go: http://golang.org/doc/install"; false)
	@go version | grep -q 'go version go1.3' || (echo "Go version 1.3 required, you have a version of go that is unsupported. See http://golang.org/doc/install"; false)


clean:
	-@rm -fr build bin pkg

deps-clean:
	rm -fr src/code.google.com/
	rm -fr src/github.com

conntrack-logger: go-check 
	go build .

deb: conntrack-logger
	fpm -s dir -t deb -n conntrack-logger -v ${VERSION} conntrack-logger=/usr/bin/