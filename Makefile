SHELL := /bin/bash

# APP Info
NAME := favorite-apps-tray
APP_VERSION := 11.2.2
DISPLAY_NAME := "Favorite Apps Tray"
# Make file data
GO := go
ROOT := $(CURDIR)
SCRIPTS_DIR := $(ROOT)/scripts
SO_TYPE := "linux"
RELEASE := "1"

.PHONY: all build deploy check-deps clean help

help:
	@echo "Usage: make [target]"
	@echo
	@echo "Targets:"
	@echo "  build						Build windows and linux binaries"
	@echo "  deploy						Generate installer packages"
	@echo "  clean						Remove build outputs"
	@echo "  check-deps					Verify required tools (go)"
	@echo

check-deps:
	@command -v $(GO) >/dev/null 2>&1 || { echo "[ERROR] Please install golang!"; exit 1; }

build: check-deps
	@cd "$(ROOT)"
	@bash $(SCRIPTS_DIR)/build.sh "$(NAME)" "$(RELEASE)"

deploy: check-deps
	@cd "$(ROOT)"
	@bash $(SCRIPTS_DIR)/deploy.sh "$(SO_TYPE)" "$(NAME)" "$(APP_VERSION)" "$(DISPLAY_NAME)"

clean:
	@bash $(SCRIPTS_DIR)/cleaner.sh
