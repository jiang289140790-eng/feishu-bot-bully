package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkevent "github.com/larksuite/oapi-sdk-go/v3/event"
	"github.com/larksuite/oapi-sdk-go/v3/event/dispatcher"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/larksuite/oapi-sdk-go/v3/ws"
)

// 机器人1配置
var client *lark.Client

func main() {
	// 从环境变量读取配置
	appID := os.Getenv("FEISHU_APP_ID")
	appSecret := os.Getenv("FEISHU_APP_SECRET")

	// 如果环境变量为空，使用默认值（机器人1）
	if appID == "" {
		appID = "cli_a9f1dae58d39dcd6"
		log.Println("⚠️  未设置 FEISHU_APP_ID 环境变量，使用默认值（机器人1）")
	}
	if appSecret == "" {
		appSecret = "wKp9u9Ys2YhtaPSotuOoheIdPBJFp0za"
		log.Println("⚠️  未设置 FEISHU_APP_SECRET 环境变量，使用默认值（机器人1）")
	}

	log.Printf("🤖 机器人1 App ID: %s", appID)

	// 创建飞书客户端
	client = lark.NewClient(appID, appSecret)

	// 创建事件处理器
	handler := dispatcher.NewEventDispatcher("", "")
	
	// 注册消息接收事件
	handler.OnCustomizedEvent("im.message.receive_v1", func(ctx context.Context, eventReq *larkevent.EventReq) error {
		log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		log.Println("💬 收到消息事件！")
		
		// 解析 Body
		if len(eventReq.Body) > 0 {
			var bodyData map[string]interface{}
			if err := json.Unmarshal(eventReq.Body, &bodyData); err == nil {
				if event, ok := bodyData["event"].(map[string]interface{}); ok {
					if message, ok := event["message"].(map[string]interface{}); ok {
						messageId := getString(message, "message_id")
						content := getString(message, "content")
						
						log.Printf("   MessageID: %s", messageId)
						log.Printf("   Content: %s", content)
						
						// 解析文本
						var contentMap map[string]interface{}
						if err := json.Unmarshal([]byte(content), &contentMap); err == nil {
							if text, ok := contentMap["text"].(string); ok {
								log.Printf("   文本: %s", text)
								// 回复消息（机器人1的回复）
								go replyMessage(messageId, fmt.Sprintf("【机器人1】收到你的消息：%s", text))
							}
						}
					}
				}
			}
		}
		
		log.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		return nil
	})

	// 创建 WebSocket 客户端
	cli := ws.NewClient(appID, appSecret,
		ws.WithLogLevel(larkcore.LogLevelInfo),
		ws.WithEventHandler(handler),
	)

	log.Println("🚀 正在启动飞书事件长链接监听...")
	log.Println("📱 这是机器人1的服务")

	// 启动长链接
	err := cli.Start(context.Background())
	if err != nil {
		log.Fatalf("❌ 启动失败: %v", err)
	}

	log.Println("✅ 长链接已成功建立，正在监听事件...")
	log.Println("📝 监听事件类型: im.message.receive_v1")
	log.Println("提示: 按 Ctrl+C 退出程序")

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏹️  正在关闭连接...")
	log.Println("👋 程序已退出")
}

// 回复消息
func replyMessage(messageId, text string) {
	if messageId == "" {
		log.Println("⚠️  消息 ID 为空，跳过回复")
		return
	}

	log.Printf("📤 准备回复消息: %s", text)

	req := larkim.NewReplyMessageReqBuilder().
		MessageId(messageId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType("text").
			Content(fmt.Sprintf(`{"text":"%s"}`, text)).
			Build()).
		Build()

	resp, err := client.Im.Message.Reply(context.Background(), req)
	if err != nil {
		log.Printf("❌ 回复失败: %v", err)
		return
	}

	if resp.Success() {
		log.Printf("✅ 回复成功: %s", text)
	} else {
		log.Printf("❌ 回复失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}
}

// 辅助函数：从 map 中获取字符串
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
