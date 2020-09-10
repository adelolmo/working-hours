MAKEFLAGS += --silent
BINDIR=$(HOME)/bin

install:
	@echo Installing...
	go build -o wh
	install -Dm755 wh $(DESTDIR)$(BINDIR)/wh
	rm -rf wh

uninstall:
	@echo Uninstalling...
	rm -rf $(DESTDIR)$(BINDIR)/wh
