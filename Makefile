SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY = philosopher

VERSION = $(shell date +%Y%m%d)
BUILD = $(shell  date +%Y%m%d%H%M)

TAG = v4.5.1

LDFLAGS = -ldflags "-w -s -extldflags -static -X main.version=${TAG} -X main.build=${BUILD}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} main.go

.PHONY: deploy
deploy:
	unzip -o lib/dat/bindata.go.zip -d lib/dat/

	unzip -o lib/ext/cdhit/unix/bindata.go.zip -d lib/ext/cdhit/unix/
	unzip -o lib/ext/cdhit/win/bindata.go.zip -d lib/ext/cdhit/win/

	unzip -o lib/ext/comet/unix/bindata.go.zip -d lib/ext/comet/unix/
	unzip -o lib/ext/comet/win/bindata.go.zip -d lib/ext/comet/win/rm -rf /usr/local/go && tar -C /usr/local -xzf go1.18.2.linux-amd64.tar.gz
	unzip -o lib/ext/rawfilereader/reh64/bindata.go.zip -d lib/ext/rawfilereader/reh64/
	unzip -o lib/ext/rawfilereader/win/bindata.go.zip -d lib/ext/rawfilereader/win/

	unzip -o lib/pip/bindata.go.zip -d  lib/pip/

	unzip -o lib/dat/bindata.go.zip -d  lib/dat/

	unzip -o lib/obo/unimod/bindata.go.zip -d  lib/obo/unimod/

.PHONY: test
test:
	go test ./... -v

.PHONY: rc
rc:
	env CGO_ENABLED=0 gox -os="linux" ${LDFLAGS} -arch=amd64 -output philosopher-${TAG}-RC
	env CGO_ENABLED=0 gox -os="windows" ${LDFLAGS} -arch=amd64 -output philosopher-${TAG}-RC
	mv philosopher-${TAG}-RC ~/bin/
	mv philosopher-${TAG}-RC.exe ~/bin/

.PHONY: linux
linux:
	env CGO_ENABLED=0 gox -os="linux" ${LDFLAGS} -arch=amd64 -output philosopher

.PHONY: windows
windows:
	env CGO_ENABLED=0 gox -os="windows" ${LDFLAGS} -arch=amd64 -output philosopher

.PHONY: all
all:
	env CGO_ENABLED=0 gox -os="linux" ${LDFLAGS} -arch=amd64 -output philosopher
	env CGO_ENABLED=0 gox -os="windows" ${LDFLAGS} -arch=amd64 -output philosopher

.PHONY: push
push:
	git tag -a ${TAG} -m "Philosopher ${TAG}"
	git push origin master -f --tags
	
.PHONY: draft
draft:
	goreleaser --skip-publish --snapshot --release-notes=Changelog --rm-dist

.PHONY: release
release:
	goreleaser --release-notes=Changelog --rm-dist
