
PKG=github.com/chandresh-pancholi/csi-gce
IMAGE=reg.onec.co/csi-gce
VERSION=0.1.0-alpha
GIT_COMMIT?=$(shell git rev-parse HEAD)
BUILD_DATE?=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS?="-X ${PKG}/pkg/driver.driverVersion=${VERSION} -X ${PKG}/pkg/driver.gitCommit=${GIT_COMMIT} -X ${PKG}/pkg/driver.buildDate=${BUILD_DATE}"

.PHONY: csi-gce
csi-gce:
	mkdir -p bin
	CGO_ENABLED=0 GOOS=linux go build -ldflags ${LDFLAGS} -o bin/csi-gce ./cmd/

.PHONY: test
test:
	go test -v -race ./pkg/...

.PHONY: test-sanity
test-sanity:
	go test -v ./tests/sanity/...

.PHONY: test-integration
test-integration:
	go test -c ./tests/integration/... -o bin/integration.test && \
	sudo -E bin/integration.test -ginkgo.v

.PHONY: image
image:
	docker build -t $(IMAGE):$(VERSION) .

.PHONY: push
push:
	docker push $(IMAGE):$(VERSION)
