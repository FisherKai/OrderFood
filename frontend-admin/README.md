# 点餐系统 - PC端管理后台

基于 Vue 3 + Element Plus 的现代化管理后台系统。

## ✨ 功能特性

### 🏠 仪表盘
- 数据统计概览
- 图表展示（订单趋势、分类统计）
- 最近订单列表
- 系统信息

### 🍽️ 菜品管理
- 菜品增删改查
- 图片上传和管理
- 批量操作
- 状态管理（上架/下架）
- 库存管理

### 📂 分类管理
- 分类增删改查
- 图标选择器
- 排序管理
- 菜品数量统计

### 📋 订单管理
- 订单列表查看
- 订单状态更新
- 订单详情展开
- 批量处理
- 预约订单管理

### ⭐ 评价管理
- 评价列表查看
- 图片预览
- 评价删除
- 批量操作

### 📢 公告管理
- 公告增删改查
- 富文本编辑
- 时间控制
- 优先级设置
- 公告预览

### 👥 用户管理
- 用户列表查看
- 用户状态管理
- 用户详情查看
- 订单统计
- 批量操作

### 🖼️ 图片管理
- 图片上传
- 网格/列表视图
- 图片预览
- 软删除和恢复
- 垃圾清理
- 使用情况统计

## 🚀 技术栈

- **框架**: Vue 3 (Composition API)
- **UI库**: Element Plus
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP客户端**: Axios
- **图表**: ECharts + Vue-ECharts
- **构建工具**: Vite
- **CSS预处理器**: SCSS

## 📦 安装依赖

```bash
npm install
```

## 🛠️ 开发运行

```bash
npm run dev
```

访问地址：http://localhost:3001

## 🏗️ 构建部署

```bash
npm run build
```

## 🔐 默认账户

- 用户名：`admin`
- 密码：`admin123`

## 📁 项目结构

```
frontend-admin/
├── public/                 # 静态资源
├── src/
│   ├── api/               # API接口
│   │   ├── request.js     # axios配置
│   │   ├── auth.js        # 认证接口
│   │   ├── dish.js        # 菜品接口
│   │   ├── category.js    # 分类接口
│   │   ├── order.js       # 订单接口
│   │   ├── review.js      # 评价接口
│   │   ├── announcement.js # 公告接口
│   │   └── upload.js      # 上传接口
│   ├── layout/            # 布局组件
│   │   └── Layout.vue     # 主布局
│   ├── router/            # 路由配置
│   │   └── index.js       # 路由文件
│   ├── stores/            # 状态管理
│   │   └── admin.js       # 管理员状态
│   ├── styles/            # 样式文件
│   │   └── main.scss      # 全局样式
│   ├── views/             # 页面组件
│   │   ├── Login.vue      # 登录页
│   │   ├── Dashboard.vue  # 仪表盘
│   │   ├── Dishes.vue     # 菜品管理
│   │   ├── Categories.vue # 分类管理
│   │   ├── Orders.vue     # 订单管理
│   │   ├── Reviews.vue    # 评价管理
│   │   ├── Announcements.vue # 公告管理
│   │   ├── Users.vue      # 用户管理
│   │   └── Images.vue     # 图片管理
│   ├── App.vue            # 根组件
│   └── main.js            # 入口文件
├── index.html             # HTML模板
├── package.json           # 依赖配置
├── vite.config.js         # Vite配置
└── README.md              # 说明文档
```

## 🎨 界面预览

### 登录页面
- 现代化登录界面
- 表单验证
- 响应式设计

### 仪表盘
- 数据统计卡片
- 可视化图表
- 最近活动列表

### 管理页面
- 统一的搜索栏
- 表格展示
- 批量操作
- 分页功能

## 📝 开发规范

### 代码风格
- 使用 Composition API
- 组件采用 `<script setup>` 语法
- 样式使用 SCSS
- 遵循 Vue 3 最佳实践

### 命名规范
- 组件名使用 PascalCase
- 文件名使用 PascalCase
- 变量使用 camelCase
- 常量使用 UPPER_CASE

### 目录结构
- 按功能模块组织代码
- API接口按业务分类
- 组件复用性考虑

## 🔧 配置说明

### 环境配置
- 开发环境：`npm run dev`
- 生产构建：`npm run build`
- 预览构建：`npm run preview`

### 代理配置
开发环境下，API请求会代理到 `http://localhost:8080`

### Element Plus
- 使用中文语言包
- 按需引入组件
- 支持暗黑主题

## 🚨 注意事项

1. **权限控制**: 所有页面都需要登录后访问
2. **数据模拟**: 当前使用模拟数据，实际使用需要对接真实API
3. **图片上传**: 需要配置正确的上传接口地址
4. **浏览器兼容**: 支持现代浏览器，IE不支持

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支
3. 提交代码
4. 发起 Pull Request

## 📄 许可证

MIT License