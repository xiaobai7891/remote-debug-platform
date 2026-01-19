# API 文档

## 基础信息

- 基础路径: `/api`
- 认证方式: Bearer Token (JWT)
- 响应格式: JSON

### 通用响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

### 错误码

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 参数错误 |
| 401 | 未授权 |
| 403 | 禁止访问 |
| 404 | 资源不存在 |
| 500 | 服务器错误 |
| 1001 | 用户名已存在 |
| 1002 | 邮箱已被使用 |
| 1003 | 用户名或密码错误 |
| 1004 | 原密码错误 |
| 2001 | 已达到端口数量上限 |
| 2002 | 余额不足 |
| 2003 | 无可用端口 |
| 3001 | 设备不在线 |
| 3002 | 发送命令失败 |

## 认证接口

### 注册

```
POST /api/auth/register
```

请求体:
```json
{
  "username": "string",
  "password": "string",
  "email": "string"
}
```

响应:
```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "jwt-token",
    "user": {
      "id": 1,
      "username": "test",
      "email": "test@example.com",
      "balance": 0,
      "role": "user"
    }
  }
}
```

### 登录

```
POST /api/auth/login
```

请求体:
```json
{
  "username": "string",
  "password": "string"
}
```

## 端口接口

### 获取端口列表

```
GET /api/ports
Authorization: Bearer <token>
```

### 申请新端口

```
POST /api/ports
Authorization: Bearer <token>
```

请求体:
```json
{
  "protocol": "autojs",
  "name": "端口名称",
  "days": 30
}
```

### 续费端口

```
POST /api/ports/:id/renew
Authorization: Bearer <token>
```

请求体:
```json
{
  "days": 30
}
```

### 更换端口

```
POST /api/ports/:id/change
Authorization: Bearer <token>
```

### 重新生成 Token

```
POST /api/ports/:id/regenerate-token
Authorization: Bearer <token>
```

## 设备接口

### 获取设备列表

```
GET /api/devices
Authorization: Bearer <token>
```

### 获取设备日志

```
GET /api/devices/:id/logs?limit=100
Authorization: Bearer <token>
```

### 远程执行代码

```
POST /api/devices/:id/exec
Authorization: Bearer <token>
```

请求体:
```json
{
  "code": "toast('Hello!');"
}
```

## 管理接口

### 获取系统统计

```
GET /api/admin/stats
Authorization: Bearer <admin-token>
```

### 用户充值

```
POST /api/admin/users/:id/recharge
Authorization: Bearer <admin-token>
```

请求体:
```json
{
  "amount": 100,
  "description": "管理员充值"
}
```

## WebSocket 协议

### 连接

```
ws://host:port/
```

### 认证消息

```json
{
  "type": "register",
  "token": "端口连接密钥",
  "device": {
    "id": "设备ID",
    "name": "设备名称",
    "brand": "品牌",
    "model": "型号",
    "os": "系统版本",
    "version": "应用版本"
  }
}
```

### 心跳

发送:
```json
{"type": "ping"}
```

响应:
```json
{"type": "pong"}
```

### 日志上报

```json
{
  "type": "log",
  "level": "info",
  "message": "日志内容"
}
```

### 执行代码（服务端下发）

```json
{
  "type": "exec",
  "code": "要执行的代码"
}
```
