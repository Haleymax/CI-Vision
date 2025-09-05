# Protocol Buffers 代码生成

## 生成 Go 代码

在项目根目录下执行以下命令来生成 Go 代码：

```bash
# 生成 Go 代码到 backend/proto 目录
protoc --go_out=backend --go_opt=paths=source_relative \
       --go-grpc_out=backend --go-grpc_opt=paths=source_relative \
       proto/hello.proto
```

## 生成 Python 代码

在项目根目录下执行以下命令来生成 Python 代码：

```bash
# 生成 Python 代码到 ci_agent 目录
python -m grpc_tools.protoc -I proto --python_out=ci_agent --grpc_python_out=ci_agent proto/hello.proto
```

## 说明

- Go 代码会生成到 `backend/proto/` 目录
- Python 代码会生成到 `ci_agent/` 目录
- 确保已安装相应的工具：
  - Go: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
  - Go gRPC: `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`
  - Python: `pip install grpcio-tools`
