CMDDIR := ./cmd
OBJDIR := ./output

# Build commands
.PHONY: build
build: $(OBJDIR)/tran

$(OBJDIR)/%: $(CMDDIR)/%/main.go deps
	go build -o $@ $<

# Clean commands
.PHONY: clean
clean:
	go mod tidy
	go clean -testcache
	rm -rf $(OBJDIR)

# Install commands
.PHONY: install
install:
	go install ./...

# Uninstall commands
.PHONY: uninstall
uninstall:
	go clean -i ./...

# Lint
.PHONY: lint
lint: devdeps
	go vet ./...
	golint -set_exit_status -min_confidence 0 ./...

# Run tests
.PHONY: test
test: deps dir
	go test -v -coverprofile=$(OBJDIR)/cover.out  ./...
	go tool cover -html=$(OBJDIR)/cover.out -o $(OBJDIR)/cover.html

#Install dependencies
.PHONY: deps
deps:
	go get github.com/mattn/go-isatty
	go get github.com/morikuni/aec
	go get github.com/peterh/liner

.PHONY: devdeps
devdeps: deps
	go get golang.org/x/lint

#Make directory
.PHONY: dir
dir:
	@if [ ! -d $(OBJDIR) ]; \
		then echo "mkdir -p $(OBJDIR)"; mkdir -p $(OBJDIR); \
	fi

