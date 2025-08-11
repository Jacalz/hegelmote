APPID = io.github.jacalz.hegelmote
NAME = hegelmote

# If PREFIX isn't provided, default to /usr.
PREFIX ?= /usr

release:
	go build -tags no_emoji,no_metadata,migrated_fynedo -trimpath -ldflags="-s -w" -buildvcs=false -o $(NAME)
.PHONY: release

debug:
	go build -tags no_emoji,no_metadata -trimpath -o $(NAME)
.PHONY: debug

install:
	install -Dm00755 $(NAME) $(DESTDIR)$(PREFIX)/bin/$(NAME)
	install -Dm00644 assets/img/icon.svg $(DESTDIR)$(PREFIX)/share/icons/hicolor/scalable/apps/$(APPID).svg
	install -Dm00644 assets/img/icon-512.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/512x512/apps/$(APPID).png
	install -Dm00644 assets/img/icon-256.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/256x256/apps/$(APPID).png
	install -Dm00644 assets/img/icon-128.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/128x128/apps/$(APPID).png
	install -Dm00644 assets/img/icon-64.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/64x64/apps/$(APPID).png
	install -Dm00644 assets/img/icon-48.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/48x48/apps/$(APPID).png
	install -Dm00644 assets/img/icon-32.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/32x32/apps/$(APPID).png
	install -Dm00644 assets/img/icon-24.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/24x24/apps/$(APPID).png
	install -Dm00644 assets/img/icon-16.png $(DESTDIR)$(PREFIX)/share/icons/hicolor/16x16/apps/$(APPID).png
	install -Dm00644 assets/xdg/$(APPID).desktop $(DESTDIR)$(PREFIX)/share/applications/$(APPID).desktop
	install -Dm00644 assets/xdg/$(APPID).appdata.xml $(DESTDIR)$(PREFIX)/share/appdata/$(APPID).appdata.xml
	# NOTE: You might want to update your gtk icon cache by running `make update-icon-cache` afterwards.
	# Not doing this might result in the application not showing up in the application menu.
.PHONY: install

update-icon-cache:
	sudo gtk-update-icon-cache -f /usr/share/icons/hicolor/
.PHONY: update-icon-cache

uninstall:
	-rm $(DESTDIR)$(PREFIX)/bin/$(NAME)
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/512x512/apps/$(APPID).png
	-rm $(DESTDIR)$(PREFIX)/share/icons/hicolor/scalable/apps/$(APPID).svg
	-rm $(DESTDIR)$(PREFIX)/share/applications/$(APPID).desktop
	-rm $(DESTDIR)$(PREFIX)/share/appdata/$(APPID).appdata.xml
.PHONY: uninstall

wasm:
	rm -rf wasm cmd/webmote/wasm
	~/go/bin/fyne package -os wasm -release -tags no_emoji
	cp assets/img/favicon.png wasm/icon.png
	mv wasm cmd/webmote
.PHONY: wasm

wasm-opt: wasm
	wasm-opt cmd/webmote/wasm/Hegelmote.wasm --enable-bulk-memory-opt -O4 -o cmd/webmote/wasm/Hegelmote.wasm
.PHONY: wasm-opt
