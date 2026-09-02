.DEFAULT_GOAL := help
.PHONY: help setup validate lint test test-conformance check-contracts docs generate clean

# BongoPay root developer interface.
#
# This Makefile is the STABLE contract for contributors and CI, regardless of what
# language(s) implementations/adapters/SDKs end up using underneath. Targets delegate to
# scripts/ so the interface here never has to change just because tooling does.
#
# See docs/development/README.md for the full local development guide and
# AGENTS.md for what AI agents may/should run.

help: ## Show this help
	@echo "BongoPay developer commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'

setup: ## Install local validation tooling (no heavy language toolchains required)
	@bash scripts/setup.sh

validate: ## Validate specs, schemas, and contracts
	@bash scripts/validate-specs.sh
	@bash scripts/validate-schemas.sh
	@bash scripts/validate-contracts.sh

lint: ## Lint markdown, YAML, and JSON across the repository
	@bash scripts/lint.sh

test: ## Run available unit/integration tests
	@bash scripts/test.sh

test-conformance: ## Run the conformance suite(s) as they come online
	@bash scripts/test-conformance.sh

check-contracts: ## Verify generated contracts match their source specs
	@bash scripts/check-contracts.sh

docs: ## Build/check documentation (link checks, structure checks)
	@bash scripts/docs.sh

generate: ## Regenerate generated artifacts from source specs/contracts
	@bash scripts/generate.sh

clean: ## Remove local build/validation artifacts
	@bash scripts/clean.sh
