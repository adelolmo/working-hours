MAKEFLAGS += --silent
BINDIR=$(DESTDIR)/usr/bin
BIN=wh

.PHONY: test tidy vendor install uninstall

$(BIN): test
	@echo Compiling...
	go build -o $(BIN)

test:
	go clean -testcache
	go test ./...

tidy:
	go mod tidy

vendor: tidy
	go mod vendor

install: $(BIN)
	@echo Installing...
	install -Dm755 $(BIN) $(BINDIR)/$(BIN)
	rm -rf $(BIN)

uninstall:
	@echo Uninstalling...
	rm -rf $(BINDIR)/$(BIN)
