PKG = github.com/gopheracademy/manager
CMD = alltag
PREFIX = /usr

all: build/$(CMD)

################################################################################
# building/bundling CSS/JS artifacts

hash-fonts: FORCE
	bash ./src/hash-artifacts.sh build/fonts-by-hash src/*.otf
build/fonts-by-hash/fonts.scss: hash-fonts
	sh ./src/render-font-css.sh > $@

build/alltag.css: src/alltag.scss | build/fonts-by-hash/fonts.scss
	sassc -t compressed -I vendor/github.com/majewsky/xyrillian.css -I build/fonts-by-hash $< $@

hashed-artifacts: build/alltag.css src/alltag.js | src/hash-artifacts.sh hash-fonts FORCE
	bash ./src/hash-artifacts.sh build/by-hash $^ src/*.otf
build/bindata/bindata.go: hashed-artifacts
	@mkdir -p build/bindata
	go-bindata -modtime 1 -pkg bindata -prefix build/by-hash/ -o $@ build/by-hash/*

admin: FORCE
	@cd admin && npm run build

################################################################################
# compiling and installing the binary

# NOTE: This repo uses Go modules, and uses a synthetic GOPATH at
# $(CURDIR)/.gopath that is only used for the build cache. $GOPATH/src/ is
# empty.
GO            = GOPATH=$(CURDIR)/.gopath GOBIN=$(CURDIR)/build go
GO_BUILDFLAGS =
GO_LDFLAGS    = -s -w

build/$(CMD): build/bindata/bindata.go
	$(GO) install $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' '$(PKG)'

install: FORCE all
	install -D -m 0755 "build/$(CMD)" "$(DESTDIR)$(PREFIX)/bin/$(CMD)"

################################################################################
# utilities

# convenience target for developers: `make run` runs the application with
# environment options sourced from $PWD/.env
run: build/$(CMD)
	set -euo pipefail && source ./.env && ./build/$(CMD) $*

vendor: FORCE
	$(GO) mod tidy
	$(GO) mod vendor

generate: FORCE
	$(GO) generate
.PHONY: FORCE