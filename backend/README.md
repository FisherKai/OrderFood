# OrderFood Backend

点餐系统后端服务

## 技术栈

- Go 1.21+
- Gin Web Framework
- GORM (ORM)
- MySQL
- JWT认证

## 快速开始

### 1. 安装依赖

```bash
go mod download
```

### 2. 配置数据库

复制 `config.yaml` 并修改数据库配置：

```yaml
database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "your_password"
  database: "orderfood_db"
```

### 3. 创建数据库

```sql
CREATE DATABASE orderfood_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. 运行服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动

## API文档

### 用户端

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录
- `GET /api/dishes` - 获取菜品列表
- `GET /api/dishes/:id` - 获取菜品详情
- `GET /api/categories` - 获取分类列表
- `GET /api/announcements/active` - 获取有效公告
- `POST /api/orders` - 创建订单（需认证）
- `POST /api/orders/reserve` - 创建预约订单（需认证）
- `POST /api/reviews` - 提交评价（需认证）

### 管理端

- `POST /api/admin/login` - 管理员登录
- `POST /api/admin/dishes` - 新增菜品（需认证）
- `PUT /api/admin/dishes/:id` - 编辑菜品（需认证）
- `DELETE /api/admin/dishes/:id` - 删除菜品（需认证）
- `POST /api/admin/upload` - 上传图片（需认证）
- `POST /api/admin/announcements` - 新增公告（需认证）

## 项目结构

```
backend/
├── config/          # 配置管理
├── controllers/     # 控制器
├── database/        # 数据库连接
├── middleware/      # 中间件
├── models/          # 数据模型
├── router/          # 路由
├── utils/           # 工具函数
├── uploads/         # 上传文件目录
├── config.yaml      # 配置文件
├── main.go          # 入口文件
└── go.mod           # 依赖管理
```

## 创建默认管理员

首次运行后，可以通过以下SQL创建默认管理员账户：

```sql
-- 密码: admin123 (已加密)
INSERT INTO admins (username, password, role, created_at, updated_at) 
VALUES ('admin', '$2a$10$XBq7vXz5VJdY.uqFQqWQZ.oQE5p5kC8mKXZvXQJ5K6Y5qX5Q5Q5Q5', 'admin', NOW(), NOW());
```

注意：实际使用时需要用 bcrypt 加密后的密码。
