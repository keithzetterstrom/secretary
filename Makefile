ROOT_PATH := $(PWD)
BIN_PATH := bin

.PHONY: \
	all \
	run \
	devrun \
	build_secretary \
	build_secretary_local \
	tools \
	imports \
	vendor

all: build_secretary

run:
	$(BIN_PATH)/secretary -config=configs/config.yaml -env=prod

devrun: build_secretary_local
	$(BIN_PATH)/secretary -config=configs/config.yaml \
	-env=local \
 	-env_file=./.env \
	-addr=127.0.0.1 \

build_secretary:
	go build -o $(BIN_PATH)/secretary ./cmd/

build_secretary_local:
	go build -tags dynamic -o $(BIN_PATH)/secretary ./cmd/

tools:
	go mod download golang.org/x/tools
	go install golang.org/x/tools/cmd/goimports@latest

imports:
	goimports -l -w .

vendor:
	go mod tidy && go mod vendor
