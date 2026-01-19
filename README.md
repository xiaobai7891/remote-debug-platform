# 🚀 远程调试 SaaS 平台

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://golang.org)
[![Vue](https://img.shields.io/badge/Vue-3.4+-4FC08D?logo=vue.js)](https://vuejs.org)

一个支持 Auto.js、懒人精灵、EasyClick 等自动化工具的远程调试平台，支持多用户、端口动态分配、计费等功能。

![Dashboard Preview](docs/images/dashboard.png)

## ✨ 功能特性

- 🔐 **用户系统** - 注册、登录、JWT 认证、权限管理
- 🔌 **端口管理** - 动态分配、自动续费、一键更换、状态监控
- 📱 **设备管理** - 实时上下线、延迟检测、远程执行代码、日志查看
- 💰 **计费系统** - 余额管理、消费记录、灵活定价
- 🛠 **管理后台** - 用户管理、充值、数据统计、系统监控
- 🔗 **多协议支持** - Auto.js、懒人精灵、EasyClick、通用 WebSocket

## 🚀 一键部署

### Linux (推荐)

```bash
# 克隆项目
git clone https://github.com/your-username/remote-debug-platform.git
cd remote-debug-platform

# 一键部署
chmod +x install.sh
sudo ./install.sh
```

### Windows

```bash
# 克隆项目后，双击运行
install.bat
```

### 手动部署

```bash
# 使用 Docker Compose
docker-compose up -d

# 访问 http://localhost
```

## 📦 本地开发

### 后端

```bash
cd backend
go mod tidy
go run main.go
# 服务运行在 http://localhost:3000
```

### 前端

```bash
cd frontend
npm install
npm run dev
# 访问 http://localhost:5173
```

## ⚙️ 配置说明

编辑 `backend/config.yaml` 或 `deploy/config.yaml`:

```yaml
server:
  http_port: 3000           # API 端口

port_pool:
  min: 10000                # 端口池起始
  max: 60000                # 端口池结束

pricing:
  port_per_day: 1.0         # 每端口每天价格

admin:
  username: "admin"         # 管理员用户名
  password: "admin123"      # 管理员密码（首次部署后请修改）
```

## 📱 客户端接入

### Auto.js

```javascript
var ws = new WebSocket("ws://your-server:10086");

ws.on("open", function() {
    ws.send(JSON.stringify({
        type: "register",
        token: "your-token",
        device: {
            id: device.getAndroidId(),
            name: device.brand + " " + device.model,
            os: "Android " + device.release
        }
    }));
});

ws.on("message", function(msg) {
    var data = JSON.parse(msg);
    if (data.type === "exec") {
        eval(data.code);
    }
});
```

更多客户端示例请查看 [客户端接入文档](docs/CLIENT.md)

## 📖 文档

- [API 文档](docs/API.md)
- [客户端接入文档](docs/CLIENT.md)
- [部署文档](docs/DEPLOY.md)

## 🏗 项目结构

```
remote-debug-platform/
├── backend/                # Go 后端
│   ├── handlers/           # API 处理器
│   ├── models/             # 数据模型
│   ├── services/           # 业务服务
│   └── protocols/          # 协议适配器
├── frontend/               # Vue 前端
│   └── src/views/          # 页面组件
├── deploy/                 # 部署配置
├── install.sh              # Linux 一键部署
└── install.bat             # Windows 一键部署
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

请阅读 [贡献指南](CONTRIBUTING.md) 了解详情。

## 📄 许可证

[MIT License](LICENSE)

## 🙏 鸣谢

- [Gin](https://github.com/gin-gonic/gin) - Go Web 框架
- [Vue.js](https://vuejs.org/) - 前端框架
- [Element Plus](https://element-plus.org/) - UI 组件库
