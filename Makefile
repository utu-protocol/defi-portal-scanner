GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)
GIT_DESCR = $(shell git describe --tags --always)
APP=defi-portal-scanner
# build output folder
OUTPUTFOLDER = dist
RELEASEFOLDER = release
# docker image
DOCKER_PROVIDER = 002208042662.dkr.ecr.eu-central-1.amazonaws.com
DOCKER_REGISTRY = $(DOCKER_PROVIDER)
DOCKER_IMAGE = defi-portal-scanner
DOCKER_TAG = $(GIT_DESCR)
# build paramters
OS = linux
ARCH = amd64
# K8S
K8S_NAMESPACE = defi-portal
K8S_DEPLOYMENT = defi-portal-scanner

.PHONY: list
list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs


default: build

build: build-dist

build-dist: $(GOFILES)
	@echo build binary to $(OUTPUTFOLDER)
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go build -ldflags '-s -w -extldflags "-static" -X main.Version=$(GIT_DESCR)' -o $(OUTPUTFOLDER)/$(APP) .
	@echo copy resources
	cp -r README.md LICENSE $(OUTPUTFOLDER)
	@echo done

install:
	@echo installing $(APP)
	GOOS=$(OS) GOARCH=$(ARCH) CGO_ENABLED=0 go install -ldflags '-s -w -extldflags "-static" -X main.Version=$(GIT_DESCR)' .
	@echo done


test: test-all

test-all:
	@go test -v $(GOPACKAGES) -race -coverprofile=cover.out -covermode=atomic

bench: bench-all

bench-all:
	@go test -bench -v $(GOPACKAGES)

lint: lint-all

lint-all:
	@echo run static checks and linting
	staticcheck $(GOPACKAGES)
	golint -set_exit_status $(GOPACKAGES)
	@echo done

clean:
	@echo remove $(OUTPUTFOLDER), $(RELEASEFOLDER) folder
	rm -rf $(OUTPUTFOLDER) $(RELEASEFOLDER)
	@echo done

docker: docker-build

docker-build:
	@echo copy resources
	docker build --platform linux/amd64 --build-arg DOCKER_TAG='$(GIT_DESCR)' -t $(DOCKER_IMAGE)  .
	@echo done

docker-login:
	# docker login $(DOCKER_PROVIDER)
	aws --profile utu.live ecr get-login-password --region eu-central-1 | docker login --username AWS --password-stdin $(DOCKER_REGISTRY)

docker-push: docker-login
	@echo push image
	docker tag $(DOCKER_IMAGE):latest $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	docker tag $(DOCKER_IMAGE):latest $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest
	docker push $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):latest
	@echo done

docker-run: 
	docker run -p 2011:2011 $(DOCKER_IMAGE):latest

debug-start:
	@go run main.go start

k8s-deploy:
	@echo deploy k8s
	kubectl -n $(K8S_NAMESPACE) set image deployment/$(K8S_DEPLOYMENT) $(K8S_DEPLOYMENT)=$(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo done

k8s-rollback:
	@echo deploy k8s
	kubectl -n $(K8S_NAMESPACE) rollout undo deployment/$(K8S_DEPLOYMENT)
	@echo done

changelog:
	git-chglog --output CHANGELOG.md

git-release:
	@echo making release
	git tag $(GIT_DESCR)
	git-chglog --output CHANGELOG.md
	git tag $(GIT_DESCR) --delete
	git add CHANGELOG.md && git commit -m "$(GIT_DESCR)" -m "Changelog: https://github.com/utu-crowdsale/$(APP)/blob/master/CHANGELOG.md"
	git tag -s -a "$(GIT_DESCR)" -m "Changelog: https://github.com/utu-crowdsale/$(APP)/blob/master/CHANGELOG.md"
	@echo release complete


_release-patch:
	$(eval GIT_DESCR = $(shell git describe --tags | awk -F '("|")' '{ print($$1)}' | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g'))
release-patch: _release-patch git-release

_release-minor:
	$(eval GIT_DESCR = $(shell git describe --tags | awk -F '("|")' '{ print($$1)}' | awk -F. '{$$(NF-1) = $$(NF-1) + 1;} 1' | sed 's/ /./g' | awk -F. '{$$(NF) = 0;} 1' | sed 's/ /./g'))
release-minor: _release-minor git-release

_release-major:
	$(eval GIT_DESCR = $(shell git describe --tags | awk -F '("|")' '{ print($$1)}' | awk -F. '{$$(NF-2) = $$(NF-2) + 1;} 1' | sed 's/ /./g' | awk -F. '{$$(NF-1) = 0;} 1' | sed 's/ /./g' | awk -F. '{$$(NF) = 0;} 1' | sed 's/ /./g' ))
release-major: _release-major git-release 

gh-publish-release: clean build
	@echo publish release
	mkdir -p $(RELEASEFOLDER)
	zip -rmT $(RELEASEFOLDER)/$(APP)-$(GIT_DESCR).zip $(OUTPUTFOLDER)/
	sha256sum $(RELEASEFOLDER)/$(APP)-$(GIT_DESCR).zip | tee $(RELEASEFOLDER)/$(APP)-$(GIT_DESCR).zip.checksum
	gh release create $(GIT_DESCR) $(RELEASEFOLDER)/* -t v$(GIT_DESCR) -F CHANGELOG.md
	@echo done


################# CUSTOM
# adduser utu --shell=/bin/false --system --no-create-home --group 

SSH_HOST = utu.$(APP)
SSH_BIN_FOLDER = /data/$(APP)/bin

deploy-ssh: clean build
	rsync --progress --checksum --archive --human-readable --verbose $(OUTPUTFOLDER)/ $(SSH_HOST):$(SSH_BIN_FOLDER)/

deploy-ssh-restart:
	@echo stopping $(APP) on $(SSH_HOST)
	ssh -t $(SSH_HOST) "systemctl start $(APP)"
	@echo deploy complete

