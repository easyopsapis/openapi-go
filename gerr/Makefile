regenerate:
	go install github.com/gogo/protobuf/protoc-gen-gogo
	protoc \
	--proto_path=../../../../ \
	--gogo_out=../../../../ \
	github.com/easyopsapis/openapi-go/gerr/message.proto \
