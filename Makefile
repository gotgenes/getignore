DIST_DIRS := find dist -depth 1 -type d -execdir
VERSION := $(patsubst v%,%,$(shell git describe --tags))
PACKAGES := $(shell glide novendor)

build:
	go build

install:
	go install

test:
	go test $(PACKAGES)
	go vet $(PACKAGES)

clean:
	rm -f ./getignore
	rm -rf ./dist

build-all:
	gox \
	-osarch="darwin/amd64 linux/386 linux/amd64 linux/arm linux/arm64 windows/amd64" \
	-output="dist/getignore-${VERSION}-{{.OS}}-{{.Arch}}/{{.Dir}}" .

dist: build-all
	$(DIST_DIRS) cp -r ../LICENSE ../README.md ../completions {} \; && \
	$(DIST_DIRS) tar -zcf {}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r {}.zip {} \;

.PHONY: build test install clean build-all dist
