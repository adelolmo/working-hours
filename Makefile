MAKEFLAGS += --silent
BINDIR=$(HOME)/bin

wh:
	@echo Compiling...
	go build -o wh

install: wh
	@echo Installing...
	install -Dm755 wh $(DESTDIR)$(BINDIR)/wh
	rm -rf wh

uninstall:
	@echo Uninstalling...
	rm -rf $(DESTDIR)$(BINDIR)/wh
