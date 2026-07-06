.PHONY: proto

# Находим все файлы с расширением .proto внутри папки api/proto/
PROTO_FILES := $(shell find api/proto -name "*.proto")

proto:
	@if [ -z "$(PROTO_FILES)" ]; then \
		echo "No .proto files found!"; \
	else \
		protoc --go_out=. --go_opt=paths=source_relative \
		       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
		       $(PROTO_FILES); \
		echo "Successfully compiled all proto files."; \
	fi