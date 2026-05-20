# 前端项目

族谱数字化管理平台 - 前端项目集合

## 项目结构

```
frontend/
├── admin/              # 管理后台 (React + Ant Design)
├── visualization/      # 可视化前端 (Vue 3 + D3.js)
└── shared/            # 共享组件 / TypeScript 类型
```

## 项目说明

### 🎛️ admin - 管理后台
**技术栈**: React 18 + TypeScript + Ant Design 5 + Vite

**主要功能**:
- 人物数据管理 (CRUD + 批量导入)
- 族谱关系编辑
- 用户权限管理
- 内容审核 (家族故事、文献)
- 系统配置
- 数据统计看板

**开发命令**:
```bash
cd frontend/admin
npm install
npm run dev      # 启动开发服务器
npm run build    # 生产构建
npm run lint     # 代码检查
```

---

### 📊 visualization - 可视化前端
**技术栈**: Vue 3 + TypeScript + Element Plus + D3.js + ECharts + Three.js

**主要功能**:
- 族谱树可视化展示 (D3.js)
- 旭日图、桑基图等高级可视化
- 迁徙地图
- 3D 祠堂漫游 (Three.js)
- 家族文化展览馆

**开发命令**:
```bash
cd frontend/visualization
npm install
npm run dev      # 启动开发服务器
npm run build    # 生产构建
```

---

### 📦 shared - 共享资源

**内容**:
- TypeScript 类型定义
- 公共工具函数
- API 客户端封装
- 通用 UI 组件
- 样式变量

## 开发规范

### 通用规范
- 使用 TypeScript，避免 `any`
- ESLint + Prettier 格式化
- 组件采用 PascalCase 命名
- 使用函数式组件 + Hooks

### Git 提交
遵循项目根目录的 [Git 工作流规范](../docs/GIT_WORKFLOW.md)

```bash
# 前端相关提交示例
feat(admin): add person batch import
fix(visualization): fix tree render error
style: format code with prettier
```
