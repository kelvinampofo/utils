BIN_DIR := bin
CMD_DIR := cmd
CMDS := $(notdir $(wildcard $(CMD_DIR)/*))
UTIL ?=

.PHONY: help list fmt build build-all clean

help:
	@echo "Targets:"
	@echo "  make list                 # list available CLI commands"
	@echo "  make fmt                  # run gofmt on all Go source files"
	@echo "  make build UTIL=<name>    # build one CLI into ./bin"
	@echo "  make build-all            # build all CLIs into ./bin"
	@echo "  make clean                # remove built binaries from ./bin"

list:
	@printf "%s\n" $(CMDS)

fmt:
	gofmt -w $$(find $(CMD_DIR) -name '*.go')

build:
	@if [ -z "$(UTIL)" ]; then \
		echo "Usage: make build UTIL=<name>"; \
		echo "Available: $(CMDS)"; \
		exit 1; \
	fi
	@if [ ! -d "$(CMD_DIR)/$(UTIL)" ]; then \
		echo "Unknown util: $(UTIL)"; \
		echo "Available: $(CMDS)"; \
		exit 1; \
	fi
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(UTIL) ./$(CMD_DIR)/$(UTIL)
	@echo "Built $(BIN_DIR)/$(UTIL)"

build-all:
	@mkdir -p $(BIN_DIR)
	@set -e; for c in $(CMDS); do \
		go build -o $(BIN_DIR)/$$c ./$(CMD_DIR)/$$c; \
		echo "Built $(BIN_DIR)/$$c"; \
	done

clean:
	@set -e; for c in $(CMDS); do \
		rm -f $(BIN_DIR)/$$c; \
	done
	@echo "Removed built binaries for: $(CMDS)"
