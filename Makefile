DIST_DIRS := find dist -type d -execdir
VERSION := $(shell git describe --tags)

build:
	go build

install:
	go install

test:
	go test -v .

clean:
	rm -f ./getignore
	rm -rf ./dist

build-all:
	gox -os="linux darwin windows" \
	-arch="amd64 386" \
	-output="dist/getignore-${VERSION}-{{.OS}}-{{.Arch}}/{{.Dir}}" .

dist: build-all
	$(DIST_DIRS) cp LICENSE README.md {} \; && \
	$(DIST_DIRS) tar -zcf {}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r {}.zip {} \;

.PHONY: build test install clean build-all dist
