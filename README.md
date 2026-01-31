# 飞书事件监听器（长链接模式）

这是一个使用 Go 语言实现的飞书机器人事件监听器，采用长链接（WebSocket）方式接收事件。

## 功能特性

- ✅ 使用长链接方式接收飞书事件（无需公网 IP 和域名）
- ✅ 自动处理消息接收事件
- ✅ 支持自动回复功能
- ✅ 完整的错误处理和日志记录
- ✅ 优雅的启动和关闭

## 前置要求

1. Go 1.21 或更高版本
2. 飞书应用的 App ID 和 App Secret
3. 已在飞书开放平台配置事件订阅

## 快速开始

### 1. 修改配置

编辑 `main.go` 文件，替换以下配置：

```go
const (
    APP_ID     = "your_app_id"     // 替换为您的 App ID
    APP_SECRET = "your_app_secret" // 替换为您的 App Secret
)
```

### 2. 安装依赖

```bash
go mod tidy
```

### 3. 运行程序

```bash
go run main.go
```

如果一切正常，您会看到：

```
正在启动飞书事件长链接监听...
✅ 长链接已成功建立，正在监听事件...
```

### 4. 在飞书开放平台配置

1. 进入飞书开放平台 -> 您的应用 -> 事件与回调
2. 选择「事件配置」标签页
3. 订阅方式选择：**使用长链接接收事件（推荐）**
4. 点击「保存」

### 5. 订阅事件

在「事件订阅」中订阅您需要的事件，例如：
- `im.message.receive_v1` - 接收消息
- `application.bot.menu_v6` - 机器人菜单

## 事件处理

程序当前支持处理以下事件：

### 1. 接收消息事件 (`im.message.receive_v1`)

当用户向机器人发送消息时：
- 记录消息内容、类型、发送者等信息
- 自动回复"收到您的消息了！"

您可以在 `handleMessageReceive` 函数中自定义处理逻辑。

### 2. 添加新的事件处理

在 `main()` 函数的 `switch` 语句中添加新的 case：

```go
case "your.event.type":
    // 您的处理逻辑
    handleYourEvent(ctx, client, eventReq)
```

## 项目结构

```
feishu-event-listener/
├── main.go          # 主程序文件
├── go.mod           # Go 模块定义
├── go.sum           # 依赖校验文件
└── README.md        # 项目说明文档
```

## 常见问题

### Q: 长链接断开怎么办？
A: SDK 会自动重连，您不需要手动处理。

### Q: 如何调试事件？
A: 查看控制台日志，所有接收到的事件都会被打印出来。

### Q: 可以同时监听多个应用吗？
A: 需要为每个应用创建独立的客户端和 WebSocket 连接。

### Q: 如何部署到服务器？
A: 编译后直接运行：
```bash
go build -o feishu-listener main.go
./feishu-listener
```

建议使用 systemd 或 supervisor 管理进程。

## 技术支持

- 飞书开放平台文档：https://open.feishu.cn/document/
- Go SDK 文档：https://github.com/larksuite/oapi-sdk-go

## 许可证

MIT License
