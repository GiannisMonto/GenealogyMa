# Git 工作流规范

## 概述

本文档规定了族谱数字化管理平台的 Git 工作流程、分支策略、提交规范和代码审查流程。

---

## 一、分支策略

### 1.1 核心分支

| 分支名 | 说明 | 保护规则 |
|--------|------|----------|
| `main` | 主分支，生产环境代码，随时可发布 | ✅ 受保护，仅通过 PR 合并 |
| `develop` | 开发分支，最新开发状态 | ✅ 受保护，仅通过 PR 合并 |

### 1.2 短期分支（用完即删）

| 前缀 | 用途 | 命名示例 | 生命周期 |
|------|------|----------|----------|
| `feature/` | 新功能开发 | `feature/person-tree-visualization` | 从 develop 创建，合并回 develop |
| `bugfix/` | 修复 develop 中的 bug | `bugfix/fix-search-performance` | 从 develop 创建，合并回 develop |
| `hotfix/` | 紧急修复生产环境 bug | `hotfix/fix-login-error` | 从 main 创建，合并回 main 和 develop |
| `release/` | 版本发布准备 | `release/v1.0.0` | 从 develop 创建，合并回 main 和 develop |
| `docs/` | 文档更新 | `docs/update-api-documentation` | 从 develop 创建，合并回 develop |

### 1.3 分支命名规范

```
<类型>/<简短描述>

# 示例
feature/person-search-enhance
bugfix/family-tree-render-error
hotfix/critical-security-patch
release/v1.2.0
docs/api-endpoints-documentation
```

**命名规则**：
- 使用小写字母和连字符 `-`
- 简短描述使用英文，不超过 50 字符
- 使用动词开头：`add-`, `fix-`, `update-`, `refactor-`
- 避免使用中文或特殊字符

---

## 二、提交规范 (Conventional Commits)

### 2.1 提交格式

```
<type>(<scope>): <subject>
<空行>
<body>
<空行>
<footer>
```

### 2.2 Type 类型

| 类型 | 说明 | 示例 |
|------|------|------|
| `feat` | 新功能 | `feat: add person search API` |
| `fix` | 修复 bug | `fix: resolve family tree render issue` |
| `docs` | 文档更新 | `docs: update API documentation` |
| `style` | 代码格式（不影响逻辑） | `style: format Go code with gofmt` |
| `refactor` | 重构（非新增功能也非修复bug） | `refactor: simplify person repository` |
| `perf` | 性能优化 | `perf: optimize ltree query performance` |
| `test` | 测试相关 | `test: add unit tests for person service` |
| `chore` | 构建/工具/依赖更新 | `chore: update gin to v1.9.1` |
| `ci` | CI/CD 配置变更 | `ci: add GitHub Actions workflow` |
| `revert` | 回滚提交 | `revert: revert previous commit` |

### 2.3 Scope 范围（可选）

指定影响的模块：

- `backend`: 后端整体
- `person`: 人物模块
- `auth`: 认证模块
- `frontend`: 前端整体
- `admin`: 管理后台
- `visualization`: 可视化前端
- `miniprogram`: 小程序
- `3d`: 3D 祠堂模块
- `docker`: Docker 相关
- `deps`: 依赖更新

**示例**：
```
feat(person): add batch import API
fix(visualization): fix sunburst chart render
docs(backend): update API docs
```

### 2.4 Subject 主题

- 使用祈使句，一般现在时："add" 而非 "added" 或 "adds"
- 首字母小写
- 结尾不使用句号
- 不超过 72 字符

**✅ 良好示例**：
```
feat(person): add fuzzy search by name
fix(tree): resolve parent-child relation bug
docs(api): update swagger documentation
```

**❌ 不良示例**：
```
fix bug                     # 太模糊，类型不明确
Added new feature           # 使用过去式
feat: add a really really long description that exceeds 72 characters... # 太长
```

### 2.5 Body 正文（可选）

详细描述变更的内容、原因和影响，每行不超过 72 字符。

**示例**：
```
feat(tree): add 5-generation family tree query

- Implement recursive query using PostgreSQL ltree
- Add depth limit parameter (up to 10 generations)
- Add caching for frequently accessed trees
- Benchmark shows 3x performance improvement

Closes #123
```

### 2.6 Footer 页脚（可选）

- **关联 Issue**：`Closes #123`, `Fixes #456`, `Resolves #789`
- **破坏性变更**：`BREAKING CHANGE: API endpoint changed`

---

## 三、开发工作流

### 3.1 功能开发流程

```
1. 拉取最新 develop
   └─ git checkout develop
   └─ git pull origin develop

2. 创建功能分支
   └─ git checkout -b feature/your-feature-name

3. 开发与提交
   └─ 编码...
   └─ git add .
   └─ git commit -m "feat: your commit message"

4. 同步最新代码（定期）
   └─ git checkout develop
   └─ git pull origin develop
   └─ git checkout feature/your-feature-name
   └─ git rebase develop

5. 推送到远程
   └─ git push origin feature/your-feature-name

6. 创建 Pull Request
   └─ 目标分支: develop
   └─ 填写 PR 模板
   └─ 请求 Code Review

7. 合并后删除分支
   └─ 本地: git branch -d feature/your-feature-name
   └─ 远程: git push origin --delete feature/your-feature-name
```

### 3.2 紧急修复流程 (Hotfix)

```
1. 从 main 创建 hotfix 分支
   └─ git checkout main
   └─ git pull origin main
   └─ git checkout -b hotfix/critical-fix

2. 修复并提交
   └─ git commit -m "fix: critical security patch"

3. 合并回 main 和 develop
   └─ 创建 PR 到 main
   └─ 合并后创建 PR 到 develop

4. 打标签发布
   └─ git tag v1.0.1
   └─ git push origin v1.0.1
```

### 3.3 版本发布流程

```
1. 从 develop 创建 release 分支
   └─ git checkout -b release/v1.0.0

2. 版本准备
   └─ 更新 CHANGELOG.md
   └─ 更新版本号
   └─ 最终测试和 bug 修复

3. 合并到 main 并打标签
   └─ PR 合并到 main
   └─ git tag v1.0.0
   └─ git push origin v1.0.0

4. 同步回 develop
   └─ PR 合并 release 到 develop
```

---

## 四、Pull Request 规范

### 4.1 PR 标题格式

```
<type>(<scope>): <subject>

# 示例
feat(person): add batch import functionality
fix(tree): resolve rendering performance issue
docs: update API documentation
```

### 4.2 PR 描述模板

```markdown
## 变更概述
<!-- 简要描述这个 PR 的主要变更 -->

## 变更类型
- [ ] ✨ 新功能 (feat)
- [ ] 🐛 Bug 修复 (fix)
- [ ] 📝 文档更新 (docs)
- [ ] 🎨 代码格式 (style)
- [ ] ♻️ 重构 (refactor)
- [ ] ⚡ 性能优化 (perf)
- [ ] ✅ 测试 (test)
- [ ] 🔧 构建/工具 (chore)

## 技术细节
<!-- 详细描述实现细节、设计决策等 -->

## 测试情况
- [ ] 单元测试已通过
- [ ] 集成测试已通过
- [ ] 手动测试已完成

## 关联 Issue
<!-- Closes #123, Fixes #456 -->

## 截图（如适用）
<!-- 添加相关截图或演示 -->
```

### 4.3 Code Review 清单

**审查者需要检查**：
- [ ] 代码逻辑正确，没有明显的 bug
- [ ] 符合项目的编码规范
- [ ] 命名清晰、语义明确
- [ ] 没有代码异味（Code Smell）
- [ ] 测试覆盖率足够
- [ ] 文档已更新
- [ ] 没有安全漏洞
- [ ] 性能影响已评估

**审查结论**：
- ✅ **Approve**：同意合并
- 🟡 **Comment**：提出意见，需要改进
- ❌ **Request Changes**：需要重大修改

---

## 五、Git 常用命令速查

### 5.1 基础操作

```bash
# 克隆仓库
git clone https://github.com/your-org/genealogy-ma.git

# 查看状态
git status

# 查看变更
git diff

# 添加文件
git add .                      # 所有文件
git add path/to/file.go       # 单个文件
git add -p                     # 交互式添加

# 提交
git commit -m "feat: your message"

# 推送
git push origin branch-name

# 拉取
git pull origin branch-name
```

### 5.2 分支操作

```bash
# 创建并切换分支
git checkout -b feature/new-feature

# 切换分支
git checkout develop

# 查看分支
git branch -a

# 删除本地分支
git branch -d branch-name

# 删除远程分支
git push origin --delete branch-name

# 重命名分支
git branch -m old-name new-name
```

### 5.3 变基与合并

```bash
# 变基（保持提交历史线性）
git rebase develop

# 解决冲突后继续
git rebase --continue

# 中止变基
git rebase --abort

# 合并分支
git merge feature/branch-name

# 压缩合并
git merge --squash feature/branch-name
```

### 5.4 查看历史

```bash
# 查看提交历史
git log

# 简洁的图形化历史
git log --oneline --graph --all

# 查看某文件变更历史
git log --oneline path/to/file.go

# 查看具体变更
git show commit-hash
```

### 5.5 撤销操作

```bash
# 撤销工作区变更
git checkout -- file.go

# 撤销暂存
git reset HEAD file.go

# 撤销最近一次提交（保留变更）
git reset --soft HEAD~1

# 撤销最近一次提交（丢弃变更）
git reset --hard HEAD~1

# 回滚到指定提交
git revert commit-hash
```

---

## 六、Git 最佳实践

### 6.1 提交最佳实践

✅ **Do**：
- 小步提交，每个提交一个独立变更
- 提交前自行 review 代码
- 确保代码可以编译和运行
- 写有意义的提交信息
- 定期同步远程分支

❌ **Don't**：
- 不要提交不完整的代码
- 不要提交密码、密钥等敏感信息
- 不要提交 node_modules、dist 等构建产物
- 不要直接在 main/develop 分支提交
- 不要提交 large binary 文件（使用 LFS 或 CDN）

### 6.2 分支最佳实践

- 分支生命周期越短越好
- 定期同步 develop 分支
- 合并后立即删除分支
- 一个分支只做一件事
- 避免在分支上积累大量变更

### 6.3 团队协作

- 每天至少 pull 一次 develop
- 功能开发前先沟通
- PR 尽量小而聚焦
- Review 响应不超过 24 小时
- 使用 Draft PR 进行早期反馈

---

## 七、版本号规范

遵循 [Semantic Versioning 2.0](https://semver.org/lang/zh-CN/)

```
主版本号.次版本号.修订号
  │        │        │
  │        │        └─ 向后兼容的 bug 修复
  │        └─ 向后兼容的新功能
  └─ 不向后兼容的破坏性变更
```

**示例**：
- `v1.0.0` - 初始版本发布
- `v1.0.1` - 修复 bug
- `v1.1.0` - 新增功能
- `v2.0.0` - 破坏性变更

---

## 八、工具推荐

### 8.1 Git 客户端

- **GitLens** - VS Code 插件，增强 Git 功能
- **SourceTree** - 图形化 Git 客户端
- **GitHub Desktop** - GitHub 官方客户端

### 8.2 提交规范工具

- **commitlint** - 检查提交信息是否符合规范
- **husky** - Git hooks 管理
- **cz-git** - 交互式提交工具

---

## 参考资料

- [Conventional Commits](https://www.conventionalcommits.org/)
- [Git Pro 中文版](https://git-scm.com/book/zh/v2)
- [GitHub Flow](https://guides.github.com/introduction/flow/)
