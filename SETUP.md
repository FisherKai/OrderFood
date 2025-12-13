# 点餐系统快速设置指南

## 🚀 快速开始（5分钟搞定）

### 步骤1: 安装MySQL数据库

确保MySQL已安装并运行。如果没有，请先安装：

```bash
# macOS
brew install mysql
brew services start mysql

# 或使用Docker
docker run --name mysql -e MYSQL_ROOT_PASSWORD=123456 -p 3306:3306 -d mysql:8.0
```

### 步骤2: 初始化数据库

```bash
# 进入项目目录
cd /Users/shiyijiang/Documents/code/OrderFood

# 导入数据库结构和初始数据
mysql -u root -p < init.sql
```

### 步骤3: 生成管理员密码

```bash
cd backend/scripts
go run create_admin.go
```

复制输出的SQL语句，执行插入管理员账户。

### 步骤4: 配置后端

编辑 `backend/config.yaml`：

```yaml
database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "your_mysql_password"  # 修改为你的MySQL密码
  database: "orderfood_db"

jwt:
  secret: "your-secret-key-change-this"  # 修改为随机字符串
```

### 步骤5: 安装依赖

```bash
# 安装后端依赖
cd backend
go mod download

# 安装前端依赖
cd ../frontend-mobile
npm install
```

### 步骤6: 启动项目

**方式一：使用启动脚本（推荐）**

```bash
cd /Users/shiyijiang/Documents/code/OrderFood
./start.sh
```

**方式二：分别启动**

```bash
# 终端1 - 启动后端
cd backend
mkdir -p uploads
go run main.go

# 终端2 - 启动移动端前端
cd frontend-mobile
npm run dev
```

### 步骤7: 访问系统

- 📱 **移动端**: http://localhost:3000
- 📡 **后端API**: http://localhost:8080

## 📋 默认账户

### 用户账户
需要自己注册

### 管理员账户
- 用户名: `admin`
- 密码: `admin123`

## ✅ 功能测试清单

### 用户端测试
- [ ] 用户注册
- [ ] 用户登录
- [ ] 浏览菜品
- [ ] 查看公告
- [ ] 添加购物车
- [ ] 创建订单
- [ ] 预约点餐
- [ ] 提交评价

### 管理端测试（需要开发PC管理后台）
- [ ] 管理员登录
- [ ] 新增菜品
- [ ] 上传图片
- [ ] 管理分类
- [ ] 查看订单
- [ ] 发布公告

## 🛠️ 常见问题

### 1. 数据库连接失败
- 检查MySQL是否启动
- 检查 `config.yaml` 中的数据库配置
- 确认数据库 `orderfood_db` 已创建

### 2. 前端无法访问后端
- 确认后端已启动在8080端口
- 检查防火墙设置
- 查看浏览器控制台错误信息

### 3. 图片上传失败
- 确认 `backend/uploads` 目录存在
- 检查目录写入权限
- 确认文件大小不超过5MB

### 4. Go依赖下载慢
```bash
# 设置Go代理
go env -w GOPROXY=https://goproxy.cn,direct
```

### 5. npm安装慢
```bash
# 设置npm镜像
npm config set registry https://registry.npmmirror.com
```

## 📦 项目打包

### 后端打包

```bash
cd backend
go build -o orderfood main.go
```

### 前端打包

```bash
cd frontend-mobile
npm run build
```

生成的文件在 `dist` 目录。

## 🌐 生产环境部署

### 后端部署

1. 修改 `config.yaml` 中的 `mode` 为 `release`
2. 配置生产环境数据库
3. 运行编译后的二进制文件

```bash
./orderfood
```

### 前端部署

1. 构建生产版本：`npm run build`
2. 将 `dist` 目录部署到Nginx或其他Web服务器
3. 配置反向代理到后端API

## 📝 下一步

- [ ] 开发PC端管理后台
- [ ] 添加支付功能
- [ ] 添加短信通知
- [ ] 性能优化
- [ ] 安全加固

## 🆘 获取帮助

如有问题，请查看：
1. 项目README.md
2. 需求文档 需求.md
3. 提交Issue

---

**祝你使用愉快！** 🎉
