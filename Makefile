SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY = philosopher

VERSION = $(shell date +%Y%m%d)
BUILD = $(shell  date +%Y%m%d%H%M)

TAG = v2.1.0

LDFLAGS = -ldflags "-w -s -X main.version=${TAG} -X main.build=${BUILD}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} main.go

.PHONY: deploy
deploy:
	unzip -o lib/dat/bindata.go.zip -d  lib/dat/

	unzip -o lib/ext/cdhit/unix/bindata.go.zip -d  lib/ext/cdhit/unix/
	unzip -o lib/ext/cdhit/win/bindata.go.zip -d  lib/ext/cdhit/win/

	unzip -o lib/ext/comet/unix/bindata.go.zip -d  lib/ext/comet/unix/
	unzip -o lib/ext/comet/win/bindata.go.zip -d  lib/ext/comet/win/

	unzip -o lib/ext/interprophet/unix/bindata.go.zip -d  lib/ext/interprophet/unix/
	unzip -o lib/ext/interprophet/win/bindata.go.zip -d  lib/ext/interprophet/win/

	unzip -o lib/ext/peptideprophet/unix/bindata.go.zip -d  lib/ext/peptideprophet/unix/
	unzip -o lib/ext/peptideprophet/win/bindata.go.zip -d  lib/ext/peptideprophet/win/

	unzip -o lib/ext/ptmprophet/unix/bindata.go.zip -d  lib/ext/ptmprophet/unix/
	unzip -o lib/ext/ptmprophet/win/bindata.go.zip -d  lib/ext/ptmprophet/win/

	unzip -o lib/ext/proteinprophet/unix/bindata.go.zip -d  lib/ext/proteinprophet/unix/
	unzip -o lib/ext/proteinprophet/win/bindata.go.zip -d  lib/ext/proteinprophet/win/

	unzip -o lib/pip/bindata.go.zip -d  lib/pip/

	unzip -o lib/dat/bindata.go.zip -d  lib/dat/

	unzip -o lib/obo/unimod/bindata.go.zip -d  lib/obo/unimod/

.PHONY: test
test:
	go test ./... -v

.PHONY: linux
linux:
	gox -os="linux" ${LDFLAGS} -arch=amd64 -output philosopher

.PHONY: windows
windows:
	gox -os="windows" ${LDFLAGS} -arch=amd64 -output philosopher

.PHONY: all
all:
	gox -os="linux" ${LDFLAGS} -arch=amd64 -output philosopher
	gox -os="windows" ${LDFLAGS} -arch=amd64 -output philosopher

.PHONY: draft
draft:
	goreleaser --skip-publish --snapshot --release-notes=Changelog --rm-dist

.PHONY: push
push:
	git tag -a ${TAG} -m "Philosopher ${TAG}"
	git push origin master -f --tags

.PHONY: release
release:
	goreleaser --release-notes=Changelog --rm-dist