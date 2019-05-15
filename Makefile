.PHONY: all build

VERSION := $(shell cat VERSION)
REPO_URL := $(shell cat REPO_URL)

build:
	mkdir -p build
	env GOOS=linux GOARCH=amd64 go build -o build/rump

build_container:
	docker rmi -f rump:$(VERSION) >/dev/null 2>&1 || true
	docker build -t rump:$(VERSION) -f Dockerfile .

clean_build_container:
	docker rmi -f rump:$(VERSION) >/dev/null 2>&1 || true
	docker build --no-cache -t rump:$(VERSION) -f Dockerfile .

upload:
	$(eval CONTAINER_ID := $(shell docker image ls rump | grep '$(VERSION)' | awk '{print $$3}'))
	docker tag $(CONTAINER_ID) $(REPO_URL):$(VERSION)
	docker push $(REPO_URL):$(VERSION)
