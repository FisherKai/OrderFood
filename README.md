# 点餐系统 (OrderFood System)

基于 Gin + Vue3 的前后端分离点餐系统，支持移动端和PC端。

## 📋 项目概述

### 技术栈

**后端**
- Go 1.21+
- Gin Web Framework
- GORM (MySQL)
- JWT 认证
- Viper 配置管理

**前端**
- Vue 3 + Vite
- Pinia 状态管理
- Vue Router
- Vant 4 (移动端)
- Element Plus (PC端管理后台)
- Axios
- Swiper

### 主要功能

#### 用户端功能
- ✅ 用户注册/登录
- ✅ 菜品浏览（分类筛选、搜索）
- ✅ 购物车管理
- ✅ 在线点餐
- ✅ 预约点餐
- ✅ 订单查看
- ✅ 菜品评价
- ✅ 公告轮播

#### 管理端功能  
- ✅ 管理员登录
- ✅ 数据统计仪表盘
- ✅ 菜品管理（增删改查）
- ✅ 菜品图片管理（上传、删除、历史图片恢复）
- ✅ 分类管理
- ✅ 订单管理
- ✅ 预约订单管理
- ✅ 评价管理
- ✅ 公告管理
- ✅ 用户管理
- ✅ 图片存储管理

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Node.js 16+
- MySQL 5.7+

### 1. 克隆项目

```bash
cd /Users/shiyijiang/Documents/code/OrderFood
```

### 2. 启动后端

```bash
cd backend

# 安装依赖
go mod download

# 创建数据库
mysql -u root -p
CREATE DATABASE orderfood_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
exit;

# 修改配置文件
# 编辑 config.yaml，设置数据库连接信息

# 运行服务
go run main.go
```

后端服务运行在: `http://localhost:8080`

### 3. 启动移动端前端

```bash
cd frontend-mobile

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

移动端运行在: `http://localhost:3000`

### 4. 启动PC端管理后台

```bash
cd frontend-admin

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

管理后台运行在: `http://localhost:3001`

### 5. 创建管理员账户

首次运行后，需要创建管理员账户。可以直接插入SQL或使用Go代码加密密码：

```go
// 生成密码哈希的示例代码
package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    password := "admin123"
    hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    fmt.Println(string(hash))
}
```

然后插入数据库：

```sql
INSERT INTO admins (username, password, role, created_at, updated_at) 
VALUES ('admin', '生成的密码哈希', 'admin', NOW(), NOW());
```

## 📁 项目结构

```
OrderFood/
├── backend/                 # 后端项目
│   ├── config/             # 配置管理
│   ├── controllers/        # 控制器
│   ├── database/           # 数据库
│   ├── middleware/         # 中间件
│   ├── models/             # 数据模型
│   ├── router/             # 路由
│   ├── utils/              # 工具函数
│   ├── uploads/            # 上传文件
│   ├── config.yaml         # 配置文件
│   ├── main.go             # 入口文件
│   └── go.mod              # 依赖管理
│
├── frontend-mobile/        # 移动端前端
│   ├── src/
│   │   ├── api/           # API接口
│   │   ├── components/    # 组件
│   │   ├── router/        # 路由
│   │   ├── stores/        # 状态管理
│   │   ├── styles/        # 样式
│   │   ├── views/         # 页面
│   │   ├── App.vue        # 根组件
│   │   └── main.js        # 入口文件
│   ├── package.json
│   └── vite.config.js
│
├── frontend-admin/         # PC端管理后台
│   ├── src/
│   │   ├── api/           # API接口
│   │   ├── layout/        # 布局组件
│   │   ├── router/        # 路由
│   │   ├── stores/        # 状态管理
│   │   ├── styles/        # 样式
│   │   ├── views/         # 页面
│   │   ├── App.vue        # 根组件
│   │   └── main.js        # 入口文件
│   ├── package.json
│   └── vite.config.js
│
└── 需求.md                 # 需求文档
```

## 🗄️ 数据库设计

### 核心表
- `users` - 用户表
- `admins` - 管理员表
- `categories` - 菜品分类表
- `dishes` - 菜品表
- `dish_images` - 菜品图片表
- `orders` - 订单表
- `order_items` - 订单详情表
- `reviews` - 评价表
- `review_images` - 评价图片表
- `announcements` - 公告表

## 📡 API接口

### 用户端
- `POST /api/auth/register` - 注册
- `POST /api/auth/login` - 登录
- `GET /api/dishes` - 菜品列表
- `GET /api/dishes/:id` - 菜品详情
- `GET /api/categories` - 分类列表
- `GET /api/announcements/active` - 有效公告
- `POST /api/orders` - 创建订单
- `POST /api/orders/reserve` - 创建预约
- `POST /api/reviews` - 提交评价

### 管理端
- `POST /api/admin/login` - 管理员登录
- `POST /api/admin/dishes` - 新增菜品
- `PUT /api/admin/dishes/:id` - 编辑菜品
- `DELETE /api/admin/dishes/:id` - 删除菜品
- `POST /api/admin/upload` - 上传图片
- `POST /api/admin/announcements` - 新增公告
- 更多API请查看 `backend/router/router.go`

## 🔧 配置说明

### 后端配置 (backend/config.yaml)

```yaml
server:
  port: 8080
  mode: debug

database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "your_password"
  database: "orderfood_db"

jwt:
  secret: "your-secret-key"
  expire_hours: 168

storage:
  local_path: "./uploads"
  image_max_size: 5242880

announcement:
  auto_play_interval: 3000
  max_display_count: 10
```

## 📝 开发计划

- [x] 后端项目初始化
- [x] 移动端前端初始化
- [x] 用户认证系统
- [x] 菜品管理功能
- [x] 图片上传和管理
- [x] 订单和购物车功能
- [x] 预约点餐功能
- [x] 评价系统
- [x] 公告系统
- [x] PC端管理后台
- [ ] 部署上线

## 📄 License

MIT

## 👥 贡献

欢迎提交 Issue 和 Pull Request！

---

**开发日期**: 2025-11-13  
**版本**: v1.0.0
