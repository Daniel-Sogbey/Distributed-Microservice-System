.PHONY: proto

PROTO_DIR=proto

proto:
	protoc -I . \
	  --go_out=paths=source_relative:. \
	  --go-grpc_out=paths=source_relative:. \
	  $(shell find $(PROTO_DIR) -name '*.proto')
