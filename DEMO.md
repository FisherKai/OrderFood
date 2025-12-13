# 🎬 点餐系统演示指南

## 📱 移动端用户操作流程

### 场景1: 新用户注册并点餐

#### 步骤1: 注册账号
1. 访问 http://localhost:3000
2. 点击"注册账号"
3. 填写信息:
   - 用户名: testuser
   - 密码: 123456
   - 手机号: 13800138000
   - 邮箱: test@example.com
   - 昵称: 测试用户
4. 点击"注册"

#### 步骤2: 登录系统
1. 返回登录页
2. 输入用户名和密码
3. 登录成功后进入首页

#### 步骤3: 浏览菜品
1. 首页可以看到:
   - 顶部公告轮播
   - 分类导航
   - 推荐菜品
2. 点击"查看全部"或底部"菜品"标签
3. 可以:
   - 搜索菜品
   - 按分类筛选
   - 滚动加载更多

#### 步骤4: 查看菜品详情
1. 点击任意菜品
2. 查看:
   - 菜品图片轮播
   - 价格和描述
   - 用户评价
3. 点击"加入购物车"

#### 步骤5: 管理购物车
1. 点击底部"购物车"标签
2. 可以:
   - 调整菜品数量
   - 删除菜品(左滑)
   - 选择要结算的菜品
3. 点击"结算"

#### 步骤6: 查看订单
1. 订单创建成功
2. 自动跳转到"我的订单"
3. 查看订单状态和详情

---

### 场景2: 预约点餐

#### 步骤1: 添加菜品到购物车
1. 浏览菜品并添加到购物车
2. 至少选择1个菜品

#### 步骤2: 创建预约
1. 点击个人中心 -> "预约点餐"
2. 选择预约时间(至少提前2小时)
3. 输入用餐人数
4. 查看已选菜品和总价
5. 点击"提交预约"

#### 步骤3: 查看预约订单
1. 在"我的订单"中查看
2. 预约订单会显示预约时间

---

### 场景3: 评价菜品

#### 步骤1: 完成订单
1. 订单状态必须是"已完成"
2. 管理员需要将订单状态更新为3

#### 步骤2: 提交评价
1. (功能已实现，可通过API直接调用)
2. POST /api/reviews
```json
{
  "dish_id": 1,
  "order_id": 1,
  "rating": 5,
  "content": "非常好吃！"
}
```

#### 步骤3: 查看评价
1. 进入菜品详情页
2. 滚动到评价区域
3. 查看其他用户的评价

---

## 🔧 管理员操作流程

### 使用Postman/API工具测试

#### 1. 管理员登录

```bash
POST http://localhost:8080/api/admin/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**响应:**
```json
{
  "message": "登录成功",
  "token": "eyJhbGc...",
  "admin": {
    "id": 1,
    "username": "admin",
    "role": "admin"
  }
}
```

保存返回的 token，后续请求需要在 Header 中添加：
```
Authorization: Bearer <token>
```

---

#### 2. 创建菜品分类

```bash
POST http://localhost:8080/api/admin/categories
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "川菜",
  "icon": "",
  "sort": 10
}
```

---

#### 3. 上传菜品图片

```bash
POST http://localhost:8080/api/admin/upload
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: [选择图片文件]
```

**响应:**
```json
{
  "message": "上传成功",
  "url": "/uploads/1234567890_abc.jpg"
}
```

---

#### 4. 创建菜品

```bash
POST http://localhost:8080/api/admin/dishes
Authorization: Bearer <token>
Content-Type: application/json

{
  "category_id": 1,
  "name": "宫保鸡丁",
  "price": 38.00,
  "description": "经典川菜，香辣可口",
  "status": 1,
  "stock": 100
}
```

**响应后获得菜品ID，再添加图片:**

```bash
# 需要先在数据库手动插入，或通过代码实现
INSERT INTO dish_images (dish_id, image_url, is_main, sort) 
VALUES (1, '/uploads/1234567890_abc.jpg', true, 0);
```

---

#### 5. 查看所有订单

```bash
GET http://localhost:8080/api/admin/orders?page=1&page_size=20
Authorization: Bearer <token>
```

---

#### 6. 更新订单状态

```bash
PUT http://localhost:8080/api/admin/orders/1/status
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": 2
}
```

状态码:
- 1: 待处理
- 2: 制作中
- 3: 已完成
- 4: 已取消

---

#### 7. 发布公告

```bash
POST http://localhost:8080/api/admin/announcements
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "双十一优惠活动",
  "content": "全场8折，欢迎选购！",
  "type": 2,
  "start_time": "2025-11-11T00:00:00Z",
  "end_time": "2025-11-11T23:59:59Z",
  "status": 1,
  "sort": 10
}
```

公告类型:
- 1: 普通
- 2: 重要
- 3: 紧急

---

#### 8. 查看所有评价

```bash
GET http://localhost:8080/api/admin/reviews?page=1&page_size=20
Authorization: Bearer <token>
```

---

#### 9. 软删除菜品图片

```bash
DELETE http://localhost:8080/api/admin/images/1
Authorization: Bearer <token>
```

---

#### 10. 恢复历史图片

```bash
PUT http://localhost:8080/api/admin/dishes/1/images/restore
Authorization: Bearer <token>
Content-Type: application/json

{
  "image_id": 1
}
```

---

## 🧪 完整测试流程

### 准备测试数据

```sql
-- 1. 创建测试分类
INSERT INTO categories (name, icon, sort) VALUES
('热门推荐', '', 100),
('川菜', '', 90),
('粤菜', '', 80),
('湘菜', '', 70);

-- 2. 创建测试菜品
INSERT INTO dishes (category_id, name, price, description, status, stock) VALUES
(1, '宫保鸡丁', 38.00, '经典川菜，香辣可口', 1, 100),
(1, '鱼香肉丝', 32.00, '酸甜适中，色香味俱全', 1, 100),
(2, '白切鸡', 58.00, '皮爽肉滑，清淡鲜美', 1, 50),
(2, '清蒸鲈鱼', 68.00, '鲜嫩美味，营养丰富', 1, 30);

-- 3. 添加菜品图片(需要先上传图片)
INSERT INTO dish_images (dish_id, image_url, is_main, sort) VALUES
(1, '/uploads/gongbao.jpg', true, 0),
(2, '/uploads/yuxiang.jpg', true, 0),
(3, '/uploads/baiqie.jpg', true, 0),
(4, '/uploads/luyu.jpg', true, 0);
```

---

## 📊 测试检查清单

### 用户端功能
- [ ] 用户注册
- [ ] 用户登录
- [ ] 查看首页公告
- [ ] 浏览菜品列表
- [ ] 搜索菜品
- [ ] 分类筛选
- [ ] 查看菜品详情
- [ ] 添加购物车
- [ ] 管理购物车
- [ ] 创建订单
- [ ] 预约点餐
- [ ] 查看订单列表
- [ ] 退出登录

### 管理端功能
- [ ] 管理员登录
- [ ] 创建分类
- [ ] 创建菜品
- [ ] 上传图片
- [ ] 编辑菜品
- [ ] 删除菜品
- [ ] 查看订单
- [ ] 更新订单状态
- [ ] 发布公告
- [ ] 编辑公告
- [ ] 删除公告
- [ ] 查看评价

---

## 🎥 演示建议

### 演示顺序
1. **展示移动端首页** - 公告轮播、分类、推荐菜品
2. **注册新用户** - 演示注册流程
3. **浏览和搜索** - 展示菜品列表功能
4. **加入购物车** - 演示购物流程
5. **创建订单** - 完成下单
6. **预约点餐** - 演示预约功能
7. **API演示** - 使用Postman展示管理员功能

### 演示要点
- 强调**响应式设计**
- 展示**公告轮播**效果
- 演示**购物车**实时计算
- 说明**预约时间限制**(至少提前2小时)
- 展示**图片历史管理**功能
- 强调**前后端分离**架构

---

## 🔍 调试技巧

### 查看后端日志
```bash
cd backend
go run main.go
# 观察控制台输出的SQL和请求日志
```

### 查看前端Network
1. 打开浏览器开发者工具
2. 切换到Network标签
3. 筛选XHR请求
4. 查看请求和响应数据

### 数据库调试
```bash
mysql -u root -p orderfood_db

# 查看订单
SELECT * FROM orders ORDER BY created_at DESC LIMIT 10;

# 查看用户
SELECT * FROM users;

# 查看菜品
SELECT d.*, c.name as category_name FROM dishes d 
LEFT JOIN categories c ON d.category_id = c.id;
```

---

**祝演示成功！** 🎉
