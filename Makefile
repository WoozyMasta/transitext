GO          ?= go
LINTER      ?= golangci-lint
ALIGNER     ?= betteralign
BENCHSTAT   ?= benchstat
SCHEMADOC   ?= schemadoc
BENCH_COUNT ?= 6
BENCH_REF   ?= bench_baseline.txt
SCHEMA_FILE ?= schema.json
DOC_FILE    ?= CONFIG.md

.PHONY: test test-race test-short bench bench-fast bench-reset verify vet check ci \
	fmt fmt-check lint lint-fix align align-fix tidy tidy-check download \
	tools tools-ci tool-golangci-lint tool-betteralign tool-benchstat tool-schemadoc \
	schema config-doc config-doc-check release-check release-notes

check: verify tidy fmt vet lint-fix align-fix test-short config-doc
ci: download tools-ci verify tidy-check fmt-check vet lint align test-short

fmt:
	gofmt -w .

fmt-check:
	@files=$$(gofmt -l .); \
	if [ -n "$$files" ]; then \
		echo "$$files" 1>&2; \
		echo "gofmt: files need formatting" 1>&2; \
		exit 1; \
	fi

vet:
	$(GO) vet ./...

test:
	$(GO) test ./...

test-race:
	$(GO) test -race ./...

test-short:
	$(GO) test -short ./...

bench:
	@tmp=$$(mktemp); \
	$(GO) test ./... -run=^$$ -bench 'Benchmark' -benchmem -count=$(BENCH_COUNT) | tee "$$tmp"; \
	if [ -f "$(BENCH_REF)" ]; then \
		$(BENCHSTAT) "$(BENCH_REF)" "$$tmp"; \
	else \
		cp "$$tmp" "$(BENCH_REF)" && echo "Baseline saved to $(BENCH_REF)"; \
	fi; \
	rm -f "$$tmp"

bench-fast:
	$(GO) test ./... -run=^$$ -bench 'Benchmark' -benchmem

bench-reset:
	rm -f "$(BENCH_REF)"

verify:
	$(GO) mod verify

tidy-check:
	@$(GO) mod tidy
	@git diff --stat --exit-code -- go.mod go.sum || ( \
		echo "go mod tidy: repository is not tidy"; \
		exit 1; \
	)

tidy:
	$(GO) mod tidy

download:
	$(GO) mod download

lint:
	$(LINTER) run ./...

lint-fix:
	$(LINTER) run --fix ./...

align:
	$(ALIGNER) ./...

align-fix:
	-$(ALIGNER) -apply ./...
	$(ALIGNER) ./...

tools: tool-golangci-lint tool-betteralign tool-benchstat tool-schemadoc
tools-ci: tool-golangci-lint tool-betteralign tool-schemadoc

tool-golangci-lint:
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

tool-betteralign:
	$(GO) install github.com/dkorunic/betteralign/cmd/betteralign@latest

tool-benchstat:
	$(GO) install golang.org/x/perf/cmd/benchstat@latest

tool-schemadoc:
	$(GO) install github.com/woozymasta/schemadoc/cmd/schemadoc@latest

release-notes:
	@awk '\
	/^<!--/,/^-->/ { next } \
	/^## \[[0-9]+\.[0-9]+\.[0-9]+\]/ { if (found) exit; found=1; next } \
	found { \
		if (/^## \[/) { exit } \
		if (/^$$/) { flush(); print; next } \
		if (/^\* / || /^- /) { flush(); buf=$$0; next } \
		if (/^###/ || /^\[/) { flush(); print; next } \
		sub(/^[ \t]+/, ""); sub(/[ \t]+$$/, ""); \
		if (buf != "") { buf = buf " " $$0 } else { buf = $$0 } \
		next \
	} \
	function flush() { if (buf != "") { print buf; buf = "" } } \
	END { flush() } \
	' CHANGELOG.md

schema:
	@echo ">> schema: $(SCHEMA_FILE)"
	$(SCHEMADOC) mod2schema --module-root . \
		--package github.com/woozymasta/transitext/config \
		--type Config github.com/woozymasta/transitext $(SCHEMA_FILE)

config-doc: schema
	@echo ">> config doc: $(DOC_FILE)"
	$(SCHEMADOC) schema2md --title "transitext config reference" --template table \
		--list-marker '*' --mode all --format yaml \
		"$(SCHEMA_FILE)" "$(DOC_FILE)"

config-doc-check: config-doc
	@git diff --stat --exit-code -- "$(SCHEMA_FILE)" "$(DOC_FILE)" || ( \
		echo "schema/docs are out of date; run 'make config-doc' and commit changes"; \
		exit 1; \
	)
