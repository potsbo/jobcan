NAME := jobcan
SRCS := $(shell find . -type f -name '*.go')

bin/$(NAME): $(SRCS) vendor/*
	go build -o bin/$(NAME)

.PHONY: dep
dep:
ifeq ($(shell command -v dep 2> /dev/null),)
	go get -u github.com/golang/dep/cmd/dep
endif

.PHONY: deps
deps: dep
	dep ensure -v

vendor/*: Gopkg.toml Gopkg.lock
	@make deps
