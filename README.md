# 海参链溯 SeaTrace

毕业设计示例：**Go（Gin + GORM）后端** + **Vue 3 + Tailwind** 前端 + **Hyperledger Fabric** 链码锚定溯源哈希。

## 功能概览

- 批次与溯源事件管理（JWT 角色：管理员 / 企业）
- 事件内容计算 SHA-256，调用账本接口写入 Fabric（未配置网络时使用 **Mock 交易号**）
- 公开溯源页：按批次号查询时间线与链上交易摘要
- 公开校验：`GET /api/verify/:batchNo` 比对哈希

## 目录结构

```
sea-cucumber-trace/
  backend/           # Go API
  frontend/          # Vue 3 + Vite + Tailwind
  chaincode/traceability/  # Fabric Go 链码骨架
```

## 本地运行

### 1. 后端

```powershell
cd backend
go run ./cmd/server
```

默认监听 `:8080`，SQLite 文件 `trace.db`（驱动为 **glebarez/sqlite**，纯 Go、无需 CGO）。演示账号：

| 用户 | 密码 | 说明 |
|------|------|------|
| admin | admin123 | 管理员 |
| orguser | org123 | 企业用户（绑定演示 orgId=1） |

### 2. 前端

```powershell
cd frontend
npm install
npm run dev
```

浏览器打开 `http://127.0.0.1:5173`（已配置将 `/api` 代理到后端）。

### 3. Hyperledger Fabric

1. 使用 Fabric 测试网络或自有联盟链，安装链码 `chaincode/traceability`（链码名如 `traceability`）。
2. 在后端 `internal/fabric/gateway.go` 中使用 **fabric-gateway**（Go）实现 `SubmitTransaction("PutTrace", payload)`，与链码 `PutTrace` 对齐。
3. 设置环境变量（示例）：

```text
FABRIC_ENABLED=true
FABRIC_CHANNEL=mychannel
FABRIC_CHAINCODE=traceability
FABRIC_PEER_ENDPOINT=localhost:7051
FABRIC_CERT_PATH=...
FABRIC_KEY_PATH=...
FABRIC_TLS_PATH=...
```

未接通 Fabric 时保持 `FABRIC_ENABLED=false`，系统使用 Mock 账本，便于答辩演示。

## API 摘要

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | /api/auth/login | 登录 |
| GET | /api/trace/:batchNo | 公开溯源时间线 |
| GET | /api/verify/:batchNo | 公开哈希校验 |
| GET | /api/batches | 需登录：批次列表 |
| POST | /api/batches | 需登录：创建批次 |
| POST | /api/batches/:id/events | 需登录：新增事件并上链锚定 |

## 数据库

开发环境使用 **SQLite**（`DB_PATH` 可改路径）。生产可换 **MySQL/PostgreSQL**：替换 GORM `dialector` 即可，表结构由 `AutoMigrate` 生成，亦可将 `internal/model` 映射为论文中的 ER 图。

## 许可证

示例代码用于毕业设计学习，按需修改。
