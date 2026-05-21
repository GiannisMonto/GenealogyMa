# 族谱数字化管理平台 - 开发计划

> **最后更新**: 2026-05-21
> **版本**: v1.0

---

## 📋 项目概览

### 项目背景
为姓氏宗族提供完整的数字化族谱管理解决方案，包含人物数据管理、世系关系可视化、家族文化传承、在线宗祠等核心功能。

### 当前状态
- **后端架构**: ✅ Go + Gin + DDD 四层架构完成
- **核心模块**: ✅ Person 人物领域 CRUD + 族谱树完成
- **权限系统**: ✅ RBAC 后端实现 + 数据库迁移脚本完成
- **管理后台**: ⏳ 待开发（React + Ant Design）
- **可视化前端**: ⏳ 待开发（Vue 3 + D3.js）
- **微信小程序**: ⏳ 待开发（uni-app）

---

## 🎯 开发阶段规划

### 阶段一：权限系统完整实现 ✅ (89%)
**目标**: 完成用户认证与权限管理核心功能

| 任务 | 关键文件路径 | 状态 | 完成日期 | 备注 |
|------|-------------|------|---------|------|
| 1.1 用户领域实体创建 | `backend/internal/domain/user/entity.go` | ✅ | 2026-05-21 | User 聚合根、角色关联 |
| 1.2 用户仓储接口与实现 | `backend/internal/domain/user/repository.go` | ✅ | 2026-05-21 | 接口定义 + PostgreSQL实现 |
| 1.3 用户应用服务 | `backend/internal/application/service/user_service.go` | ✅ | 2026-05-21 | 注册、登录、角色管理 |
| 1.4 认证控制器 | `backend/internal/interfaces/http/controller/auth_controller.go` | ✅ | 2026-05-21 | 登录、注册、刷新令牌、登出 |
| 1.5 数据库迁移脚本 | `backend/internal/infrastructure/persistence/migrations/` | ✅ | 2026-05-21 | users、roles、permissions、关联表 |
| 1.6 密码加密服务 | `backend/pkg/auth/password.go` | ✅ | 2026-05-21 | bcrypt 实现 |
| 1.7 令牌黑名单 | `backend/internal/infrastructure/cache/token_blacklist.go` | ✅ | 2026-05-21 | Redis 实现令牌注销 |
| 1.8 初始化数据脚本 | `backend/internal/infrastructure/persistence/migrations/006_init_data.sql` | ✅ | 2026-05-21 | 默认角色、权限、超级管理员 |
| 1.9 权限管理 API | `backend/internal/interfaces/http/controller/permission_controller.go` | ⏳ | | 角色、权限 CRUD 接口完善 |

---

### 阶段二：管理后台基础架构 ⏳ (0%)
**目标**: 从零搭建 React + Ant Design 5 管理后台框架

| 任务 | 关键文件路径 | 状态 | 完成日期 | 负责人 | 备注 |
|------|-------------|------|---------|--------|------|
| 2.1 管理后台项目初始化 | `frontend/admin/package.json` | ⏳ | | Claude Code | Vite 6 + React 18 + TypeScript 5 |
| 2.2 项目基础配置 | `frontend/admin/vite.config.ts`, `frontend/admin/tsconfig.json` | ⏳ | | | 路径别名 @、ESLint、Prettier、构建优化 |
| 2.3 AntD 主题配置 | `frontend/admin/src/theme/index.ts` | ⏳ | | | 定制族谱主题色、全局样式、CSS 变量 |
| 2.4 路由系统搭建 | `frontend/admin/src/router/index.tsx` | ⏳ | | | React Router v6 + 权限路由守卫、动态菜单 |
| 2.5 状态管理集成 | `frontend/admin/src/store/` | ⏳ | | | Zustand - 用户状态、权限状态、全局配置 |
| 2.6 API 请求封装 | `frontend/admin/src/api/request.ts` | ⏳ | | | Axios 封装、请求/响应拦截器、错误处理、类型定义 |
| 2.7 全局布局组件 | `frontend/admin/src/layouts/BasicLayout.tsx` | ⏳ | | | 侧边栏、顶部导航、面包屑、标签页、页脚 |
| 2.8 登录页面开发 | `frontend/admin/src/pages/login/index.tsx` | ⏳ | | | 登录表单、验证码、忘记密码、对接后端 API |
| 2.9 通用组件库 | `frontend/admin/src/components/` | ⏳ | | | 可复用组件：PageHeader、StatusTag、SearchForm 等 |

---

### 阶段三：数据查看与可视化 ⏳ (0%)
**目标**: 实现管理后台的数据查看与族谱可视化功能

| 任务 | 关键文件路径 | 状态 | 完成日期 | 负责人 | 备注 |
|------|-------------|------|---------|--------|------|
| 3.1 数据看板首页 | `frontend/admin/src/pages/dashboard/index.tsx` | ⏳ | | | ECharts 5 统计图表、关键指标卡片、数据概览 |
| 3.2 人物列表页面 | `frontend/admin/src/pages/person/list.tsx` | ⏳ | | | AntD Table、高级搜索、多维度筛选、分页、操作列 |
| 3.3 人物详情页面 | `frontend/admin/src/pages/person/detail.tsx` | ⏳ | | | 详情展示、亲属关系卡片、生平信息、时间线 |
| 3.4 族谱树可视化组件 | `frontend/admin/src/components/GenealogyTree/` | ⏳ | | | D3.js v7 可交互族谱树、支持缩放拖拽、节点高亮 |
| 3.5 世代分布可视化 | `frontend/admin/src/pages/statistics/generation.tsx` | ⏳ | | | 旭日图、柱状图、饼图等多种展示形式 |
| 3.6 人物 API 封装 | `frontend/admin/src/api/person.ts` | ⏳ | | | 人物 CRUD、搜索、族谱树查询、类型定义 |
| 3.7 统计 API 封装 | `frontend/admin/src/api/statistics.ts` | ⏳ | | | 统计数据接口封装、图表数据适配 |
| 3.8 高级搜索组件 | `frontend/admin/src/components/AdvancedSearch/` | ⏳ | | | 多条件组合搜索、保存搜索条件 |
| 3.9 族谱树页面 | `frontend/admin/src/pages/genealogy/tree.tsx` | ⏳ | | | 完整族谱树视图、快速定位、导出图片 |

---

### 阶段四：数据勘误与补充功能 ⏳ (0%)
**目标**: 实现数据管理的高级功能 - 批量操作、导入导出、审计追踪

| 任务 | 关键文件路径 | 状态 | 完成日期 | 负责人 | 备注 |
|------|-------------|------|---------|--------|------|
| 4.1 人物表单编辑页面 | `frontend/admin/src/pages/person/form.tsx` | ⏳ | | | 创建/编辑人物、亲属选择器、表单验证 |
| 4.2 批量操作 API 完善 | `backend/internal/interfaces/http/controller/person_controller.go` | ⏳ | | | 人物批量创建、更新、删除接口 |
| 4.3 批量操作前端页面 | `frontend/admin/src/pages/person/batch.tsx` | ⏳ | | | 批量选择、批量编辑弹窗、批量删除确认 |
| 4.4 Excel 导入功能 - 后端 | `backend/internal/application/service/import_service.go` | ⏳ | | | Excel/CSV 解析、数据验证、冲突处理、导入队列 |
| 4.5 数据导入前端页面 | `frontend/admin/src/pages/data/import.tsx` | ⏳ | | | 文件上传、数据预览、冲突处理、导入进度条 |
| 4.6 数据导出功能 - 后端 | `backend/internal/application/service/export_service.go` | ⏳ | | | Excel/CSV/PDF 导出、按条件筛选导出 |
| 4.7 数据导出前端界面 | `frontend/admin/src/pages/data/export.tsx` | ⏳ | | | 导出条件配置、字段选择、导出进度、下载管理 |
| 4.8 审计日志领域模块 | `backend/internal/domain/audit/` | ⏳ | | | 审计实体、仓储、服务、GORM 中间件 |
| 4.9 变更历史页面 | `frontend/admin/src/pages/data/history.tsx` | ⏳ | | | 变更记录列表、详情查看、版本对比、回滚操作 |
| 4.10 数据质量检测服务 | `backend/internal/application/service/data_quality_service.go` | ⏳ | | | 重复数据检测、完整性检查、关系异常检测 |
| 4.11 数据质量报告页面 | `frontend/admin/src/pages/data/quality.tsx` | ⏳ | | | 质量报告展示、问题列表、修复引导 |

---

### 阶段五：权限管理界面 ⏳ (0%)
**目标**: 实现完整的 RBAC 权限管理后台界面

| 任务 | 关键文件路径 | 状态 | 完成日期 | 负责人 | 备注 |
|------|-------------|------|---------|--------|------|
| 5.1 用户管理列表页面 | `frontend/admin/src/pages/user/list.tsx` | ⏳ | | | 用户列表、搜索筛选、状态管理、重置密码 |
| 5.2 用户编辑表单页面 | `frontend/admin/src/pages/user/form.tsx` | ⏳ | | | 创建/编辑用户、角色分配、资料编辑 |
| 5.3 角色管理列表页面 | `frontend/admin/src/pages/role/list.tsx` | ⏳ | | | 角色列表、角色状态管理、系统角色保护 |
| 5.4 角色权限配置页面 | `frontend/admin/src/pages/role/permission.tsx` | ⏳ | | | 树形权限分配界面、权限继承、批量配置 |
| 5.5 操作日志页面 | `frontend/admin/src/pages/audit/logs.tsx` | ⏳ | | | 审计日志查询、筛选、导出、详情查看 |
| 5.6 权限管理 API 完善 | `backend/internal/interfaces/http/controller/permission_controller.go` | ⏳ | | | 角色、权限的完整 CRUD 接口 |
| 5.7 权限 API 封装 | `frontend/admin/src/api/permission.ts` | ⏳ | | | 用户、角色、权限 API 类型定义 |

---

### 阶段六：其他管理后台功能 ⏳ (0%)
**目标**: 完成系统管理、内容审核、配置等辅助功能

| 任务 | 关键文件路径 | 状态 | 完成日期 | 负责人 | 备注 |
|------|-------------|------|---------|--------|------|
| 6.1 系统配置页面 | `frontend/admin/src/pages/settings/` | ⏳ | | | 系统参数配置、网站设置、功能开关 |
| 6.2 配置管理 API | `backend/internal/interfaces/http/controller/config_controller.go` | ⏳ | | | 配置 CRUD 接口、配置缓存 |
| 6.3 内容审核页面 | `frontend/admin/src/pages/review/` | ⏳ | | | 待审核数据列表、审核操作、审核记录 |
| 6.4 系统日志页面 | `frontend/admin/src/pages/system/logs.tsx` | ⏳ | | | 系统日志查询、筛选、导出、实时查看 |
| 6.5 数据库备份/恢复服务 | `backend/internal/application/service/backup_service.go` | ⏳ | | | 自动备份、手动备份、备份列表、恢复功能 |
| 6.6 备份管理页面 | `frontend/admin/src/pages/system/backup.tsx` | ⏳ | | | 备份列表、手动备份、恢复操作、下载备份 |
| 6.7 文件管理功能 | `frontend/admin/src/pages/filemanager/` | ⏳ | | | 文件上传、浏览、管理、存储配置 |

---

### 阶段七：可视化前端项目 ⏳ (0%)
**目标**: 搭建独立的 Vue 3 + D3.js 族谱可视化项目

| 任务 | 关键文件路径 | 状态 | 完成日期 | 负责人 | 备注 |
|------|-------------|------|---------|--------|------|
| 7.1 可视化项目初始化 | `frontend/visualization/package.json` | ⏳ | | | Vue 3 + Vite + TypeScript + D3.js |
| 7.2 族谱树可视化组件 | `frontend/visualization/src/components/GenealogyTree/` | ⏳ | | | D3.js 高性能族谱树、支持 10k+ 节点 |
| 7.3 人物详情弹窗 | `frontend/visualization/src/components/PersonDetail/` | ⏳ | | | 人物信息展示、关系跳转 |
| 7.4 搜索定位功能 | `frontend/visualization/src/components/SearchPanel/` | ⏳ | | | 搜索人物、快速定位到节点 |
| 7.5 世代旭日图 | `frontend/visualization/src/components/SunburstChart/` | ⏳ | | | 世代分布可视化、点击下钻 |
| 7.6 迁徙地图组件 | `frontend/visualization/src/components/MigrationMap/` | ⏳ | | | 基于地图的族人迁徙路径展示 |
| 7.7 3D 宗祠模块 | `frontend/visualization/src/components/AncestralHall/` | ⏳ | | | Three.js 3D 宗祠场景、第一人称漫游 |

---

### 阶段八：微信小程序 ⏳ (0%)
**目标**: 基于 uni-app 的族谱微信小程序

| 任务 | 关键文件路径 | 状态 | 完成日期 | 负责人 | 备注 |
|------|-------------|------|---------|--------|------|
| 8.1 小程序项目初始化 | `mobile/miniprogram/package.json` | ✅ | 2026-05-21 | Claude Code | uni-app + Vue 3 + TypeScript |
| 8.2 登录/认证模块 | `mobile/miniprogram/src/pages/login/` | ✅ | 2026-05-21 | Claude Code | 登录页、微信授权、JWT 认证 |
| 8.3 人物列表页 | `mobile/miniprogram/src/pages/person/list.vue` | ✅ | 2026-05-21 | Claude Code | 移动端人物列表、搜索、筛选 |
| 8.4 人物详情页 | `mobile/miniprogram/src/pages/person/detail.vue` | ✅ | 2026-05-21 | Claude Code | 人物信息、亲属关系展示 |
| 8.5 族谱树页面 | `mobile/miniprogram/src/pages/genealogy/tree.vue` | ✅ | 2026-05-21 | Claude Code | 移动端族谱树可视化 |
| 8.6 个人中心 | `mobile/miniprogram/src/pages/user/profile.vue` | ✅ | 2026-05-21 | Claude Code | 个人信息、设置、退出登录 |

---

## 📊 总体进度仪表盘

| 阶段 | 总任务数 | 已完成 | 进行中 | 待开始 | 完成度 |
|------|---------|--------|--------|--------|--------|
| 阶段一：权限系统 | 9 | 8 | 0 | 1 | 89% |
| 阶段二：管理后台基础 | 9 | 0 | 0 | 9 | 0% |
| 阶段三：数据查看与可视化 | 9 | 0 | 0 | 9 | 0% |
| 阶段四：数据勘误与补充 | 11 | 0 | 0 | 11 | 0% |
| 阶段五：权限管理界面 | 7 | 0 | 0 | 7 | 0% |
| 阶段六：其他管理功能 | 7 | 0 | 0 | 7 | 0% |
| 阶段七：可视化前端项目 | 7 | 0 | 0 | 7 | 0% |
| 阶段八：微信小程序 | 6 | 6 | 0 | 0 | 100% |
| **总计** | **65** | **14** | **0** | **51** | **22%** |

---

## 📝 更新记录

| 日期 | 更新内容 | 更新人 |
|------|---------|--------|
| 2026-05-21 | 创建初始开发阶段计划 | Claude Code |
| 2026-05-21 | 完成阶段一权限系统后端实现（User实体、仓储、服务、控制器、令牌黑名单） | Claude Code |
| 2026-05-21 | 完成数据库迁移脚本与初始化数据，更新计划至8个阶段，新增可视化前端与小程序开发计划 | Claude Code |

---

## 📌 使用说明

1. **每次功能开发前**：检查本计划中的任务状态，确认任务归属
2. **功能完成后**：更新对应任务的「状态」和「完成日期」
3. **新增功能时**：在对应阶段添加任务行，并更新进度仪表盘
4. **里程碑完成时**：在「更新记录」中添加里程碑记录
