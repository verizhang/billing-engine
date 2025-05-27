# billing-engine

## How to run
1. Make sure your device can generating proto file
```shell
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
2. Install git tools (https://git-scm.com/downloads)
3. Install PostgreSQL on your device
4. Run all migration inside migrations folder sequentially
5. Setup the environment on run-dev.sh
6. run the run-dev.sh file
```shell
sh run-dev.sh
```