HELM_3_PLUGINS := $(shell bash -c 'eval $$(helm env); echo $$HELM_PLUGINS')

.PHONY: install
install: build
	mkdir -p $(HELM_3_PLUGINS)/helm-init/bin
	cp bin/init $(HELM_3_PLUGINS)/helm-init/bin
	cp plugin.yaml $(HELM_3_PLUGINS)/helm-init/

.PHONY: build
build:
	mkdir -p bin/
	go build -v -o bin/init main.go

.PHONY: dist
dist: export COPYFILE_DISABLE=1 #teach OSX tar to not put ._* files in tar archive
dist: export CGO_ENABLED=0
dist:
	rm -rf build/init/* release/*
	mkdir -p build/init/bin release/
	cp README.md LICENSE plugin.yaml build/init
	GOOS=linux GOARCH=amd64 go build -o build/init/bin/init -trimpath -ldflags="$(LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-init-linux.tgz init/
	GOOS=freebsd GOARCH=amd64 go build -o build/init/bin/init -trimpath -ldflags="$(LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-init-freebsd.tgz init/
	GOOS=darwin GOARCH=amd64 go build -o build/init/bin/init -trimpath -ldflags="$(LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-init-macos.tgz init/
	rm build/init/bin/init
	GOOS=windows GOARCH=amd64 go build -o build/init/bin/init.exe -trimpath -ldflags="$(LDFLAGS)"
	tar -C build/ -zcvf $(CURDIR)/release/helm-init-windows.tgz init/

.PHONY: release
release: dist
ifndef GITHUB_TOKEN
	$(error GITHUB_TOKEN is undefined)
endif
	scripts/release.sh v$(VERSION) master
