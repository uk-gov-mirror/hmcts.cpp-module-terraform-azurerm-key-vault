SHELL := /bin/bash

TERRATEST_DIR := tests
TERRATEST_TIMEOUT ?= 30

.PHONY: \
	az-login \
	terratest-test \
	terratest-test-gotestsum \
	terratest-all \
	terratest-all-gotestsum \
	terratest-role-assignments \
	terratest-role-assignments-gotestsum

az-login:
	@az account show >/dev/null 2>&1 || az login

terratest-all: az-login
	cd $(TERRATEST_DIR) && go test -v -timeout $(TERRATEST_TIMEOUT)m ./terratest

terratest-all-gotestsum: az-login
	cd $(TERRATEST_DIR) && go run gotest.tools/gotestsum@v1.12.0 -- -v -timeout $(TERRATEST_TIMEOUT)m ./terratest

# Universal target: pass TEST_NAME as a regex pattern for -run, e.g. TEST_NAME='^TestRoleAssignments_'
terratest-test: az-login
	cd $(TERRATEST_DIR) && if [ -n "$(TEST_NAME)" ]; then \
		go test -v -timeout $(TERRATEST_TIMEOUT)m ./terratest -run "$(TEST_NAME)"; \
	else \
		go test -v -timeout $(TERRATEST_TIMEOUT)m ./terratest; \
	fi

terratest-test-gotestsum: az-login
	cd $(TERRATEST_DIR) && if [ -n "$(TEST_NAME)" ]; then \
		go run gotest.tools/gotestsum@v1.12.0 -- -v -timeout $(TERRATEST_TIMEOUT)m ./terratest -run "$(TEST_NAME)"; \
	else \
		go run gotest.tools/gotestsum@v1.12.0 -- -v -timeout $(TERRATEST_TIMEOUT)m ./terratest; \
	fi

terratest-role-assignments:
	$(MAKE) terratest-test TEST_NAME='^TestRoleAssignments_'

terratest-role-assignments-gotestsum:
	$(MAKE) terratest-test-gotestsum TEST_NAME='^TestRoleAssignments_'
