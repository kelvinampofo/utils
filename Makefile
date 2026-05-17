SHELL = /bin/sh

BIN_DIR := bin
CMD_DIR := cmd
SCRIPT_DIR := scripts
CMDS := $(notdir $(wildcard $(CMD_DIR)/*))
SCRIPTS := $(notdir $(wildcard $(SCRIPT_DIR)/*))
TOOLS := $(sort $(CMDS) $(SCRIPTS))

.SUFFIXES:

.PHONY: help list fmt test check build clean

help:
	@echo "Targets:"
	@echo "  make list                 # list available CLI commands"
	@echo "  make fmt                  # run gofmt on all Go source files"
	@echo "  make test                 # run Go tests"
	@echo "  make check                # run formatting and tests"
	@echo "  make build                # build/copy tools into ./bin"
	@echo "  make clean                # remove ./bin"

list:
	@printf "%s\n" $(TOOLS)

fmt:
	gofmt -w $$(find $(CMD_DIR) -name '*.go')

test:
	go test ./...

check: fmt test

build:
	@mkdir -p $(BIN_DIR)
	@set -e; for c in $(CMDS); do \
		go build -o $(BIN_DIR)/$$c ./$(CMD_DIR)/$$c; \
		echo "Built $(BIN_DIR)/$$c"; \
	done
	@set -e; for s in $(SCRIPTS); do \
		install -m 0755 $(SCRIPT_DIR)/$$s $(BIN_DIR)/$$s; \
		echo "Installed $(BIN_DIR)/$$s"; \
	done

clean:
	rm -rf $(BIN_DIR)
