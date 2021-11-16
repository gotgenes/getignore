DIST_DIRS := find dist -depth 1 -type d -execdir
VERSION := $(patsubst v%,%,$(shell git describe --tags))
LDFLAGS := -ldflags "-X 'github.com/gotgenes/getignore/pkg/getignore.Version=${VERSION}'"

build:
	go build ${LDFLAGS} ./cmd/getignore

install:
	go install ${LDFLAGS} ./cmd/getignore

dev-install:
	go install github.com/onsi/ginkgo/ginkgo@v1.16.5

test:
	go vet ./...
	ginkgo -r ${LDFLAGS}

tag:
	git tag -a -m "Release $(version)" v$(version)

clean:
	rm -f ./getignore
	rm -rf ./dist

build-all:
	gox \
	${LDFLAGS} \
	-osarch="darwin/amd64 darwin/arm64 linux/386 linux/amd64 linux/arm linux/arm64 windows/amd64" \
	-output="dist/getignore-${VERSION}-{{.OS}}-{{.Arch}}/{{.Dir}}" \
	./cmd/getignore

dist: build-all
	$(DIST_DIRS) cp -r ../LICENSE ../README.md ../completions {} \; && \
	$(DIST_DIRS) tar -zcf {}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r {}.zip {} \;

.PHONY: build test install dev-install tag clean build-all dist
