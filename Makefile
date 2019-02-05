SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY = philosopher

VERSION = $(shell date +%Y%m%d)
BUILD = $(shell  date +%Y%m%d%H%M)

LDFLAGS = -ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD}"

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
	#go get -u github.com/konsorten/go-windows-terminal-sequences

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

	unzip -o lib/uni/bindata.go.zip -d  lib/uni/

	unzip -o lib/pip/bindata.go.zip -d  lib/pip/

	unzip -o lib/dat/bindata.go.zip -d  lib/dat/

.PHONY: test
test:
	ginkgo -r -cover -outputdir test/profiles

.PHONY: install
install:
	gox -os="linux" ${LDFLAGS} -arch=amd64 -output philosopher.${VERSION}
	cp philosopher.${VERSION} ${GOBIN}/philosopher;
	mv philosopher.${VERSION} ${GOBIN}/philosopher.${VERSION};

.PHONY: linux
linux:
	gox -os="linux" ${LDFLAGS} -arch=amd64 -output philosopher.${VERSION}

.PHONY: windows
windows:
	gox -os="windows" ${LDFLAGS} -arch=amd64 -output philosopher.${VERSION}
	#cp philosopher.${VERSION}.exe /home/prvst/Public/philosopher.exe

.PHONY: release
release:
	gox ${LDFLAGS}

.PHONY: all
all:
	gox -os="linux" ${LDFLAGS} -arch=amd64 -output philosopher.${VERSION}
	cp philosopher.${VERSION} ${GOBIN}/philosopher;
	#cp philosopher.${VERSION} /home/felipevl/Servers/castor/home/felipevl/bin/philosopher
	#cp philosopher.${VERSION} /home/felipevl/Servers/pathbio/bin/philosopher
	mv philosopher.${VERSION} ${GOBIN}/philosopher.${VERSION};
	gox ${LDFLAGS} .

.PHONY: clean
clean:
	if [ -f ${BINARY} ]; then rm ${BINARY}; fi
	if [ -f ${BINARY}_linux_386 ]; then rm ${BINARY}_linux_386; fi
	if [ -f ${BINARY}_linux_amd64 ]; then rm ${BINARY}_linux_amd64; fi
	if [ -f ${BINARY}_linux_arm ]; then rm ${BINARY}_linux_arm; fi
	if [ -f ${BINARY}_windows_386.exe ]; then rm ${BINARY}_windows_386.exe; fi
	if [ -f ${BINARY}_windows_amd64.exe ]; then rm ${BINARY}_windows_amd64.exe; fi
	if [ -f ${BINARY}_darwin_amd64 ]; then rm ${BINARY}_darwin_amd64; fi
	if [ -f ${BINARY}_darwin_386 ]; then rm ${BINARY}_darwin_386; fi
	if [ -f ${BINARY}_freebsd_amd64 ]; then rm ${BINARY}_freebsd_amd64; fi
	if [ -f ${BINARY}_freebsd_386 ]; then rm ${BINARY}_freebsd_386; fi
	if [ -f ${BINARY}_freebsd_arm ]; then rm ${BINARY}_freebsd_arm; fi
	if [ -f ${BINARY}_netbsd_amd64 ]; then rm ${BINARY}_netbsd_amd64; fi
	if [ -f ${BINARY}_netbsd_386 ]; then rm ${BINARY}_netbsd_386; fi
	if [ -f ${BINARY}_netbsd_arm ]; then rm ${BINARY}_netbsd_arm; fi
	if [ -f ${BINARY}_openbsd_amd64 ]; then rm ${BINARY}_openbsd_amd64; fi
	if [ -f ${BINARY}_openbsd_386 ]; then rm ${BINARY}_openbsd_386; fi
	if [ -f ${BINARY}_linux_mips ]; then rm ${BINARY}_linux_mips; fi
	if [ -f ${BINARY}_linux_mips64 ]; then rm ${BINARY}_linux_mips64; fi
	if [ -f ${BINARY}_linux_mips64le ]; then rm ${BINARY}_linux_mips64le; fi
	if [ -f ${BINARY}_linux_mipsle ]; then rm ${BINARY}_linux_mipsle; fi
	if [ -f ${BINARY}_linux_s390x ]; then rm ${BINARY}_linux_s390x; fi
