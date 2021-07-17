HELM_3_PLUGINS := $(shell bash -c 'eval $$(helm env); echo $$HELM_PLUGINS')

PKG:= github.com/billyshambrook/helm-init
LDFLAGS := -X $(PKG)/cmd.Version=$(VERSION)

# Clear the "unreleased" string in BuildMetadata
LDFLAGS += -X k8s.io/helm/pkg/version.BuildMetadata=
LDFLAGS += -X k8s.io/helm/pkg/version.Version=$(shell ./scripts/dep-helm-version.sh)

.PHONY: install
install: build
	mkdir -p $(HELM_3_PLUGINS)/helm-init/bin
	cp bin/init $(HELM_3_PLUGINS)/helm-init/bin
	cp plugin.yaml $(HELM_3_PLUGINS)/helm-init/

.PHONY: build
build:
	mkdir -p bin/
	go build -v -o bin/init -ldflags="$(LDFLAGS)" main.go
