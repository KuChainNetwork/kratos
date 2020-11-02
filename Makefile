#!/usr/bin/make -f

PACKAGES_SIMTEST=$(shell go list ./... | grep '/simulation')
VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
COMMIT := $(shell git log -1 --format='%H')
LEDGER_ENABLED ?= true
SDK_PACK := $(shell go list -m github.com/cosmos/cosmos-sdk | sed  's/ /\@/g')
TIME_BEGIN := $(shell date -u +"%Y-%m-%d %H:%M.%S")
BRANCH := $(shell echo $(shell git rev-parse --abbrev-ref HEAD) | sed 's/^v//')

$(info "Kuchain Version: ${VERSION} ${SDK_PACK} in ${TIME_BEGIN}")
$(info "Current branch: ${BRANCH}")

MAIN_SYMBOL := kuchain
CORE_SYMBOL := sys

export GO111MODULE = on

# process build tags

build_tags = netgo
ifeq ($(LEDGER_ENABLED),true)
  ifeq ($(OS),Windows_NT)
    GCCEXE = $(shell where gcc.exe 2> NUL)
    ifeq ($(GCCEXE),)
      $(error gcc.exe not installed for ledger support, please install or set LEDGER_ENABLED=false)
    else
      build_tags += ledger
    endif
  else
    UNAME_S = $(shell uname -s)
    ifeq ($(UNAME_S),OpenBSD)
      $(warning OpenBSD detected, disabling ledger support (https://github.com/cosmos/cosmos-sdk/issues/1988))
    else
      GCC = $(shell command -v gcc 2> /dev/null)
      ifeq ($(GCC),)
        $(error gcc not installed for ledger support, please install or set LEDGER_ENABLED=false)
      else
        build_tags += ledger
      endif
    endif
  endif
endif

ifeq ($(WITH_CLEVELDB),yes)
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace :=
whitespace += $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=$(MAIN_SYMBOL) \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=kucd \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=kucli \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)" \
		  -X github.com/KuChainNetwork/kuchain/chain/constants.KuchainBuildVersion=$(VERSION) \
		  -X github.com/KuChainNetwork/kuchain/chain/constants.KuchainBuildBranch=$(BRANCH) \
		  -X "github.com/KuChainNetwork/kuchain/chain/constants.KuchainBuildTime=$(TIME_BEGIN)" \
		  -X github.com/KuChainNetwork/kuchain/chain/constants.KuchainBuildSDKVersion=$(SDK_PACK) \
		  -X github.com/KuChainNetwork/kuchain/chain/constants/keys.ChainNameStr=$(CORE_SYMBOL) \
		  -X github.com/KuChainNetwork/kuchain/chain/constants/keys.ChainMainNameStr=$(MAIN_SYMBOL)

ifeq ($(WITH_CLEVELDB),yes)
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

OS_NAME := $(shell uname -s | tr A-Z a-z)/$(shell uname -m)

os:
	@echo $(OS_NAME)

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

BUILD_OS := -osarch="linux/amd64 windows/amd64"
ifneq ($(OS_NAME), darwin/x86_64)
	BUILD_OS = -osarch="linux/amd64 windows/amd64 darwin/amd64"
endif

all: clear-build build

build: go.sum
ifeq ($(OS),Windows_NT)
	go build -mod=readonly $(BUILD_FLAGS) -o build/kucd.exe ./cmd/kucd
	go build -mod=readonly $(BUILD_FLAGS) -o build/kucli.exe ./cmd/kucli
else
	go build -mod=readonly $(BUILD_FLAGS) -o build/kucd ./cmd/kucd
	go build -mod=readonly $(BUILD_FLAGS) -o build/kucli ./cmd/kucli
endif

build-gox: go.sum
	gox $(BUILD_FLAGS) $(BUILD_OS) -output "build/kucd_{{.OS}}_{{.Arch}}" ./cmd/kucd
	gox $(BUILD_FLAGS) $(BUILD_OS) -output "build/kucli_{{.OS}}_{{.Arch}}" ./cmd/kucli

ifeq ($(OS_NAME), darwin/x86_64)
	go build -mod=readonly $(BUILD_FLAGS) -o build/kucd_darwin_amd64 ./cmd/kucd
	go build -mod=readonly $(BUILD_FLAGS) -o build/kucli_darwin_amd64 ./cmd/kucli
endif

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/kucd
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/kucli

go-mod-cache: go.sum
	@echo "--> Download go modules to local cache"
	@go mod download

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

draw-deps:
	@# requires brew install graphviz or apt-get install graphviz
	go get github.com/RobotsAndPencils/goviz
	@goviz -i ./cmd/kucd -d 2 | dot -Tpng -o dependency-graph.png

clear-build:
	rm -rf ./build

