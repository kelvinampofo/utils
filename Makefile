BIN_DIR := bin
CMD_DIR := cmd
CMDS := $(notdir $(wildcard $(CMD_DIR)/*))

.PHONY: help list fmt test check build clean

help:
	@echo "Targets:"
	@echo "  make list                 # list available CLI commands"
	@echo "  make fmt                  # run gofmt on all Go source files"
	@echo "  make test                 # run Go tests"
	@echo "  make check                # run formatting and tests"
	@echo "  make build                # build Go CLIs into ./bin"
	@echo "  make clean                # remove ./bin"

list:
	@printf "%s\n" $(CMDS)

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

clean:
	rm -rf $(BIN_DIR)
