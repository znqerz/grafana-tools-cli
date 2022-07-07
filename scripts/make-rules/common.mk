PROJECTNAME  := grafana-tools-cli
MAKE 		 :=  make
PROJECT_HOME_PATH := $(HOME)/.${PROJECTNAME}
COMPLETION_FILE_PATH := $(PROJECT_HOME_PATH)/.${PROJECTNAME}-completion.sh

GO           := go
GOFMT        := $(GO)fmt
PKGS         := $(shell $(GO) list ./... | grep -v -e "vendor")


# include the common make file
COMMON_SELF_DIR = $(dir $(lastword $(MAKEFILE_LIST)))

ifeq ($(origin ROOT_DIR),undefined)
ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR)/../.. && pwd -P))
endif

GRAFANA_SECRETS=$(ROOT_DIR)/deployments/secrets
DIST_DIR     := ${ROOT_DIR}/dist
CLI_BIN 	 := ${ROOT_DIR}/dist/$(PROJECTNAME)


build.bin:
	CGO_ENABLED=0 $(GO) build -v -i -o ${DIST_DIR}/${PROJECTNAME} ./
	@chmod a+x ${DIST_DIR}/${PROJECTNAME}