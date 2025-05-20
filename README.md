# 1. 进行 protoc-gen-go-grpc 的安装

```
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

# 2. 既生成 pb.go 又生成 grpc.pb.go

```
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative  .\protobuf\hello.proto
```