# OrderFood 项目改动记录

## 2025-12-13 更新

### 一、管理端移动端适配 ✅

实现了管理后台的移动端响应式适配，支持通过 UA 和屏幕宽度判断设备类型。

#### 新增文件

| 文件 | 说明 | 状态 |
|------|------|------|
| `frontend-admin/src/composables/useDevice.js` | 设备检测 composable（UA + 屏幕宽度） | ✅ 完成 |

#### 修改文件

| 文件 | 修改内容 | 状态 |
|------|----------|------|
| `frontend-admin/src/layout/Layout.vue` | 移动端抽屉式侧边栏、响应式头部 | ✅ 完成 |
| `frontend-admin/src/styles/main.scss` | 全局响应式样式、移动端组件样式 | ✅ 完成 |
| `frontend-admin/src/views/Dashboard.vue` | Dashboard 移动端适配、卡片式订单列表 | ✅ 完成 |
| `frontend-admin/index.html` | 添加移动端 meta 标签 | ✅ 完成 |

#### 功能特性

1. **设备检测**
   - 基于 User-Agent 检测移动设备
   - 基于屏幕宽度判断（< 768px 为移动端）
   - 自动响应窗口大小变化

2. **移动端布局**
   - 抽屉式侧边菜单（从左滑出）
   - 遮罩层点击关闭
   - 底部退出登录按钮
   - 简化的头部导航

3. **响应式组件**
   - 统计卡片自适应尺寸
   - 图表自动调整配置
   - 订单列表移动端使用卡片形式
   - 表格、弹窗、分页等全局适配

#### 使用示例

```javascript
import { useDevice } from '@/composables/useDevice'

const { isMobile, isTablet, isDesktop, screenWidth } = useDevice()

// 在模板中使用
<div :class="{ 'mobile': isMobile }">
  <template v-if="isMobile">移动端内容</template>
  <template v-else>桌面端内容</template>
</div>
```

---

### 二、Docker 打包方案 ✅

创建了完整的 Docker 容器化部署方案。

#### 新增文件

| 文件 | 说明 | 状态 |
|------|------|------|
| `Dockerfile` | 多阶段构建，前后端打包到一个镜像 | ✅ 完成 |
| `docker-compose.yml` | 编排应用服务 + MySQL 数据库 | ✅ 完成 |
| `docker/nginx.conf` | Nginx 反向代理配置 | ✅ 完成 |
| `docker/entrypoint.sh` | 启动脚本，从环境变量生成配置 | ✅ 完成 |
| `docker/build.sh` | 构建辅助脚本 | ✅ 完成 |
| `.env.example` | 环境变量模板 | ✅ 完成 |
| `.dockerignore` | Docker 构建排除文件 | ✅ 完成 |

#### 配置管理方式

通过环境变量注入配置，启动时自动生成 `config.yaml`：

```bash
# 主要环境变量
DB_HOST          # 数据库地址
DB_PORT          # 数据库端口（默认3306）
DB_USER          # 数据库用户（默认root）
DB_PASSWORD      # 数据库密码（必填）
DB_NAME          # 数据库名（默认orderfood_db）
JWT_SECRET       # JWT密钥（必填）
SERVER_MODE      # 运行模式（默认release）
```

#### 使用方式

```bash
# docker-compose 一键启动
cp .env.example .env
# 编辑 .env 修改密码和密钥
docker-compose up -d

# 或单独构建镜像
./docker/build.sh
```

#### 访问地址

| 服务 | 地址 |
|------|------|
| 移动端 | http://localhost:3000 |
| 管理后台 | http://localhost:3001 |
| 统一入口 | http://localhost |
| 后端API | http://localhost:8080 |

---

### 二、中文乱码修复 ✅

解决 Docker 部署后中文显示乱码问题。

#### 修改文件

| 文件 | 修改内容 | 状态 |
|------|----------|------|
| `docker/nginx.conf` | 添加 `charset utf-8;` 字符集配置 | ✅ 完成 |
| `Dockerfile` | 添加 `ENV LANG=C.UTF-8` 和 `LC_ALL=C.UTF-8` | ✅ 完成 |
| `docker/entrypoint.sh` | 添加 `export LANG=C.UTF-8` | ✅ 完成 |
| `docker-compose.yml` | MySQL 添加 `--character-set-server=utf8mb4` 启动参数 | ✅ 完成 |

---

### 三、一周菜谱功能修复 ✅

解决移动端一周菜谱页面无法显示的问题。

#### 问题1：后端时间查询逻辑

**问题原因**：原代码精确匹配 `week_start = 周一日期`，但数据库存储的时间包含时区信息，导致匹配失败。

**解决方案**：改用范围查询。

| 文件 | 修改内容 | 状态 |
|------|----------|------|
| `backend/controllers/weekly_menu.go` | `GetCurrentWeekMenu` 函数改用范围查询 | ✅ 完成 |
| `backend/controllers/weekly_menu.go` | `GetWeekMenuByDate` 函数改用范围查询 | ✅ 完成 |

**修改前**：
```go
WHERE DATE(week_start) = '计算出的周一日期' AND status = 已发布
```

**修改后**：
```go
WHERE DATE(请求日期) >= DATE(week_start) 
  AND DATE(请求日期) <= DATE(week_end) 
  AND status = 已发布
```

#### 问题2：前端 dayjs 插件缺失

**问题原因**：dayjs 默认没有 `week()` 方法，需要引入插件。

**错误信息**：
```
TypeError: currentWeekStart.value.week is not a function
```

| 文件 | 修改内容 | 状态 |
|------|----------|------|
| `frontend-mobile/src/views/WeeklyMenu.vue` | 添加 dayjs 插件 `weekOfYear` 和 `isoWeek` | ✅ 完成 |

**新增代码**：
```javascript
import weekOfYear from 'dayjs/plugin/weekOfYear'
import isoWeek from 'dayjs/plugin/isoWeek'

dayjs.extend(weekOfYear)
dayjs.extend(isoWeek)

// 使用 ISO 标准方法
const currentWeekStart = ref(dayjs().isoWeekday(1)) // ISO 周一
const week = currentWeekStart.value.isoWeek()       // ISO 周数
```

---

### 四、其他机器使用镜像说明 ✅

#### 方式1：推送到 Docker Hub

```bash
docker login
docker tag orderfood:latest your-username/orderfood:latest
docker push your-username/orderfood:latest

# 其他机器拉取
docker pull your-username/orderfood:latest
```

#### 方式2：导出镜像文件（离线传输）

```bash
# 导出
docker save orderfood:latest | gzip > orderfood.tar.gz

# 传输后加载
docker load < orderfood.tar.gz
```

#### 方式3：私有镜像仓库

```bash
docker tag orderfood:latest registry.example.com/orderfood:latest
docker push registry.example.com/orderfood:latest
```

---

## 待办事项

- [ ] 生产环境部署测试
- [ ] 性能优化
- [ ] 日志收集配置

---

## 注意事项

1. 部署前务必修改 `.env` 文件中的密码和密钥
2. 首次部署需要执行 `init.sql` 初始化数据库
3. 上传的图片存储在 Docker volume `uploads_data` 中，需要做好备份
