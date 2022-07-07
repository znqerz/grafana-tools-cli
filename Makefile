# ==============================================================================
# Includes

include scripts/make-rules/common.mk

# ==============================================================================
# Usage



##@ General
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)



##@ ZSH command
zsh_deploy: ## Deploy to zsh linked folder
	CGO_ENABLED=0 $(GO) build -v -o ${DIST_DIR}/${PROJECTNAME} ./
	chmod a+x ${DIST_DIR}/${PROJECTNAME}
	mkdir -p ${PROJECT_HOME_PATH}
	${DIST_DIR}/${PROJECTNAME} doc generate
	${DIST_DIR}/${PROJECTNAME} completion zsh > ${COMPLETION_FILE_PATH}
	cp -R ${DIST_DIR}/${PROJECTNAME} ${PROJECT_HOME_PATH}/

zsh_source: ## Source file to zshrc
	find $(HOME)/.zshrc -name '.zshrc' -exec grep -L $(COMPLETION_FILE_PATH) '{}' ';' -print | xargs -I {} sh -c "echo 'source $(COMPLETION_FILE_PATH)' >> {}"


##@ Release
grafana-tools-cli: ## Make grafana-tools-cli
	CGO_ENABLED=0 $(GO) build -v -o ${DIST_DIR}/${PROJECTNAME} ./

