# 移动端项目

族谱数字化管理平台 - 移动端项目

## 项目结构

```
mobile/
└── miniprogram/          # 微信小程序
```

---

## 📱 miniprogram - 微信小程序

**技术栈**: uni-app + Vue 3 + TypeScript + uView UI

**目标平台**:
- ✅ 微信小程序
- ✅ H5 (可选)
- ✅ App (可选)

## 主要功能

### 1. 族谱阅览
- 人物搜索（姓名、字号）
- 人物详情卡片
- 族谱树可视化
- 上下五代世系图

### 2. 3D 祠堂
- 简化版 3D 祠堂漫游
- 牌位查询与供奉
- 虚拟祭祀功能

### 3. 家族文化
- 家族故事浏览
- 文献库查阅
- 家训家规展示

### 4. 交流社区
- 家族动态流
- 寻根问祖板块
- 族人通讯录（授权可见）
- 私信聊天

### 5. 个人中心
- 用户信息管理
- 实名认证
- 收藏夹
- 浏览历史
- 设置与通知

## 开发准备

### 环境要求

- Node.js 18+
- HBuilderX 或 VS Code
- 微信开发者工具
- 微信小程序 AppID

### 开发命令

```bash
cd mobile/miniprogram

# 安装依赖
npm install

# H5 开发模式
npm run dev:h5

# 微信小程序开发模式
npm run dev:mp-weixin

# 构建发布 - 微信小程序
npm run build:mp-weixin
```

## 目录结构

```
miniprogram/
├── src/
│   ├── pages/           # 页面
│   │   ├── home/       # 首页
│   │   ├── person/     # 人物相关
│   │   ├── tree/       # 族谱树
│   │   ├── hall/       # 3D 祠堂
│   │   ├── culture/    # 家族文化
│   │   ├── community/  # 社区交流
│   │   └── mine/       # 个人中心
│   ├── components/     # 公共组件
│   ├── api/            # API 接口封装
│   ├── store/          # 状态管理 (Pinia)
│   ├── utils/          # 工具函数
│   ├── styles/         # 全局样式
│   ├── static/         # 静态资源
│   └── App.vue
├── manifest.json        # uni-app 配置
├── pages.json           # 页面路由配置
└── package.json
```

## 性能优化

### 小程序端注意事项

1. **包体积控制**
   - 主包控制在 2MB 以内
   - 采用分包加载策略
   - 图片资源使用 CDN

2. **渲染性能**
   - 列表使用虚拟滚动
   - 图片懒加载
   - 减少 setData 调用

3. **3D 优化**
   - 小程序端简化模型面数
   - 使用压缩纹理
   - 考虑降级为 2D 方案

## 发布流程

1. 执行 `npm run build:mp-weixin`
2. 在微信开发者工具中导入 `dist/build/mp-weixin`
3. 上传代码
4. 提交审核
5. 发布上线

## 注意事项

### 小程序审核相关

- ❌ 不能有虚拟支付
- ❌ 不能诱导分享
- ✅ 内容合规（家族文化内容需正向）
- ✅ 用户隐私协议

### 后端对接

API 域名需要在微信公众平台配置 request 合法域名：

```
https://api.your-genealogy-domain.com
```

WebSocket 域名也需要单独配置。
