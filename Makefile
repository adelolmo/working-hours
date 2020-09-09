MAKEFLAGS += --silent
BINDIR=$(HOME)/bin

install:
	@echo Installing...
	go build -o working-hours
	install -Dm755 working-hours $(DESTDIR)$(BINDIR)/working-hours
	rm -rf working-hours

uninstall:
	@echo Uninstalling...
	rm -rf $(DESTDIR)$(BINDIR)/working-hours
