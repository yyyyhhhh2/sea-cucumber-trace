# SeaTrace Cloud

面向海参企业的追溯平台示例项目，包含：

- `backend`：Go + Gin + GORM API
- `frontend`：Vue 3 + Vite + Tailwind
- `mysql`：业务主数据库
- `redis`：缓存
- `chaincode/traceability`：Fabric 链码骨架

项目现在支持两种运行方式：

- 本地开发模式：前后端分别启动，数据库/缓存走 Docker
- 完整部署模式：前端、后端、MySQL、Redis 全部由 Docker Compose 拉起

## 1. 功能概览

- 企业登录与 JWT 鉴权
- 批次建档、追溯事件登记、批次时间线查询
- 公开消费者查询页与哈希校验页
- MySQL / SQLite 双数据库支持
- Redis 缓存
- Fabric 已接入抽象，未配置网络时自动降级为 Mock 账本
- 管理员一键导入演示数据
- Docker 一键启动完整项目

## 2. 演示账号

系统首次启动且数据库为空时，会自动写入演示数据。

| 用户名 | 密码 | 角色 |
|---|---|---|
| `admin` | `admin123` | 管理员 |
| `orguser` | `org123` | 企业演示账号 |
| `farm_user` | `farm123` | 养殖场操作员 |
| `process_user` | `process123` | 加工厂质检员 |
| `logistics_user` | `logistics123` | 冷链调度员 |
| `retail_user` | `retail123` | 门店管理员 |

## 3. 本地开发

### 3.1 启动 MySQL 和 Redis

```powershell
docker compose -f docker-compose.dev.yml up -d
```

### 3.2 启动后端

```powershell
cd backend
copy .env.example .env
go run ./cmd/server
```

默认配置：

- 后端地址：`http://127.0.0.1:8080`
- 数据库：MySQL
- Redis：开启

### 3.3 启动前端

```powershell
cd frontend
npm install
npm run dev
```

前端地址：

- `http://127.0.0.1:5173`

## 4. Docker 一键启动完整项目

先生成 Docker 运行所需本地产物：

```powershell
.\scripts\build-docker-assets.ps1
```

然后启动：

```powershell
docker compose up -d --build
```

启动后：

- 前端：`http://127.0.0.1:5180`
- 后端 API：`http://127.0.0.1:8081/api`
- MySQL：`127.0.0.1:3306`
- Redis：`127.0.0.1:6379`

默认 Compose 配置：

- MySQL 数据库名：`sea_cucumber_trace`
- MySQL 用户名：`seatrace`
- MySQL 密码：`seatrace123`
- Redis 开启 AOF 持久化

## 5. 后端环境变量

参考文件：

- [backend/.env.example](backend/.env.example)

关键变量：

| 变量 | 说明 |
|---|---|
| `PORT` | 后端端口 |
| `JWT_SECRET` | JWT 密钥 |
| `DB_DRIVER` | `sqlite` 或 `mysql` |
| `DB_PATH` | SQLite 文件路径 |
| `MYSQL_DSN` | MySQL 连接串 |
| `REDIS_ENABLED` | 是否启用 Redis |
| `REDIS_ADDR` | Redis 地址 |
| `CACHE_TTL_SECONDS` | 缓存 TTL |
| `FABRIC_ENABLED` | 是否启用 Fabric |

## 6. API 摘要

### 公开接口

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/api/health` | 服务健康检查 |
| `POST` | `/api/auth/login` | 登录 |
| `GET` | `/api/trace/:batchNo` | 公开追溯时间线 |
| `GET` | `/api/verify/:batchNo` | 公开哈希校验 |

### 登录后接口

| 方法 | 路径 | 说明 |
|---|---|---|
| `GET` | `/api/me` | 当前登录用户 |
| `GET` | `/api/orgs` | 机构列表 |
| `GET` | `/api/dashboard/summary` | 工作台摘要 |
| `GET` | `/api/batches` | 批次列表 |
| `POST` | `/api/batches` | 创建批次 |
| `GET` | `/api/batches/:id` | 批次详情 |
| `POST` | `/api/batches/:id/events` | 新增事件 |
| `GET` | `/api/batches/:id/timeline` | 批次时间线 |
| `POST` | `/api/admin/import-demo` | 管理员导入演示数据 |

## 7. 数据与缓存说明

### MySQL

正式运行建议使用 MySQL。当前项目：

- 已内置 GORM 自动建表
- 启动时自动迁移
- 启动时空库自动导入演示数据

### Redis

Redis 用于缓存：

- 批次列表
- 公开追溯页
- 公开校验结果

如果 Redis 不可用，系统会自动降级，不影响核心业务接口。

## 8. Fabric 说明

如需接入真实 Fabric：

1. 设置 `FABRIC_ENABLED=true`
2. 填入证书、私钥、TLS 和 Peer 配置
3. 链码需实现 `PutTrace`

未接通 Fabric 时，系统默认使用 Mock 账本返回交易号，适合毕业设计汇演。

## 9. 当前部署文件

- [docker-compose.yml](docker-compose.yml)：完整项目启动
- [docker-compose.dev.yml](docker-compose.dev.yml)：仅 MySQL + Redis 开发依赖
- [backend/Dockerfile](backend/Dockerfile)
- [frontend/Dockerfile](frontend/Dockerfile)
- [frontend/nginx.conf](frontend/nginx.conf)
