GOCMD=go
GOINSTALL=$(GOCMD) install
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=ginkgo

VERSION=`git describe --tags`
VER=`git describe --tags --abbrev=0`
DATE=`date -u '+%Y-%m-%d_%I:%M:%S%p'`
LDFLAGS=-ldflags "-X main.buildDate=${DATE} \
                  -X main.buildVersion=${VERSION}"


all: install

test:
	ginkgo

install:
	cd backme; \
	$(GOCLEAN); \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOINSTALL) ${LDFLAGS};


build:
	cd backme; \
	$(GOCLEAN); \
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) ${LDFLAGS};

release:
	cd backme; \
	$(GOCLEAN); \
	GOOS=linux GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
	tar zcvf /tmp/backme-${VER}-linux.tar.gz backme; \
	$(GOCLEAN); \
	GOOS=darwin GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
	tar zcvf /tmp/backme-${VER}-mac.tar.gz backme; \
	$(GOCLEAN); \
	GOOS=windows GOARCH=amd64 $(GOBUILD) ${LDFLAGS}; \
	zip -9 /tmp/backme-${VER}-win-64.zip backme.exe; \
	$(GOCLEAN);
