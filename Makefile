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

modules: kube-module openshift-module

kube-module:
	@echo "Building Kubernetes module..."
	@go build -buildmode=plugin \
		-o $(BUILD_DIR)/plugins/kube.so \
		-ldflags='-s -w' \
		plugins/kube/*.go

openshift-module:
	@echo "Building OpenShift module..."
	@go build -buildmode=plugin \
		-o $(BUILD_DIR)/plugins/openshift.so \
		-ldflags='-s -w' \
		plugins/openshift/*.go

docker-image:
	@echo "Building docker image..."
	@docker build -t crunchy-watch \
			-f $(CCP_BASEOS)/$(CCP_PGVERSION)/Dockerfile.watch.$(CCP_BASEOS) .
	@docker tag crunchy-watch \
			crunchydata/crunchy-watch:$(CCP_BASEOS)-$(CCP_PG_FULLVERSION)-$(CCP_VERSION)

setup:
	@echo "Downloading tools..."
	mkdir -p $(TOOLS_DIR)
	@echo "Downloading kubectl..."
	@curl -o $(TOOLS_DIR)/kubectl https://storage.googleapis.com/kubernetes-release/release/$(shell curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
	@chmod +x $(TOOLS_DIR)/kubectl
	@echo "Downloading oc..."
	@curl -L -o $(TOOLS_DIR)/openshift.tar.gz https://github.com/openshift/origin/releases/download/v3.6.1/openshift-origin-server-v3.6.1-008f2d5-linux-64bit.tar.gz
	@mkdir -p $(TOOLS_DIR)/openshift
	@tar --warning=no-unknown-keyword -zxvf $(TOOLS_DIR)/openshift.tar.gz -C $(TOOLS_DIR)/openshift --strip 1
