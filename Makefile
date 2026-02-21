.PHONY: gen

# Generate Go and gRPC code from proto. 
# Input: api/proto/entity.proto
# output: pb/
gen:
	protoc --proto_path=api/proto \
		--go_out=. --go_opt=module=crud \
		--go-grpc_out=. --go-grpc_opt=module=crud \
		api/proto/entity.proto