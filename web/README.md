# 短剧管理系统前端

基于 Vue 3 + Vite + Element Plus 的现代化管理后台。

## 技术栈

- **Vue 3** - 渐进式 JavaScript 框架
- **Vite** - 下一代前端构建工具
- **Element Plus** - 基于 Vue 3 的组件库
- **Vue Router** - Vue.js 官方路由管理器
- **Pinia** - Vue 的状态管理库
- **Axios** - HTTP 客户端

## 功能特性

- 🎯 现代化的管理界面
- 📱 响应式设计，支持移动端
- 🔐 完整的认证系统
- 📊 数据可视化
- 🎨 优雅的 UI 设计
- ⚡ 快速的开发体验

## 快速开始

### 安装依赖

```bash
npm install
# 或
yarn install
# 或
pnpm install
```

### 开发环境

```bash
npm run dev
```

访问 http://localhost:3000

### 构建生产版本

```bash
npm run build
```

### 预览生产版本

```bash
npm run preview
```

## 项目结构

```
frontend/
├── public/                 # 静态资源
├── src/
│   ├── components/         # 公共组件
│   ├── layouts/           # 布局组件
│   ├── views/             # 页面组件
│   ├── router/            # 路由配置
│   ├── stores/            # 状态管理
│   ├── utils/             # 工具函数
│   ├── App.vue            # 根组件
│   └── main.js            # 入口文件
├── index.html             # HTML 模板
├── vite.config.js         # Vite 配置
└── package.json           # 项目配置
```

## 开发说明

### API 接口

前端通过代理访问后端 API：
- 开发环境：http://localhost:1800
- 生产环境：根据部署配置

### 路由结构

- `/login` - 登录页面
- `/dashboard` - 仪表板
- `/dramas` - 短剧管理
- `/episodes` - 剧集管理
- `/users` - 用户管理

### 状态管理

使用 Pinia 进行状态管理：
- `auth` - 认证状态
- 其他业务状态可按需添加

## 部署

### 构建

```bash
npm run build
```

构建产物将输出到 `../web/dist` 目录，可直接被 Go 后端服务。

### 环境变量

- `VITE_API_BASE_URL` - API 基础地址
- `VITE_APP_TITLE` - 应用标题

## 开发规范

- 使用 Composition API
- 组件名使用 PascalCase
- 文件名使用 kebab-case
- 遵循 Vue 3 最佳实践
