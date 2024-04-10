# SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

TOOLS_BIN_DIR              := $(TOOLS_DIR)/bin
SKAFFOLD                   := $(TOOLS_BIN_DIR)/skaffold
KUSTOMIZE                  := $(TOOLS_BIN_DIR)/kustomize
GOLANGCI_LINT              := $(TOOLS_BIN_DIR)/golangci-lint
GOIMPORTS                  := $(TOOLS_BIN_DIR)/goimports

# default tool versions
SKAFFOLD_VERSION ?= v1.38.0
GOLANGCI_LINT_VERSION ?= v1.57.2
KUSTOMIZE_VERSION ?= v5.0.0

export TOOLS_BIN_DIR := $(TOOLS_BIN_DIR)
export PATH := $(abspath $(TOOLS_BIN_DIR)):$(PATH)

#########################################
# Tools                                 #
#########################################

$(KUSTOMIZE):
	@test -s $(TOOLS_BIN_DIR)/kustomize || GOBIN=$(abspath $(TOOLS_BIN_DIR)) go install sigs.k8s.io/kustomize/kustomize/v5@${KUSTOMIZE_VERSION}

$(GOIMPORTS):
	go build -o $(GOIMPORTS) golang.org/x/tools/cmd/goimports

$(GOLANGCI_LINT):
	@# CGO_ENABLED has to be set to 1 in order for golangci-lint to be able to load plugins
	@# see https://github.com/golangci/golangci-lint/issues/1276
	GOBIN=$(abspath $(TOOLS_BIN_DIR)) CGO_ENABLED=1 go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
