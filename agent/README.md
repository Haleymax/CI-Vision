# CI Agent

CI-Vision 项目的 Python 智能代理模块

## 项目结构

```
ci_agent/
├── .venv/                    # Python虚拟环境
├── __pycache__/             # Python字节码缓存
├── config/                  # 配置文件
├── docs/                    # 项目文档
├── logs/                    # 日志文件目录
├── src/                     # 源代码
│   ├── grpc/               # gRPC通信模块
│   │   ├── client.py       # gRPC客户端
│   │   ├── service.py      # gRPC服务端
│   │   └── proto/          # Protocol Buffers生成文件
│   ├── agent/              # AI智能代理模块
│   ├── models/             # 数据模型
│   ├── services/           # 业务服务层
│   └── utils/              # 工具函数
├── tests/                   # 测试代码
│   ├── test_agent/         # AI代理测试
│   └── test_grpc/          # gRPC通信测试
├── requirements.txt        # Python依赖包
└── README.md              # 项目说明文档
```
