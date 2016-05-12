NAME="gontinuum"
VERSION=$(shell changelog release version)
BUILD_DIR="build"

YELLOW=\033[1m\033[93m
CYAN=\033[1m\033[96m
CLEAR=\033[0m

.PHONY: build

help:
	@echo "$(YELLOW)Help page$(CLEAR)"
	@echo "$(CYAN)help$(CLEAR)   Print this help page"
	@echo "$(CYAN)deps$(CLEAR)   Install GO dependencies"
	@echo "$(CYAN)build$(CLEAR)  Build the application"
	@echo "$(CYAN)clean$(CLEAR)  Clean generated file"

deps:
	go get gopkg.in/yaml.v1

build:
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(NAME)

clean:
	rm -rf $(BUILD_DIR)
