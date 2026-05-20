# 族谱管理平台 - 后端服务

基于 Go + Gin + GORM + DDD 领域驱动设计的后端服务

## 技术栈

| 技术 | 说明 | 版本 |
|------|------|------|
| Go | 编程语言 | 1.21+ |
| Gin | Web 框架 | v1.9.1 |
| GORM | ORM 框架 | v1.25.5 |
| PostgreSQL | 数据库 | 15+ |
| Redis | 缓存 | 7+ |
| JWT | 认证 | v5 |
| Swagger | API 文档 | v1.16.2 |

## 架构设计

采用 DDD (领域驱动设计) 分层架构

```
backend/
├── cmd/
│   ├── api/              # API 服务入口
│   └── migrate/          # 数据库迁移工具
├── internal/
│   ├── domain/           # 领域层（核心业务逻辑）
│   │   ├── person/      # 人物聚合根
│   │   │   ├── entity.go
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   ├── genealogy/   # 族谱领域
│   │   ├── culture/     # 文化领域
│   │   ├── memorial/    # 祭祀领域
│   │   ├── cemetery/    # 墓葬领域
│   │   └── community/   # 社区领域
│   ├── application/      # 应用层（业务编排）
│   │   ├── dto/         # 数据传输对象
│   │   └── service/     # 应用服务
│   ├── infrastructure/   # 基础设施层
│   │   ├── persistence/ # 数据库实现
│   │   ├── cache/       # Redis 缓存
│   │   └── config/      # 配置
│   └── interfaces/       # 接口层
│       ├── http/
│       │   ├── controller/  # HTTP 控制器
│       │   └── middleware/  # 中间件
│       └── ws/         # WebSocket
└── pkg/                 # 公共工具包
    ├── auth/
    ├── cache/
    ├── utils/
    └── errors/
```

## 快速开始

### 环境要求

- Go 1.21+
- PostgreSQL 15+ (ltree extension)
- Redis 7+

### 配置

```bash
# 复制环境变量
cp .env.example .env

# 编辑配置
vim .env
```

### 运行

```bash
# 安装依赖
go mod download

# 运行服务
go run cmd/api/main.go

# 热重载开发 (需要 air)
air -c .air.toml
```

服务将启动在 `http://localhost:8080`

## API 文档

启动服务后访问 Swagger 文档：

```
http://localhost:8080/swagger/index.html
```

## 核心 API

### 人物模块

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| GET | `/api/v1/persons` | 搜索人物列表 | 公开 |
| GET | `/api/v1/persons/:id` | 获取人物详情 | 公开 |
| GET | `/api/v1/persons/:id/tree` | 获取族谱树 | 公开 |
| GET | `/api/v1/persons/statistics` | 获取统计数据 | 公开 |
| POST | `/api/v1/persons` | 创建人物 | 需要认证 |
| PUT | `/api/v1/persons/:id` | 更新人物 | 需要认证 |
| DELETE | `/api/v1/persons/:id` | 删除人物 | 需要认证 |

## 数据库

### 核心表

- `members` - 人物主表 (5950 条数据)
- `parent_child_relations` - 父子关系表
- `spouses` - 配偶信息表
- `sections` - 族谱章节表
- `source_references` - 资料来源表

### ltree 索引优化

使用 PostgreSQL ltree 扩展实现高效族谱查询：

```sql
-- 创建 GIST 索引
CREATE INDEX idx_members_ltree_path ON members USING GIST(lineage_path);

-- 查询后代
SELECT * FROM members WHERE lineage_path <@ '1.2.3';

-- 查询祖先
SELECT * FROM members WHERE lineage_path @> '1.2.3.4';
```

## 开发命令

```bash
# 运行测试
go test ./... -v

# 运行测试并生成覆盖率
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# 代码格式检查
gofmt -l .

# 代码质量检查 (需要 golangci-lint)
golangci-lint run

# 生成 Swagger 文档
swag init -g cmd/api/main.go

# 构建
go build -o bin/api cmd/api/main.go
```

## 代码规范

- 遵循 [Go 官方代码规范](https://go.dev/doc/effective_go)
- 包名使用小写，不使用下划线
- 导出函数必须有注释
- 错误处理：始终检查错误，尽早返回
- 优先使用标准库

## 目录说明

```bash
# 后端完整目录
backend/
├── cmd/                 # 可执行程序入口
├── internal/            # 内部代码（不对外暴露）
│   ├── domain/         # 领域层：实体、仓储接口、领域服务
│   ├── application/    # 应用层：DTO、应用服务、命令/查询
│   ├── infrastructure/ # 基础设施层：数据库、缓存、外部服务
│   └── interfaces/     # 接口层：HTTP控制器、中间件、路由
├── pkg/                 # 可复用公共包
├── test/               # 测试辅助代码
├── .air.toml           # air 热重载配置
├── go.mod
└── README.md
```
