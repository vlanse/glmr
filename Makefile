LOCAL_BIN := $(CURDIR)/bin
GEN_DIR="./proto"

.PHONY:
buf-deps: deps
	PATH="$(PATH):$(LOCAL_BIN)" go tool buf dep update proto

.PHONY:
deps:
	go mod tidy
	GOBIN=$(LOCAL_BIN) go install tool

.PHONY:
generate/proto: deps
	PATH="$(PATH):$(LOCAL_BIN)" go tool buf generate

.PHONY:
generate: generate/proto
