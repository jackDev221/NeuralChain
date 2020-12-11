# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: gnc android ios gnc-cross evm all test clean
.PHONY: gnc-linux gnc-linux-386 gnc-linux-amd64 gnc-linux-mips64 gnc-linux-mips64le
.PHONY: gnc-linux-arm gnc-linux-arm-5 gnc-linux-arm-6 gnc-linux-arm-7 gnc-linux-arm64
.PHONY: gnc-darwin gnc-darwin-386 gnc-darwin-amd64
.PHONY: gnc-windows gnc-windows-386 gnc-windows-amd64

GOBIN = $(shell pwd)/build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

gnc:
	$(GORUN) build/ci.go install ./cmd/gnc
	@echo "Done building."
	@echo "Run \"$(GOBIN)/gnc\" to launch gnc."

all:
	$(GORUN) build/ci.go install

android:
	$(GORUN) build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/geth.aar\" to use the library."

ios:
	$(GORUN) build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/Geth.framework\" to use the library."

test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	./build/clean_go_build_cache.sh
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

gnc-cross: gnc-linux gnc-darwin gnc-windows gnc-android gnc-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/gnc-*

gnc-linux: gnc-linux-386 gnc-linux-amd64 gnc-linux-arm gnc-linux-mips64 gnc-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-*

gnc-linux-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/gnc
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-* | grep 386

gnc-linux-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/gnc
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-* | grep amd64

gnc-linux-arm: gnc-linux-arm-5 gnc-linux-arm-6 gnc-linux-arm-7 gnc-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-* | grep arm

gnc-linux-arm-5:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/gnc
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-* | grep arm-5

gnc-linux-arm-6:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/gnc
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-* | grep arm-6

gnc-linux-arm-7:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/gnc
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-* | grep arm-7

gnc-linux-arm64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/gnc
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-* | grep arm64

gnc-linux-mips:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/gnc
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-* | grep mips

gnc-linux-mipsle:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/gnc
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-* | grep mipsle

gnc-linux-mips64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/gnc
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-* | grep mips64

gnc-linux-mips64le:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/gnc
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/gnc-linux-* | grep mips64le

gnc-darwin: gnc-darwin-386 gnc-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/gnc-darwin-*

gnc-darwin-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/gnc
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/gnc-darwin-* | grep 386

gnc-darwin-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/gnc
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gnc-darwin-* | grep amd64

gnc-windows: gnc-windows-386 gnc-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/gnc-windows-*

gnc-windows-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/gnc
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/gnc-windows-* | grep 386

gnc-windows-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/gnc
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/gnc-windows-* | grep amd64
