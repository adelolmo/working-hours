MAKEFLAGS += --silent
BINDIR=$(DESTDIR)/usr/bin
BIN=wh

.PHONY: all build test tidy vendor install uninstall

all: build

build:
	go build -o $(BIN)

test:
	go clean -testcache
	go test ./...

tidy:
	go mod tidy

vendor: tidy
	go mod vendor

install:
	install -Dm755 $(BIN) $(BINDIR)/$(BIN)
	install wh-completion.bash $(DESTDIR)/etc/bash_completion.d/

uninstall:
	rm -rf $(BINDIR)/$(BIN)
	rm -rf $(DESTDIR)/etc/bash_completion.d/wh-completion.bash
