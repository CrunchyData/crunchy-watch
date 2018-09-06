.PHONY: all build modules kube-module openshift-module clean resolve docker-image setup

PROJECT_DIR := $(shell pwd)
BUILD_DIR := $(PROJECT_DIR)/build
RELEASE_DIR := $(PROJECT_DIR)/release
TOOLS_DIR := $(PROJECT_DIR)/tools
VENDOR_DIR := $(PROJECT_DIR)/vendor
DOCS_DIR := $(PROJECT_DIR)/docs

all: clean resolve build modules

clean:
	@echo "Cleaning project..."
	@rm -rf $(BUILD_DIR)
	@rm -rf $(RELEASE_DIR)
	@rm -rf $(VENDOR_DIR)
	@go clean -i

resolve:
	@echo "Resolving dependencies..."
	@dep ensure

build:
	@echo "Building crunchy-watch..."
	@go build -i -o $(BUILD_DIR)/crunchy-watch \
		-ldflags='-s -w' \
		$(PROJECT_DIR)/*.go

modules: kube-module 

kube-module:
	@echo "Building Kubernetes module..."
	@go build -buildmode=plugin \
		-o $(BUILD_DIR)/plugins/kube.so \
		-ldflags='-s -w' \
		plugins/kube/*.go

docker-image:
	@echo "Building docker image..."
	@docker build -t crunchy-watch \
			-f $(CCP_BASEOS)/$(CCP_PGVERSION)/Dockerfile.watch.$(CCP_BASEOS) .
	@docker tag crunchy-watch \
			crunchydata/crunchy-watch:$(CCP_BASEOS)-$(CCP_PG_FULLVERSION)-$(CCP_VERSION)

setup:
	@echo "Downloading tools..."
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
	go get github.com/blang/expenv
