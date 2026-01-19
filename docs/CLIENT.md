# 客户端接入文档

## 概述

本平台支持多种自动化工具的远程调试，通过 WebSocket 协议连接到分配的端口进行通信。

## 支持的协议

| 协议 | 说明 |
|------|------|
| autojs | Auto.js / Auto.js Pro |
| lazyman | 懒人精灵 |
| easyclick | EasyClick |
| generic | 通用 WebSocket |

## 连接流程

1. 在平台申请端口，获取端口号和连接密钥
2. 客户端通过 WebSocket 连接到端口
3. 发送认证消息（包含密钥和设备信息）
4. 认证成功后保持心跳
5. 接收服务端下发的命令

## Auto.js 接入

### 基础连接

```javascript
// 配置
var CONFIG = {
    host: "your-server.com",
    port: 10086,
    token: "your-connection-token"
};

// 创建连接
var ws = new WebSocket("ws://" + CONFIG.host + ":" + CONFIG.port);

// 连接成功
ws.on("open", function() {
    log("连接成功");

    // 发送认证消息
    ws.send(JSON.stringify({
        type: "register",
        token: CONFIG.token,
        device: {
            id: device.getAndroidId(),
            name: device.brand + " " + device.model,
            brand: device.brand,
            model: device.model,
            os: "Android " + device.release,
            version: app.versionName
        }
    }));
});

// 接收消息
ws.on("message", function(msg) {
    try {
        var data = JSON.parse(msg);
        handleMessage(data);
    } catch(e) {
        log("解析消息失败: " + e);
    }
});

// 连接关闭
ws.on("close", function(code, reason) {
    log("连接关闭: " + code + " " + reason);
    // 可以在这里实现重连逻辑
});

// 错误处理
ws.on("error", function(err) {
    log("连接错误: " + err);
});

// 消息处理
function handleMessage(data) {
    switch(data.type) {
        case "auth_success":
            log("认证成功");
            startHeartbeat();
            break;

        case "auth_failed":
            log("认证失败: " + data.message);
            break;

        case "exec":
            // 执行远程代码
            executeCode(data.code);
            break;

        case "screenshot":
            // 截图请求
            takeScreenshot();
            break;

        case "pong":
            // 心跳响应
            break;
    }
}

// 执行代码
function executeCode(code) {
    try {
        var result = eval(code);
        sendLog("info", "执行成功: " + result);
    } catch(e) {
        sendLog("error", "执行失败: " + e);
    }
}

// 截图
function takeScreenshot() {
    var img = captureScreen();
    if (img) {
        var base64 = images.toBase64(img);
        ws.send(JSON.stringify({
            type: "screenshot",
            data: base64
        }));
        img.recycle();
    }
}

// 发送日志
function sendLog(level, message) {
    ws.send(JSON.stringify({
        type: "log",
        level: level,
        message: message
    }));
}

// 心跳
var heartbeatTimer = null;
function startHeartbeat() {
    heartbeatTimer = setInterval(function() {
        if (ws.readyState === WebSocket.OPEN) {
            ws.send(JSON.stringify({ type: "ping" }));
        }
    }, 30000);
}
```

### 自动重连

```javascript
var reconnectInterval = 5000;
var maxReconnects = 10;
var reconnectCount = 0;

function connect() {
    var ws = new WebSocket("ws://" + CONFIG.host + ":" + CONFIG.port);

    ws.on("open", function() {
        reconnectCount = 0;
        // ... 认证逻辑
    });

    ws.on("close", function() {
        if (reconnectCount < maxReconnects) {
            reconnectCount++;
            log("尝试重连... (" + reconnectCount + "/" + maxReconnects + ")");
            setTimeout(connect, reconnectInterval);
        }
    });

    return ws;
}
```

## 懒人精灵接入

```lua
-- 配置
local config = {
    host = "your-server.com",
    port = 10086,
    token = "your-connection-token"
}

-- 创建连接
local ws = WebSocket.new("ws://" .. config.host .. ":" .. config.port)

-- 连接成功
ws:on("open", function()
    print("连接成功")

    -- 发送认证消息
    ws:send(json.encode({
        type = "register",
        token = config.token,
        device = {
            udid = Device.udid(),
            device_name = Device.name(),
            model = Device.model(),
            os_version = Device.osVersion(),
            app_version = App.version()
        }
    }))
end)

-- 接收消息
ws:on("message", function(msg)
    local data = json.decode(msg)

    if data.type == "auth_success" then
        print("认证成功")
    elseif data.type == "exec" then
        -- 执行远程代码
        local func = load(data.code)
        if func then
            func()
        end
    elseif data.type == "screenshot" then
        -- 截图
        local img = Screen.capture()
        ws:send(json.encode({
            type = "screenshot",
            data = img:base64()
        }))
    end
end)

-- 心跳
Timer.setInterval(function()
    ws:send(json.encode({ type = "ping" }))
end, 30000)

-- 连接
ws:connect()
```

## EasyClick 接入

```javascript
// 配置
var config = {
    host: "your-server.com",
    port: 10086,
    token: "your-connection-token"
};

// 创建连接
var ws = websocket.create("ws://" + config.host + ":" + config.port);

// 连接
ws.connect({
    onOpen: function() {
        logd("连接成功");

        // 发送认证
        ws.send(JSON.stringify({
            type: "register",
            token: config.token,
            device: {
                deviceId: device.getIMEI(),
                deviceName: device.getBrand() + " " + device.getModel(),
                brand: device.getBrand(),
                model: device.getModel(),
                osVersion: device.getOSVersion(),
                sdkVersion: device.getSDKInt()
            }
        }));
    },

    onMessage: function(msg) {
        var data = JSON.parse(msg);

        if (data.type == "exec") {
            eval(data.code);
        }
    },

    onClose: function(code, reason) {
        logd("连接关闭: " + code);
    }
});
```

## 通用协议

如果你的自动化工具不在上述列表中，可以使用通用协议接入。

### 认证消息格式

```json
{
  "type": "register",
  "token": "端口连接密钥",
  "device": {
    "id": "设备唯一标识",
    "name": "设备名称",
    "brand": "品牌（可选）",
    "model": "型号（可选）",
    "os": "操作系统（可选）",
    "version": "应用版本（可选）"
  }
}
```

### 心跳消息

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
  "level": "info|warn|error",
  "message": "日志内容"
}
```

### 接收执行命令

```json
{
  "type": "exec",
  "code": "要执行的代码"
}
```

### 执行结果上报

```json
{
  "type": "result",
  "data": "执行结果"
}
```

## 常见问题

### 1. 连接被拒绝

- 检查端口号是否正确
- 检查端口是否已过期
- 检查服务器防火墙设置

### 2. 认证失败

- 检查连接密钥是否正确
- 检查密钥是否已更换

### 3. 设备显示离线

- 检查心跳是否正常发送
- 检查网络连接是否稳定
