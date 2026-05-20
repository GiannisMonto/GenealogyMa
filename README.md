# 族谱数字化管理平台

> 基于现代技术栈的族谱数字化、可视化、传承一体化解决方案

[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/your-org/genealogy-ma/blob/main/LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)]()
[![Git Commit Standard](https://img.shields.io/badge/commit-Conventional%20Commits-brightgreen)]()

## ✨ 项目特色

- 🔍 **5950+ 人物数据** - 完整的族谱数据，支持高效查询
- 🌳 **族谱树可视化** - 基于 D3.js 的交互式世系图展示
- 🌅 **旭日图分析** - 直观的族群世代分布可视化
- 🏛️ **3D 宗祠系统** - Three.js 沉浸式宗祠漫游体验
- 📜 **家族文化传承** - 文献库、故事馆、家训堂
- 📱 **多端支持** - Web、移动端、微信小程序

---

## 📁 项目结构 (Monorepo)

```
GenealogyMa/
├── backend/                          # 后端服务 (Go + DDD)
│   ├── cmd/                         # 应用入口
│   ├── internal/                    # 内部代码
│   │   ├── domain/                  # 领域层（核心业务）
│   │   ├── application/             # 应用层（业务编排）
│   │   ├── infrastructure/          # 基础设施（DB、缓存）
│   │   └── interfaces/              # 接口层（HTTP、WebSocket）
│   ├── pkg/                         # 公共工具包
│   └── README.md
│
├── frontend/                         # 前端项目
│   ├── admin/                       # 管理后台 (React + AntD)
│   ├── visualization/               # 可视化前端 (Vue 3 + D3.js)
│   ├── shared/                      # 共享组件 / 类型
│   └── README.md
│
├── mobile/                           # 移动端
│   ├── miniprogram/                 # 微信小程序 (uni-app)
│   └── README.md
│
├── docs/                             # 项目文档
│   ├── GIT_WORKFLOW.md              # Git 工作流规范
│   ├── API.md                       # API 文档
│   ├── DATABASE.md                  # 数据库设计
│   └── DEPLOYMENT.md                # 部署文档
│
├── deployments/                      # 部署配置
│   ├── docker/                      # Docker 相关
│   ├── k8s/                         # Kubernetes 配置
│   └── scripts/                     # 部署脚本
│
├── scripts/                          # 项目脚本
│   ├── setup.sh                     # 环境初始化
│   └── deploy.sh                    # 自动部署
│
├── .github/                          # GitHub 配置
│   └── workflows/                   # GitHub Actions CI/CD
│
├── .gitignore                        # Git 忽略规则
├── .gitattributes                    # Git 属性配置
├── .env.example                      # 环境变量示例
├── package.json                      # 根级 NPM 脚本
└── README.md                         # 项目说明（本文件）
```

---

## 🛠️ 技术栈

### 后端
| 技术 | 说明 | 版本 |
|------|------|------|
| <img src="https://img.shields.io/badge/Go-00ADD8?logo=go&logoColor=white"/> | 编程语言 | 1.21+ |
| <img src="https://img.shields.io/badge/Gin-00ADD8?logo=gin&logoColor=white"/> | Web 框架 | v1.9.1 |
| <img src="https://img.shields.io/badge/GORM-25B86E?logoColor=white"/> | ORM 框架 | v1.25.5 |
| <img src="https://img.shields.io/badge/PostgreSQL-316192?logo=postgresql&logoColor=white"/> | 数据库 | 15+ |
| <img src="https://img.shields.io/badge/Redis-DC382D?logo=redis&logoColor=white"/> | 缓存 | 7+ |
| <img src="https://img.shields.io/badge/JWT-000000?logo=json-web-tokens&logoColor=white"/> | 认证 | v5 |
| <img src="https://img.shields.io/badge/Swagger-85EA2D?logo=swagger&logoColor=black"/> | API 文档 | v1.16.2 |

### 前端
| 技术 | 说明 | 应用场景 |
|------|------|----------|
| React 18 + TypeScript | 前端框架 | 管理后台 |
| Ant Design 5 | UI 组件库 | 管理后台 |
| Vue 3 + TypeScript | 前端框架 | 可视化前端 |
| Element Plus | UI 组件库 | 可视化前端 |
| D3.js | 可视化库 | 族谱树、旭日图 |
| ECharts | 图表库 | 统计图表 |
| Three.js | 3D 引擎 | 宗祠漫游 |
| uni-app + Vue 3 | 跨端框架 | 微信小程序 |

---

## 🚀 快速开始

### 环境要求

- **Go** >= 1.21.0
- **Node.js** >= 18.0.0
- **PostgreSQL** >= 15 (with ltree extension)
- **Redis** >= 7.0

### 1. 克隆项目

```bash
git clone https://github.com/your-org/genealogy-ma.git
cd genealogy-ma
```

### 2. 后端启动

```bash
cd backend

# 复制环境变量
cp .env.example .env
# 编辑 .env 配置数据库连接

# 安装依赖
go mod download

# 运行服务
go run cmd/api/main.go
```

后端服务将启动在 `http://localhost:8080`

API 文档：`http://localhost:8080/swagger/index.html`

### 3. 前端启动 (可视化)

```bash
cd frontend/visualization

npm install
npm run dev
```

### 4. 管理后台启动

```bash
cd frontend/admin

npm install
npm run dev
```

### 使用根级脚本（推荐）

```bash
# 同时启动后端和可视化前端
npm run dev

# 单独启动
npm run dev:backend
npm run dev:visualization
npm run dev:admin

# 构建所有
npm run build
```

---

## 🎯 核心功能模块

### 1. 族谱阅览
- ✅ 人物 CRUD 管理
- ✅ 多维度搜索（姓名、字号、生卒年）
- ✅ 上下五代世系图
- ✅ 批量导入导出
- ✅ 数据版本控制

### 2. 族群可视化
- ✅ 族谱树 D3.js 可视化
- ✅ 旭日图代际分布
- ✅ 迁徙路径地图
- ✅ 人口统计分析

### 3. 3D 宗祠
- ✅ Three.js 3D 场景渲染
- ✅ 第一人称沉浸式漫游
- ✅ 虚拟牌位供奉系统
- ✅ 在线祭祀功能

### 4. 家族文化
- ✅ 族谱文献库
- ✅ 家族故事馆
- ✅ 家训家规展示
- ✅ 家族名人堂

### 5. 墓葬管理
- 📋 墓园 GIS 地图
- 📋 墓位信息管理
- 📋 祭扫导航

### 6. 交流社区
- 📋 家族动态流
- 📋 寻根问祖板块
- 📋 族人通讯录
- 📋 私信功能

---

## 📊 数据统计

| 指标 | 数值 |
|------|------|
| 总人数 | 5,950 |
| 世代跨度 | 0 - 22 代 |
| 配偶数量 | 2,166 |
| 年均新增 | 待统计 |
| 最旺世代 | 第 19 代 |

---

## 📖 文档

| 文档 | 说明 |
|------|------|
| [Git 工作流规范](./docs/GIT_WORKFLOW.md) | 分支策略、提交规范、PR 流程 |
| [后端 README](./backend/README.md) | 后端开发指南 |
| [前端 README](./frontend/README.md) | 前端开发指南 |
| [移动端 README](./mobile/README.md) | 小程序开发指南 |

---

## 🔧 开发命令

```bash
# ===== 后端 =====
cd backend
go run cmd/api/main.go          # 启动服务
go test ./... -v                 # 运行测试
go test ./... -cover            # 测试覆盖率
golangci-lint run                # 代码检查
swag init -g cmd/api/main.go    # 生成 Swagger 文档

# ===== 前端 =====
cd frontend/visualization
npm run dev                      # 开发模式
npm run build                    # 生产构建
npm run lint                     # 代码检查

# ===== 根级快捷命令 =====
npm run dev                      # 启动前后端
npm run build                    # 构建所有
npm run test                     # 运行所有测试
npm run lint                     # 代码检查
npm run docker:up                # Docker 启动
```

---

## 🌿 Git 工作流

本项目采用 **Git Flow** 改良工作流：

```
main --------------------------○---------------- 生产分支
         \              /
develop -------------○-----------○---------- 开发分支
         \         /     \
feature   ---○---○---○---○--- 功能分支
```

**核心分支**：
- `main`: 生产环境代码，受保护
- `develop`: 开发分支，最新状态

**临时分支**：
- `feature/*`: 功能开发
- `bugfix/*`: Bug 修复
- `hotfix/*`: 紧急生产修复
- `release/*`: 版本发布准备

**提交规范**：
```
<type>(<scope>): <subject>

# 示例
feat(person): add fuzzy search by name
fix(tree): resolve render issue
docs: update API documentation
```

详细规范参见 [Git 工作流规范](./docs/GIT_WORKFLOW.md)

---

## 🤝 贡献指南

我们非常欢迎你的贡献！

### 贡献流程

1. **Fork** 本仓库
2. 从 `develop` 创建功能分支: `git checkout -b feature/your-feature`
3. 提交你的变更: `git commit -m 'feat: add some feature'`
4. 推送到分支: `git push origin feature/your-feature`
5. 发起 **Pull Request** 到 `develop` 分支

### PR 要求

- ✅ 遵循项目代码规范
- ✅ 提交信息符合 [Conventional Commits](https://www.conventionalcommits.org/)
- ✅ 相关测试已通过
- ✅ 文档已更新（如需要）

---

## 📄 许可证

[MIT License](LICENSE)

---

## 👥 团队

| 角色 | 职责 |
|------|------|
| 后端开发 | Go + DDD |
| 前端开发 | React + Vue |
| 3D 开发 | Three.js + 建模 |
| UI/UX | 设计与交互 |

---

## 📞 联系方式

- 项目地址: [GitHub](https://github.com/your-org/genealogy-ma)
- Issue 反馈: [Issues](https://github.com/your-org/genealogy-ma/issues)

---

<div align="center">
  <sub>Built with ❤️ by Genealogy Ma Team</sub>
</div>
