BINARY_NAME=main.out

proto:
	protoc --go_out=proto/gen/go --go_opt=paths=source_relative \                              ✔  08:14:48 PM  
       --go-grpc_out=proto/gen/go --go-grpc_opt=paths=source_relative \
       proto/api/rec.proto


build:
	go build -o ./bin/${BINARY_NAME}

run:
	go build -o bin/${BINARY_NAME} main.go
	./bin/${BINARY_NAME}