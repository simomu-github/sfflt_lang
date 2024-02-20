GOVERSION=$(shell go version)
THIS_GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
THIS_GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
GOOS=$(THIS_GOOS)
GOARCH=$(THIS_GOARCH)
BUILD_TARGETS= \
	build-linux-arm64 \
	build-linux-arm \
	build-linux-amd64 \
	build-darwin-amd64 \
	build-darwin-arm64 \
	build-windows-amd64

RELEASE_TARGETS= \
	release-linux-arm64 \
	release-linux-arm \
	release-linux-amd64 \
	release-darwin-amd64 \
	release-darwin-arm64 \
	release-windows-amd64

TEST_TARGET=...

.PHONY: build deps test clean $(BUILD_TARGETS) $(RELEASE_TARGETS)

build: deps
	go build -o releases/sfflt_lang_$(GOOS)_$(GOARCH)/sfflt_lang$(SUFFIX) cmd/sfflt.go

build-windows-amd64:
	$(MAKE) build GOOS=windows GOARCH=amd64 SUFFIX=.exe

build-windows-386:
	$(MAKE) build GOOS=windows GOARCH=386 SUFFIX=.exe

build-linux-amd64:
	$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-arm:
	$(MAKE) build GOOS=linux GOARCH=arm

build-linux-arm64:
	$(MAKE) build GOOS=linux GOARCH=arm64

build-darwin-amd64:
	$(MAKE) build GOOS=darwin GOARCH=amd64

build-darwin-arm64:
	$(MAKE) build GOOS=darwin GOARCH=arm64

all: $(BUILD_TARGETS)

targz:
	tar -czf releases/sfflt_lang_$(GOOS)_$(GOARCH).tar.gz -C releases sfflt_lang_$(GOOS)_$(GOARCH)

zip:
	cd releases&& zip -9 sfflt_lang_$(GOOS)_$(GOARCH).zip sfflt_lang_$(GOOS)_$(GOARCH)/*

release-windows-amd64:
	$(MAKE) zip GOOS=windows GOARCH=amd64 SUFFIX=.exe

release-windows-386:
	$(MAKE) zip GOOS=windows GOARCH=386 SUFFIX=.exe

release-linux-amd64:
	$(MAKE) targz GOOS=linux GOARCH=amd64

release-linux-arm:
	$(MAKE) targz GOOS=linux GOARCH=arm

release-linux-arm64:
	$(MAKE) targz GOOS=linux GOARCH=arm64

release-darwin-amd64:
	$(MAKE) zip GOOS=darwin GOARCH=amd64

release-darwin-arm64:
	$(MAKE) zip GOOS=darwin GOARCH=arm64

release-all: $(RELEASE_TARGETS)

run_test_script:
	go run cmd/sfflt.go -format pretty test.sflt
	$(FFLT_LANG) test.fflt

deps:
	go mod download

test: deps
	go test -v ./$(TEST_TARGET)

clean:
	rm -r releases/*
