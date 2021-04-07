MAKEFLAGS += --silent
BINDIR=$(DESTDIR)/usr/bin

wh:
	@echo Compiling...
	go build -o wh

install: wh
	@echo Installing...
	install -Dm755 wh $(BINDIR)/wh
	rm -rf wh

uninstall:
	@echo Uninstalling...
	rm -rf $(BINDIR)/wh
