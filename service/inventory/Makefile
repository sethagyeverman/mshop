GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)
VERSION=$(shell git describe --tags --always)


INTERNAL_PROTO_FILES=$(shell find internal -name *.proto)
API_PROTO_FILES=$(shell find api -name *.proto)



.PHONY: config
# generate internal proto
config:
	protoc --proto_path=./internal \
	       --proto_path=./../../third_party \
 	       --go_out=paths=source_relative:./internal \
	       $(INTERNAL_PROTO_FILES)


.PHONY: api
# generate api proto
api:
	protoc --proto_path=./api \
	       --proto_path=./../../third_party \
 	       --go_out=paths=source_relative:./api \
 	       --go-http_out=paths=source_relative:./api \
 	       --go-grpc_out=paths=source_relative:./api \
	       --openapi_out=fq_schema_naming=true,default_response=false:. \
	       $(API_PROTO_FILES)


.PHONY: build
# build
build:
	mkdir -p bin/ && go build -ldflags "-X main.Version=$(VERSION)" -o ./bin/ ./...


.PHONY: generate
# generate
generate:
	go generate ./...
	go mod tidy


.PHONY: all
# generate all
all:
	make api;
	make config;
	make generate;
