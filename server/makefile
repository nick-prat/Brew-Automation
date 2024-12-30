proto: ../protobuf/brewautomation.proto
	protoc \
    --proto_path=../protobuf/ \
    --go_out=../server/pb --go_opt=paths=source_relative \
    --go-grpc_out=../server/pb --go-grpc_opt=paths=source_relative \
    ../protobuf/brewautomation.proto