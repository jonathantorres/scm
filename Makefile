PROG := scm
PKG := github.com/jonathantorres/scm
PREFIX := $(shell pwd)

# compile program
$(PROG):
	go build -o bin/$(PROG) -ldflags="$(LDFLAGS)" $(PKG)

# create release version
.PHONY: release
release:
	go build -o bin/$(PROG) -ldflags="-s -w $(LDFLAGS)" $(PKG)

# Run tests
.PHONY: test
test:
	go test .

.PHONY: clean
clean:
	go clean
	rm -fr ./bin
	mkdir ./bin && touch ./bin/.gitkeep
