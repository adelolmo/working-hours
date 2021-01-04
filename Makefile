MAKEFLAGS += --silent
BINDIR=$(DESTDIR)/usr/bin

compile:
	@echo Compiling...
	go build -o wh

install: compile
	@echo Installing...
	install -Dm755 wh $(BINDIR)/wh
	rm -rf wh

uninstall:
	@echo Uninstalling...
	rm -rf $(BINDIR)/wh
