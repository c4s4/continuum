NAME="continuum"
VERSION=$(shell changelog release version)
BUILD_DIR="build"

YELLOW=\033[1m\033[93m
CYAN=\033[1m\033[96m
CLEAR=\033[0m

.PHONY: build

help:
	@echo "$(YELLOW)Help page$(CLEAR)"
	@echo "$(CYAN)help$(CLEAR)     Print this help page"
	@echo "$(CYAN)deps$(CLEAR)     Install GO dependencies"
	@echo "$(CYAN)build$(CLEAR)    Build the application"
	@echo "$(CYAN)archive$(CLEAR)  Build the distribution archive"
	@echo "$(CYAN)release$(CLEAR)  Release project" 
	@echo "$(CYAN)clean$(CLEAR)    Clean generated file"

deps:
	@echo "$(YELLOW)Installing GO dependencies$(CLEAR)"
	go get gopkg.in/yaml.v1
	go get github.com/mitchellh/gox

test:
	@echo "$(YELLOW)Running tests$(CLEAR)"
	go test

build: clean test
	@echo "$(YELLOW)Building application$(CLEAR)"
	mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(NAME)

archive: build
	@echo "$(YELLOW)Building distribution archive$(CLEAR)"
	mkdir -p $(BUILD_DIR)/$(NAME)-$(VERSION)/bin/
	gox -output=$(BUILD_DIR)/$(NAME)-$(VERSION)/bin/{{.Dir}}_{{.OS}}_{{.Arch}}
	mkdir -p $(BUILD_DIR)/$(NAME)-$(VERSION)/etc/
	cp continuum.yml $(BUILD_DIR)/$(NAME)-$(VERSION)/etc/
	cp LICENSE.txt $(BUILD_DIR)/$(NAME)-$(VERSION)/
	cp README.md $(BUILD_DIR)/ && cd $(BUILD_DIR) && md2pdf README.md && cp README.pdf $(NAME)-$(VERSION)/
	cp CHANGELOG.yml $(BUILD_DIR)/ && cd $(BUILD_DIR) && changelog to html style > $(NAME)-$(VERSION)/CHANGELOG.html
	cd $(BUILD_DIR) && tar cvzf $(NAME)-bin-$(VERSION).tar.gz $(NAME)-$(VERSION)

release: archive
	@echo "$(YELLOW)Releasing project$(CLEAR)"
	release

clean:
	@echo "$(YELLOW)Cleaning generated files$(CLEAR)"
	rm -rf $(BUILD_DIR)
