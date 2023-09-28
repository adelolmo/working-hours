#MAKEFLAGS += --silent
BINDIR=$(DESTDIR)/usr/bin
BIN=wh

.PHONY: all build test tidy vendor install uninstall clean

all: $(BIN) $(BIN).1.gz

$(BIN):
	go build -o $(BIN)

$(BIN).1.gz:
	go-md2man -in man.1.md -out $(BIN).1
	gzip $(BIN).1

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
	install -Dm644 $(BIN).1.gz $(DESTDIR)/usr/share/man/man1/$(BIN).1.gz

uninstall:
	rm -rf $(BINDIR)/$(BIN)
	rm -rf $(DESTDIR)/etc/bash_completion.d/wh-completion.bash
	rm -f $(DESTDIR)/usr/share/man/man1/$(BIN).1.gz

clean:
	rm $(BIN) $(BIN).1.gz