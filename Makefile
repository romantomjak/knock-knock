SHELL = bash
PROJECT_ROOT := $(patsubst %/,%,$(dir $(abspath $(lastword $(MAKEFILE_LIST)))))
GIT_COMMIT := $(shell git rev-parse HEAD)
GO_PKGS := $(shell go list ./...)

.PHONY: build
build:
	go build -o knock-knock

.PHONY: test
test:
	go test -cover $(GO_PKGS)

.PHONY: consul
consul:
	consul agent -dev

.PHONY: vault
vault:
	vault server -dev -dev-root-token-id="root"

.PHONY: testdata
testdata:
	consul kv put services/myservice/db/host myexampledb.a1b2c3d4wxyz.us-west-2.rds.amazonaws.com
	consul kv put services/myservice/db/database awsdatabase
	vault kv put secret/services/myservice/db username=awsuser password=awssecretpassword

.PHONY: clean
clean:
	@rm -f "$(PROJECT_ROOT)/knock-knock"
