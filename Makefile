#include .env


PROJECTNAME=$(shell basename "$(PWD)")

# Go related variables.
GOBASE=$(shell pwd)
#GOPATH="$(GOBASE)/vendor:$(GOBASE)"
GOBIN=$(GOBASE)/bin
CMD_DIR=cmd
HOME_DIR=${HOME}
GOFILES=$(wildcard *.go)
PLATFORM=linux/amd64

build:
	@echo "  >  Cleaning build cache"
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) go build pkg/main.go

run:
	@echo "  >  Running main.go"
	@ENV=$(ENV) GOPATH=$(GOPATH) go run pkg/main.go
	@echo "  >  Done!"

test:
	@echo "  >  Starting tests"
	@go test -v ./pkg/usecases/
	@echo "  >  Done!"

docker-build:
	@echo "  >  Buidling docker image"
	@docker build -t ticker . 
	@echo "  >  Done!"


docker-run:
	@echo "  >  Starting docker image"
	@docker run  -tid ticker 
	@echo "  >  Done!"

docker-stop:
	@echo "  >  Stopping image with id $(id)"
	@docker stop $(id)
	@echo "  >  Done!"

docker-list:
	@echo "  >  Listing docker images"
	@docker ps
	@echo "  >  Done!"

docker-enter:
	@echo "  >  Entering image with id $(id)"
	@docker exec -it $(id) sh
	@echo "  >  Done!"


