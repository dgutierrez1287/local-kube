MAKEFILE_DIR := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
SHELL := /usr/bin/env bash
.DEFAULT_GOAL := help

# Detect OS
OS := $(shell uname)

INSTALL_DIR ?= "/usr/local/bin"

help:
	@echo "local-kube Makefile"
	@echo "============================"
	@echo ""
	@echo "help - shows all the help information"
	@echo "build - builds the local-kube executable"
	@echo "install - installs the executable"
	@echo "uninstall - uninstalls the executable"
	@echo "tidy - runs go mod tidy to update modules"

build:
	go build .

clean:
	rm -f ./local-kube

install:
	sudo cp -f ./local-kube ${INSTALL_DIR} && sudo chmod 777 ${INSTALL_DIR}/local-kube

uninstall:
	sudo rm -f ${INSTALL_DIR}/local-kube

tidy:
	go mod tidy
