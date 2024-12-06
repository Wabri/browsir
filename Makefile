.PHONY: all build clean install

BINARY_NAME=browsir
INSTALL_PATH=/usr/local/bin
CONFIG_PATH=/etc/browsir

all: build

build:
	go build -o $(BINARY_NAME)

clean:
	go clean
	rm -f $(BINARY_NAME)

install: build
	sudo mkdir -p $(CONFIG_PATH)
	sudo ln -sf "$(PWD)/.browsir.yml" "$(CONFIG_PATH)/config.yml" || true
	sudo ln -sf "$(PWD)/shortcuts" "$(CONFIG_PATH)/shortcuts" || true
	sudo ln -sf "$(PWD)/$(BINARY_NAME)" "$(INSTALL_PATH)/$(BINARY_NAME)"
	@echo "Created symlink to $(BINARY_NAME) in $(INSTALL_PATH)"
	@echo "Created symlinks to config files in $(CONFIG_PATH)"
	@echo "You can now run 'browsir' from anywhere"
