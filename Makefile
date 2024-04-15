
PACKAGES := $(shell go list ./...)
name := $(shell basename ${PWD})

all: help

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## vet: vet code
.PHONY: vet
vet:
	go vet $(PACKAGES)

## test: run unit tests
.PHONY: test
test:
	go test -race -cover $(PACKAGES)

## run: run project
.PHONY: run
run:
	go run .

## start: run project
.PHONY: start
start:
	go run .

## build: build project
.PHONY: build
build: test
	go build -o bin/${name} .

## install: install project
.PHONY: install
install: test
	go install

## docker: build project into a docker image and run it in a container
.PHONY: docker
docker: docker-build docker-run

## docker-build: build project into a docker image
.PHONY: docker-build
docker-build:
	GOPROXY=direct docker buildx build -t ${name} .

## docker-run: run project in a docker container
.PHONY: docker-run
docker-run:
	docker run -it --rm ${name}
