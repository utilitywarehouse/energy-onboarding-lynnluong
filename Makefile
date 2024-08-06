mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
base_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

SERVICE ?= $(base_dir)

DOCKER_REGISTRY?=registry.uw.systems
DOCKER_REPOSITORY_NAMESPACE?=energy
DOCKER_REPOSITORY_IMAGE=$(SERVICE)
DOCKER_REPOSITORY=$(DOCKER_REGISTRY)/$(DOCKER_REPOSITORY_NAMESPACE)/$(DOCKER_REPOSITORY_IMAGE)

K8S_NAMESPACE=$(DOCKER_REPOSITORY_NAMESPACE)
K8S_DEPLOYMENT_NAME=$(DOCKER_REPOSITORY_IMAGE)
K8S_CONTAINER_NAME=$(K8S_DEPLOYMENT_NAME)

BUILDENV :=
BUILDENV += CGO_ENABLED=0
GIT_HASH := $(GITHUB_SHA)
ifeq ($(GIT_HASH),)
  GIT_HASH := $(shell git rev-parse HEAD)
endif
LINKFLAGS :=-s -X main.gitHash=$(GIT_HASH) -extldflags "-static"
TESTFLAGS := -v -cover
LINTER := golangci-lint

EMPTY :=
SPACE := $(EMPTY) $(EMPTY)
join-with = $(subst $(SPACE),$1,$(strip $2))

LEXC :=

.PHONY: install
install:
	GOPROXY=https://proxy.golang.org GO111MODULE=on GOPRIVATE="github.com/utilitywarehouse/*" go mod download

$(LINTER):
	@ [ -e ./bin/$(LINTER) ] || wget -O - -q https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s latest

.PHONY: lint
lint: $(LINTER)
	./bin/$(LINTER) run $(LINT_FLAGS)

.PHONY: clean
clean:
	rm -f $(SERVICE)

# builds our binary
$(SERVICE): clean
	GO111MODULE=on $(BUILDENV) go build -o $(SERVICE) -a -ldflags '$(LINKFLAGS)' ./cmd/$(SERVICE)

build: $(SERVICE)

.PHONY: test
test:
	GO111MODULE=on $(BUILDENV) go test $(TESTFLAGS) ./...

.PHONY: all
all: clean $(LINTER) lint test build

docker-image:
	docker build -t $(DOCKER_REPOSITORY):local . --build-arg SERVICE=$(SERVICE) --build-arg GITHUB_TOKEN=$(GITHUB_TOKEN)

ci-docker-auth:
	@echo "Logging in to $(DOCKER_REGISTRY) as $(DOCKER_ID)"
	@docker login -u $(DOCKER_ID) -p $(DOCKER_PASSWORD) $(DOCKER_REGISTRY)

ci-docker-build: ci-docker-auth
	docker build -t $(DOCKER_REPOSITORY):$(GITHUB_SHA) . --build-arg SERVICE=$(SERVICE) --build-arg GITHUB_TOKEN=$(GITHUB_TOKEN)
	if [ -n "$(CIRCLE_BRANCH)" ]; then docker tag $(DOCKER_REPOSITORY):$(GITHUB_SHA) $(DOCKER_REPOSITORY):$(shell echo $(CIRCLE_BRANCH) | awk -F '/' '{gsub("/","_"); print $0}'); fi
	if [ -n "$(CIRCLE_TAG)" ]; then docker tag $(DOCKER_REPOSITORY):$(GITHUB_SHA) $(DOCKER_REPOSITORY):$(CIRCLE_TAG); fi
	docker tag $(DOCKER_REPOSITORY):$(GITHUB_SHA) $(DOCKER_REPOSITORY):latest
	docker push -a $(DOCKER_REPOSITORY)

K8S_URL=https://elb.master.k8s.dev.uw.systems/apis/extensions/v1beta1/namespaces/$(K8S_NAMESPACE)/deployments/$(K8S_DEPLOYMENT_NAME)
K8S_PAYLOAD={"spec":{"template":{"spec":{"containers":[{"name":"$(K8S_CONTAINER_NAME)","image":"$(DOCKER_REPOSITORY):$(GITHUB_SHA)"}]}}}}

ci-kubernetes-push:
	test "$(shell curl -o /dev/null -w '%{http_code}' -s -X PATCH -k -d '$(K8S_PAYLOAD)' -H 'Content-Type: application/strategic-merge-patch+json' -H 'Authorization: Bearer $(K8S_DEV_TOKEN)' '$(K8S_URL)')" -eq "200"

.PHONY: regen
regen:
	microservice-generator new energy-onboarding-lynnluong -d "Energy onboarding task" --kafka-in ServiceTopic --postgres --output-path . --kubernetes-manifest-dir deployment --kubernetes-namespace energy-platform
