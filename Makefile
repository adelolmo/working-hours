MAKEFLAGS += --silent
BINDIR=$(DESTDIR)/usr/bin
BIN=wh

$(BIN): test
	@echo Compiling...
	go build -o $(BIN)

test:
	go clean -testcache
	go test ./...

install: $(BIN)
	@echo Installing...
	install -Dm755 $(BIN) $(BINDIR)/$(BIN)
	rm -rf $(BIN)

uninstall:
	@echo Uninstalling...
	rm -rf $(BINDIR)/$(BIN)
