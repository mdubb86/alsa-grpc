To generate protobuf/grpc files
```shell script
protoc --go_out=. --go-grpc_out=. -I ../ alsamixer.proto
```