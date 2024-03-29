.PHONY: all clean deps build run dev watch

# locations
#
B = $(CURDIR)/build
GOBIN = $(GOPATH)/bin
NPX_BIN = $(CURDIR)/node_modules/.bin

BINDIR ?= $(GOBIN)

# config files
#
MODD_RUN_CONF = $(B)/modd-run.conf
MODD_DEV_CONF = $(B)/modd-dev.conf
MODD_WATCH_CONF = $(B)/modd-watch.conf

# tools
#
GO = go
GOFMT = gofmt
GOFMT_FLAGS = -w -l -s
GOBUILD = $(GO) build
GOBUILD_FLAGS = -v
GOGET = $(GO) get
GOGET_FLAGS = -v
GOINSTALL = $(GO) install
GOINSTALL_FLAGS = -v
GOGENERATE_FLAGS = -v
NPM = npm

MODD = $(GOBIN)/modd
MODD_FLAGS = -b
WEBPACK = $(NPX_BIN)/webpack

FILE2GO_URL = go.sancus.dev/file2go/cmd/file2go
MODD_URL = github.com/cortesi/modd/cmd/modd@latest

FILE2GO = $(GO) run $(FILE2GO_URL)

# magic constants
#
MOD = $(shell sed -n -e 's/^module \(.*\)/\1/p' go.mod)
PORT ?= 8080
DEV_PORT ?= 8081

# generated files
#
ASSETS_GO_FILE = web/assets/files.go
HTML_GO_FILE = web/html/files.go

GENERATED_GO_FILES = $(ASSETS_GO_FILE) $(HTML_GO_FILE)

# default target
all: build

# deps
#
.PHONY: go-deps npm-deps
deps: go-deps npm-deps

# go-deps
GO_DEPS = $(MODD)

go-deps: $(GO_DEPS)
	$(GOGET) $(GOGET_FLAGS) ./...

$(MODD): URL=$(MODD_URL)

$(GO_DEPS):
	$(GOINSTALL) $(GOINSTALL_FLAGS) $(URL)

# npm-deps
NPM_DEPS = $(WEBPACK)

$(NPXBIN)/%:
	$(NPM) i
	$(NPM) shrinkwrap

npm-deps:
	$(NPM) i
	$(NPM) shrinkwrap

# clean
#
clean:
	$(GO) clean -x -r -modcache
	git ls-files -o $(dir $(ASSETS_GO_FILE)) | xargs -rt rm
	rm -rf $(B) node_modules/
	rm -f $(GENERATED_GO_FILES) npm-shrinkwrap.json package-lock.json

# fmt
#
.PHONY: fmt lint npm-lint go-fmt

fmt: go-fmt npm-lint
lint: go-fmt npm-lint

go-fmt: $(GO_DEPS) FORCE
	$(GO) mod tidy -v || true
	@find -name '*.go' | xargs -r $(GOFMT) $(GOFMT_FLAGS)

npm-lint: $(NPM_DEPS) FORCE
	$(NPM) run lint

# run
#
MODD_CONF_FILES = $(MODD_RUN_CONF) $(MODD_DEV_CONF) $(MODD_WATCH_CONF)

.PHONY: modd-conf

modd-conf: $(MODD_CONF_FILES)

# TODO: rework these using patterns
$(MODD_RUN_CONF): MODE=run
$(MODD_RUN_CONF): src/modd/run.conf

$(MODD_DEV_CONF): MODE=dev
$(MODD_DEV_CONF): src/modd/dev.conf

$(MODD_WATCH_CONF): MODE=watch
$(MODD_WATCH_CONF): src/modd/watch.conf

$(MODD_CONF_FILES): Makefile
$(MODD_CONF_FILES):
	@mkdir -p $(@D)
	@sed \
		-e "s|@@BINDIR@@|$(BINDIR)|g" \
		-e "s|@@NPM@@|$(NPM)|g" \
		-e "s|@@GO@@|$(GO)|g" \
		-e "s|@@GOFMT@@|$(GOFMT) $(GOFMT_FLAGS)|g" \
		-e "s|@@GOGET@@|$(GOGET) $(GOGET_FLAGS)|g" \
		-e "s|@@GOBUILD@@|$(GOBUILD) $(GOBUILD_FLAGS)|g" \
		-e "s|@@FILE2GO@@|$(FILE2GO)|g" \
		-e "s|@@MODE@@|$(MODE)|g" \
		$< > $@~
	@mv $@~ $@
	@echo ${@F} updated.

run: $(MODD_RUN_CONF)
dev: $(MODD_DEV_CONF)
watch: $(MODD_WATCH_CONF)

run dev watch: $(MODD) go-deps $(NPM_DEPS)
run dev watch:
	env PORT=$(PORT) BACKEND=$(DEV_PORT) $(MODD) $(MODD_FLAGS) -f $<

# build
#
ASSETS_FILES_FILTER = find $(dir $(ASSETS_GO_FILE)) -type f -a ! -name '*.go' -a ! -name '.*' -a ! -name '*~'
HTML_FILES_FILTER = find $(dir $(HTML_GO_FILE)) -type f -name '*.gohtml' -o -name '*.html'
NPM_FILES_FILTER = find src/ -name '*.js' -o -name '*.scss'
GO_FILES_FILTER = find */ -name node_modules -prune -name '*.go'

ASSETS_FILES = $(shell set -x; $(ASSETS_FILES_FILTER))
HTML_FILES = $(shell set -x; $(HTML_FILES_FILTER))
NPM_FILES = $(shell set -x; $(NPM_FILES_FILTER))
GO_FILES = $(shell set -x; $(GO_FILES_FILTER)) $(GENERATED_GO_FILES)

.PHONY: npm-build go-generate go-build

build: go-build

# npm-build
NPM_BUILT_MARK = $(B)/.npm-built

$(NPM_BUILT_MARK): $(NPM_FILES) $(NPM_DEPS) Makefile
	@$(NPM) run build
	@mkdir -p $(@D)
	@touch $@

npm-build: $(NPM_DEPS) FORCE
	@$(NPM) run build
	@mkdir -p $(dir $(NPM_BUILT_MARK))
	@touch $(NPM_BUILT_MARK)

.INTERMEDIATE: $(NPM_BUILT_MARK)

# go-build
$(ASSETS_GO_FILE): $(NPM_BUILT_MARK) $(ASSETS_FILES)
	$(ASSETS_FILES_FILTER) | sort -uV | sed -e 's|^$(@D)/||' | (cd $(@D) && xargs -t $(FILE2GO) -p assets -o $(@F))

$(HTML_GO_FILE): $(HTML_FILES) Makefile
	$(HTML_FILES_FILTER) | sort -uV | sed -e 's|^$(@D)/||' | (cd $(@D) && xargs -t $(FILE2GO) -p html -T html -o $(@F))

go-build: $(GO_FILES) $(GO_DEPS) FORCE
	$(GOBUILD) $(GOBUILD_FLAGS) ./...

go-generate: $(GO_FILES) $(GO_DEPS)
	@git grep -l '^//go:generate' | sed -n -e 's|\(.*\)/[^/]\+\.go$$|\1|p' | sort -u | while read d; do \
		git grep -l '^//go:generate' "$$d" | grep '\.go$$' | xargs -r $(GO) generate $(GOGENERATE_FLAGS); \
	done

# FORCE
.PHONY: FORCE
