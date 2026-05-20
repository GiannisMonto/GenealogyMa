# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

---

## 项目概览

**族谱数字化管理平台**

- 核心数据：5,950+ 人物，22 代世系
- 架构：Monorepo + DDD 后端
- 主要功能：族谱树、3D 宗祠、族群可视化、家族文化传承、墓葬管理、交流社区

---

## 架构总览

### Monorepo 结构
- `backend/` - Go + DDD 后端服务
- `frontend/admin/` - React + AntD 管理后台
- `frontend/visualization/` - Vue 3 + D3.js 可视化前端
- `mobile/miniprogram/` - uni-app 微信小程序

### DDD 四层架构（关键设计）

```
HTTP Request → Controller → Application Service → Domain Service → Repository → PostgreSQL
```

**层级职责**：
1. **Domain Layer** (`internal/domain/`) - 领域层：实体定义、仓储接口、业务规则、领域服务
2. **Application Layer** (`internal/application/`) - 应用层：DTO、服务编排、事务边界
3. **Infrastructure Layer** (`internal/infrastructure/`) - 基础设施层：数据库、缓存、配置实现
4. **Interface Layer** (`internal/interfaces/`) - 接口层：HTTP 控制器、中间件、路由

---

## 常用开发命令

### 根级 NPM 脚本
```bash
npm run dev              # 同时启动前后端
npm run build            # 构建所有项目
npm run test             # 运行所有测试
npm run lint             # 代码检查
npm run docker:up        # Docker 启动环境
npm run db:migrate       # 数据库迁移
```

### 后端 Go 命令
```bash
# 开发运行
cd backend
go run cmd/api/main.go

# 测试
go test ./... -v         # 运行所有测试
go test ./... -cover     # 测试覆盖率
go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out

# 运行单个包测试
go test ./internal/domain/person/... -v

# 代码质量
gofmt -l .               # 代码格式检查
gofmt -w .               # 自动格式化
golangci-lint run        # 代码质量检查

# 文档
swag init -g cmd/api/main.go    # 生成 Swagger 文档

# 构建
go build -o bin/api cmd/api/main.go
```

---

## 关键文件快速定位

| 功能 | 路径 |
|------|------|
| API 主入口 | `backend/cmd/api/main.go` |
| 人物实体定义 | `backend/internal/domain/person/entity.go` |
| 人物仓储接口 | `backend/internal/domain/person/repository.go` |
| 人物应用服务 | `backend/internal/application/service/person_service.go` |
| 人物控制器 | `backend/internal/interfaces/http/controller/person_controller.go` |
| 数据库连接 | `backend/internal/infrastructure/persistence/postgres.go` |
| 配置加载 | `backend/internal/infrastructure/config/config.go` |
| JWT 认证 | `backend/pkg/auth/jwt.go` |
| 统一响应 | `backend/pkg/utils/response.go` |
| API 文档 | `http://localhost:8080/swagger/index.html` |

---

## 领域模块规划

| 领域 | 路径 | 状态 | 说明 |
|------|------|------|------|
| **Person** | `internal/domain/person/` | ✅ 已实现 | 人物聚合根、CRUD、族谱树 |
| Genealogy | `internal/domain/genealogy/` | ⏳ 待实现 | 族谱、分支、字辈 |
| Culture | `internal/domain/culture/` | ⏳ 待实现 | 文献、故事、家训 |
| Memorial | `internal/domain/memorial/` | ⏳ 待实现 | 3D 宗祠、牌位、祭祀 |
| Cemetery | `internal/domain/cemetery/` | ⏳ 待实现 | 墓园、墓位、GIS |
| Community | `internal/domain/community/` | ⏳ 待实现 | 用户、动态、私信 |

---

## 数据库关键特性

### PostgreSQL ltree 世系路径
```sql
-- 查询某人所有后代
SELECT * FROM members WHERE lineage_path <@ '1.2.3';

-- 查询某人所有祖先
SELECT * FROM members WHERE lineage_path @> '1.2.3.4';
```

**核心表**：
- `members` - 人物主表 (5,950 条)
- `parent_child_relations` - 父子关系表
- `spouses` - 配偶信息表

---

## RBAC 权限模型

**角色**：
- `super_admin` - 超级管理员
- `genealogy_admin` - 族谱管理员
- `culture_admin` - 文化管理员
- `memorial_admin` - 纪念馆管理员
- `verified_user` - 认证用户
- `guest` - 游客

**权限示例**：`person:read`, `person:write`, `person:delete`, `culture:write`, `admin:user`

---

## Git 最佳实践

### Git 工作流

| 策略 | 说明 |
|------|------|
| **Trunk-Based Development** | 高绩效团队首选，短期功能分支 + CI 持续验证 |
| **GitHub Flow** | 轻量级通用方案，main + feature 分支 |
| **Git Flow** | 适用于有计划发布周期的项目（使用逐渐减少） |

### 提交规范 (Conventional Commits v1.0.0)
```
<type>(<scope>): <subject>

# Type: feat, fix, docs, style, refactor, perf, test, build, ci, chore
# Scope: backend, person, auth, frontend, admin, visualization, 3d, docker

# 示例
feat(person): add fuzzy search by name
fix(tree): resolve render issue
```

### Git 安全实践
- GPG/SSH 提交签名配置
- 硬件安全密钥 (YubiKey) 推荐
- 机密扫描防止凭证提交

### Git 性能优化
- 稀疏检出 + 稀疏索引
- 部分克隆 (`--filter=blob:none`)
- 文件系统监视器 (`core.fsmonitor`)
- 定期 `git gc --aggressive`

### 大文件管理
- Git LFS 处理二进制文件
- git filter-repo 清理历史

### CI/CD 集成
- 事件驱动流水线
- 路径过滤器
- 依赖缓存策略
- 合并队列 (Merge Queue)

### 分支保护规则
- PR 审查要求
- 状态检查必须通过
- 强制线性历史
- 提交签名验证

---

## 环境配置

```bash
# 复制环境变量
cp .env.example .env

# 数据库（已预置）
PG_HOST=10.100.34.3
PG_PORT=5432
PG_DBNAME=genealogy
```

---

## 技术栈速查

| 层级 | 技术 | 版本 |
|------|------|------|
| 后端语言 | Go | 1.21+ |
| Web 框架 | Gin | v1.9.1 |
| ORM | GORM | v1.25.5 |
| 数据库 | PostgreSQL | 15+ |
| 缓存 | Redis | 7+ |
| 认证 | JWT | v5 |
| API 文档 | Swagger | v1.16.2 |

---

## Git 最佳实践完整指南

### 分支策略

| 策略 | 说明 |
|------|------|
| **Trunk-Based Development** | 主流推荐，短期功能分支 + CI 持续验证 |
| **GitHub Flow** | 轻量级通用，main + feature 分支 |
| **Git Flow** | 适用于有计划发布周期项目（使用减少） |

### 提交格式 (Conventional Commits v1.0.0)

```
<type>[optional scope]: <description>

# Type: feat, fix, docs, style, refactor, perf, test, build, ci, chore
# ! 表示破坏性变更
```

| 类型 | 说明 | SemVer |
|------|------|--------|
| `feat` | 新功能 | MINOR |
| `fix` | Bug 修复 | PATCH |
| `docs` | 文档更新 | - |
| `style` | 代码格式 | - |
| `refactor` | 重构 | - |
| `perf` | 性能优化 | PATCH |
| `test` | 测试相关 | - |
| `build` | 构建系统变更 | - |
| `ci` | CI/CD 配置 | - |
| `chore` | 其他杂项 | - |

### Git 安全配置

```bash
# GPG 提交签名
git config --global user.signingkey <your-key-id>!
git config --global commit.gpgsign true    # 自动签名所有提交
git config --global tag.gpgsign true       # 自动签名标签

# SSH 签名 (Git 2.34+)
git config --global gpg.format ssh
git config --global user.signingkey ~/.ssh/id_ed25519.pub

# 验证签名
git log --show-signature -1
git merge --verify-signatures -S signed-branch
```

### 性能优化配置

```bash
# 文件系统监视器（大幅提升 Windows/macOS 性能）
git config --global core.fsmonitor true

# 多线程打包
git config --global pack.threads 0

# 仓库维护
git gc --aggressive              # 深度垃圾回收
git fetch --prune               # 清理远程已删除分支引用
git prune                       # 清理无用对象

# 浅层克隆 (CI/CD 常用)
git clone --depth 1 <repo-url>
git fetch --deepen=100         # 按需获取更多历史
```

### 大仓库管理

```bash
# Git LFS
git lfs install
git lfs track "*.psd" "*.bin" "models/*.pth"
git lfs ls-files

# 稀疏检出 (Git 2.52.0+)
git sparse-checkout set --sparse-index src/app src/components
git sparse-checkout add src/utils
git sparse-checkout reapply
git sparse-checkout check-rules src/app/utils.ts

# 部分克隆
git clone --filter=blob:none <repo-url>  # 仅在需要时下载文件内容
git clone --filter=tree:0 <repo-url>     # 排除树对象（超大仓库）
```

### Pull Request 最佳实践

✅ **Do**:
- 每个 PR 小而专注，只完成一个目的
- 清晰的标题和完整描述
- 提交前自审、构建、测试
- 使用关键词关联 Issue：`Closes #123`

❌ **Don't**:
- 超大 PR（> 400 行考虑拆分）
- 混合多个无关变更
- 无描述的 "fix stuff" 式 PR

### 分支保护规则 (GitHub)

- PR 审查要求：至少 N 人批准
- 最新推送批准：防止批准后 PR 被篡改
- 撤销旧批准：PR 变更时自动撤销之前的批准
- 所有讨论必须已解决
- 状态检查必须通过
- 强制线性历史
- 要求提交签名验证
- 规则应用于管理员

### 2024-2025 Git 关键趋势

| 趋势 | 说明 |
|------|------|
| 🤖 **AI 深度集成** | AI 生成提交信息、PR 摘要、代码审查建议 |
| 🔐 **安全左移** | 更严格的签名验证、机密扫描、供应链安全 |
| ⚡ **性能优化** | 稀疏索引、部分克隆、fsmonitor 成熟易用 |
| 🚀 **自动化增强** | 合并队列、自动合并、智能 CI 跳过 |
| 📈 **Trunk-Based 普及** | 更多团队从 Git Flow 转向 TBD |

### 快速参考速查表

```bash
# 配置
git config --global user.name "Your Name"
git config --global user.email "your@email.com"
git config --global commit.gpgsign true

# 分支操作
git checkout -b feature/new-feature
git branch -d old-branch
git push origin --delete old-branch

# 查看历史
git log --oneline --graph --all
git log --show-signature -1

# 清理维护
git fetch --prune
git gc --aggressive
```

### 配置检查清单

- [ ] user.name 和 user.email 正确配置
- [ ] 提交签名（GPG/SSH）已配置并上传公钥
- [ ] 行尾设置正确（`core.autocrlf`）
- [ ] fsmonitor 已启用（Windows/macOS）
- [ ] 分支保护规则已配置
- [ ] CI 状态检查已启用

## Agent要求

- 当上下文变大时，将当前状态写入 tasks/mission.md。包括：已完成的、下一步的、被阻塞的、未解决的问题。
- 错误处理：最多重试 3 次。如果没有进展，记录到 pending_for_human.md 然后转到下一个任务。
- 压缩前，务必保存完整的已修改文件列表。
