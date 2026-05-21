# 族谱数字化管理平台 - 微信小程序

## 项目简介

基于 uni-app + Vue 3 + TypeScript 的族谱微信小程序，提供人物浏览、族谱树查看、个人中心等功能。

## 技术栈

- **框架**: uni-app 3.x (Vue 3)
- **语言**: TypeScript 5.x
- **构建工具**: Vite 5.x
- **状态管理**: Vue Composition API (reactive)

## 目录结构

```
miniprogram/
├── src/
│   ├── api/              # API 请求封装
│   │   ├── request.ts    # 通用请求方法
│   │   ├── auth.ts       # 认证相关 API
│   │   └── person.ts     # 人物相关 API
│   ├── components/        # 公共组件
│   ├── pages/            # 页面
│   │   ├── login/        # 登录页
│   │   ├── person/       # 人物相关页
│   │   ├── genealogy/     # 族谱页
│   │   └── user/         # 个人中心
│   ├── store/            # 状态管理
│   │   └── user.ts       # 用户状态
│   ├── types/            # 类型定义
│   │   └── person.ts     # 人物类型
│   ├── App.vue           # 应用入口
│   ├── main.ts           # 主入口
│   └── pages.json        # 页面配置
├── manifest.json        # uni-app 配置
├── package.json
├── tsconfig.json
└── vite.config.ts
```

## 页面说明

| 页面 | 路径 | 说明 |
|------|------|------|
| 登录 | `/pages/login/index` | 用户登录、微信授权 |
| 人物列表 | `/pages/person/list` | 人物搜索、筛选、列表 |
| 人物详情 | `/pages/person/detail` | 人物信息、亲属关系 |
| 族谱树 | `/pages/genealogy/tree` | 族谱树可视化 |
| 个人中心 | `/pages/user/profile` | 用户信息、设置 |

## 开发命令

```bash
# 安装依赖
npm install

# 开发调试（微信小程序）
npm run dev

# 构建发布
npm run build
```

## API 文档

详见 [backend API documentation](http://localhost:8080/swagger/index.html)

## 注意事项

1. 微信小程序需要配置合法的 AppID
2. 接口请求需要配合后端 JWT 认证
3. 生产环境请修改 `src/api/request.ts` 中的 `API_BASE_URL`