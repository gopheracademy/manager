PKG = github.com/gopheracademy/manager
CMD = manager

all: $(CMD)

################################################################################
# building/bundling CSS/JS artifacts
# BRIAN: These targets don't work - they're from an example on the internet
#
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


# build the cms
content: FORCE
	@cd content && npm install && npm run build

# build the website
www: FORCE
	@cd www && npm install && npm run build

################################################################################
# compiling and installing the binary

# NOTE: This repo uses Go modules, and uses a synthetic GOPATH at
# $(CURDIR)/.gopath that is only used for the build cache. $GOPATH/src/ is
# empty.
GO            = GOPATH=$(CURDIR)/.gopath GOBIN=$(CURDIR)/build go
GO_BUILDFLAGS =
GO_LDFLAGS    = -s -w

$(CMD):  deps generate www
	$(GO) build $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' '$(PKG)'


################################################################################
# utilities

# convenience target for developers: `make run` runs the application with
# environment options sourced from $PWD/.env
run:$(CMD)
	set -euo pipefail && ./$(CMD) $*

vendor: FORCE
	$(GO) mod tidy
	$(GO) mod vendor

# generate all the service files from the oto definitions in the def directory
generate: FORCE
	@./generate.sh
.PHONY: FORCE

deps: FORCE
	@go install github.com/pacedotdev/oto
.PHONY: FORCE
