SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY = philosopher

VERSION = $(shell date +%Y%m%d)
BUILD = $(shell  date +%Y%m%d%H%M)

TAG = v1.4.4

LDFLAGS = -ldflags "-w -s -X main.version=${TAG} -X main.build=${BUILD}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} main.go

.PHONY: deps
deps:
	go get -u github.com/mitchellh/gox
	go get -u github.com/inconshreveable/mousetrap
	go get -u github.com/sirupsen/logrus
	go get -u gonum.org/v1/plot
	go get -u github.com/mattn/go-colorable
	go get -u github.com/montanaflynn/stats
	go get -u github.com/pierrre/archivefile/zip
	go get -u github.com/rogpeppe/go-charset/charset
	go get -u github.com/rogpeppe/go-charset/data
	go get -u github.com/satori/go.uuid
	go get -u github.com/spf13/cobra
	go get -u github.com/spf13/viper
	go get -u golang.org/x/net/html/charset
	go get -u github.com/spf13/cobra/cobra
	go get -u github.com/nlopes/slack
	go get -u github.com/google/go-github/github
	go get -u github.com/vmihailenco/msgpack
	go get -u github.com/davecgh/go-spew/spew
	go get -u github.com/jpillora/go-ogle-analytics
	go get -u github.com/onsi/ginkgo
	go get -u github.com/onsi/gomega
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/konsorten/go-windows-terminal-sequences
	go get github.com/blang/semver

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

	unzip -o lib/ext/msconvert/unix/bindata.go.zip -d  lib/ext/msconvert/unix/
	unzip -o lib/ext/msconvert/darwin/bindata.go.zip -d  lib/ext/msconvert/darwin/

	unzip -o lib/ext/idconvert/unix/bindata.go.zip -d  lib/ext/idconvert/unix/
	unzip -o lib/ext/idconvert/darwin/bindata.go.zip -d  lib/ext/idconvert/darwin/

	unzip -o lib/pip/bindata.go.zip -d  lib/pip/

	unzip -o lib/dat/bindata.go.zip -d  lib/dat/

	unzip -o lib/obo/unimod/bindata.go.zip -d  lib/obo/unimod/

.PHONY: coverage
coverage:
	ginkgo -r -cover -outputdir test/profiles

.PHONY: test
test:
	ginkgo -r

.PHONY: linux
linux:
	gox -os="linux" ${LDFLAGS} -arch=amd64 -output philosopher

.PHONY: windows
windows:
	gox -os="windows" ${LDFLAGS} -arch=amd64 -output philosopher.exe

.PHONY: draft
draft:
	ginkgo -r
	goreleaser --skip-publish --snapshot --release-notes=Changelog --rm-dist

.PHONY: push
push:
	ginkgo -r
	git tag -a ${TAG} -m "Philosopher ${TAG}"
	git push origin master -f --tags

.PHONY: release
release:
	goreleaser --release-notes=Changelog --rm-dist