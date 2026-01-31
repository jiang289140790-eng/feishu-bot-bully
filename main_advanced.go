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
	larkws "github.com/larksuite/oapi-sdk-go/v3/ws"
	larkhelpdesk "github.com/larksuite/oapi-sdk-go/v3/service/helpdesk/v1"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
)

// 配置信息
const (
	APP_ID          = "your_app_id"          // 替换为您的 App ID
	APP_SECRET      = "your_app_secret"      // 替换为您的 App Secret
	HELPDESK_ID     = "your_helpdesk_id"     // 替换为您的服务台 ID（可选）
	HELPDESK_TOKEN  = "your_helpdesk_token"  // 替换为您的服务台 Token（可选）
)

func main() {
	// 创建飞书客户端
	client := lark.NewClient(APP_ID, APP_SECRET)

	// 创建 WebSocket 长链接客户端
	wsClient := larkws.NewClient(APP_ID, APP_SECRET,
		larkws.WithLogLevel(larkcore.LogLevelInfo),
		larkws.WithEventHandler(func(ctx context.Context, eventReq *larkws.EventReq) error {
			return handleEvent(ctx, client, eventReq)
		}),
	)

	// 启动长链接
	log.Println("🚀 正在启动飞书事件长链接监听...")
	err := wsClient.Start(context.Background())
	if err != nil {
		log.Fatalf("❌ 启动失败: %v", err)
	}

	log.Println("✅ 长链接已成功建立，正在监听事件...")
	log.Println("📝 支持的事件类型：")
	log.Println("   - im.message.receive_v1 (接收消息)")
	log.Println("   - helpdesk.ticket_v1 (服务台工单)")
	log.Println("   - helpdesk.ticket_message_v1 (工单消息)")

	// 等待退出信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⏹️  正在关闭连接...")
	wsClient.Close()
	log.Println("👋 程序已退出")
}

// 统一事件处理入口
func handleEvent(ctx context.Context, client *lark.Client, eventReq *larkws.EventReq) error {
	log.Printf("📨 收到事件: %s", eventReq.Header.EventType)

	switch eventReq.Header.EventType {
	case "im.message.receive_v1":
		// 处理接收消息事件
		return handleMessageReceive(ctx, client, eventReq)
	
	case "helpdesk.ticket_v1":
		// 处理服务台工单事件
		return handleHelpdeskTicket(ctx, client, eventReq)
	
	case "helpdesk.ticket_message_v1":
		// 处理工单消息事件
		return handleHelpdeskTicketMessage(ctx, client, eventReq)
	
	case "application.bot.menu_v6":
		// 处理机器人菜单事件
		log.Println("📋 收到机器人菜单事件")
		return nil
	
	default:
		log.Printf("⚠️  未处理的事件类型: %s", eventReq.Header.EventType)
		return nil
	}
}

// 处理接收到的消息
func handleMessageReceive(ctx context.Context, client *lark.Client, eventReq *larkws.EventReq) error {
	event := &larkim.P2ImMessageReceiveV1{}
	err := eventReq.Event.Unmarshal(event)
	if err != nil {
		return fmt.Errorf("解析消息事件失败: %w", err)
	}

	messageId := *event.Message.MessageId
	chatId := *event.Message.ChatId
	content := *event.Message.Content
	messageType := *event.Message.MessageType

	log.Printf("💬 收到消息:")
	log.Printf("   MessageID: %s", messageId)
	log.Printf("   ChatID: %s", chatId)
	log.Printf("   类型: %s", messageType)
	log.Printf("   内容: %s", content)

	// 解析消息内容
	var msgContent map[string]interface{}
	if err := json.Unmarshal([]byte(content), &msgContent); err == nil {
		if text, ok := msgContent["text"].(string); ok {
			log.Printf("   文本: %s", text)
			
			// 智能回复示例
			reply := generateReply(text)
			replyMessage(ctx, client, messageId, reply)
		}
	}

	return nil
}

// 处理服务台工单事件
func handleHelpdeskTicket(ctx context.Context, client *lark.Client, eventReq *larkws.EventReq) error {
	// 解析工单事件
	var ticketData map[string]interface{}
	if err := eventReq.Event.Unmarshal(&ticketData); err != nil {
		return fmt.Errorf("解析工单事件失败: %w", err)
	}

	log.Printf("🎫 服务台工单事件:")
	log.Printf("   数据: %+v", ticketData)

	// 这里可以添加工单处理逻辑
	// 例如：自动分配工单、发送通知等

	return nil
}

// 处理工单消息事件
func handleHelpdeskTicketMessage(ctx context.Context, client *lark.Client, eventReq *larkws.EventReq) error {
	var messageData map[string]interface{}
	if err := eventReq.Event.Unmarshal(&messageData); err != nil {
		return fmt.Errorf("解析工单消息事件失败: %w", err)
	}

	log.Printf("💬 工单消息事件:")
	log.Printf("   数据: %+v", messageData)

	return nil
}

// 生成智能回复
func generateReply(text string) string {
	// 这里可以添加更复杂的逻辑，比如调用 AI、查询数据库等
	switch text {
	case "你好", "您好", "hi", "hello":
		return "您好！我是飞书机器人，很高兴为您服务！"
	case "帮助", "help":
		return "我可以帮您处理以下内容：\n1. 回答常见问题\n2. 创建服务台工单\n3. 查询工单状态"
	default:
		return fmt.Sprintf("收到您的消息：%s\n我会尽快处理！", text)
	}
}

// 回复消息
func replyMessage(ctx context.Context, client *lark.Client, messageId, content string) {
	req := larkim.NewReplyMessageReqBuilder().
		MessageId(messageId).
		Body(larkim.NewReplyMessageReqBodyBuilder().
			MsgType("text").
			Content(fmt.Sprintf(`{"text":"%s"}`, content)).
			Build()).
		Build()

	resp, err := client.Im.Message.Reply(ctx, req)
	if err != nil {
		log.Printf("❌ 回复消息失败: %v", err)
		return
	}

	if !resp.Success() {
		log.Printf("❌ 回复消息失败: code=%d, msg=%s", resp.Code, resp.Msg)
		return
	}

	log.Printf("✅ 回复消息成功: %s", content)
}

// 创建服务台工单示例（可选功能）
func createHelpdeskTicket(ctx context.Context, client *lark.Client, description string) error {
	if HELPDESK_ID == "your_helpdesk_id" {
		log.Println("⚠️  未配置服务台 ID，跳过创建工单")
		return nil
	}

	req := larkhelpdesk.NewCreateTicketReqBuilder().
		Body(larkhelpdesk.NewCreateTicketReqBodyBuilder().
			HelpdeskId(HELPDESK_ID).
			Description(description).
			Source(1). // 1: API 创建
			Build()).
		Build()

	resp, err := client.Helpdesk.Ticket.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("创建工单失败: %w", err)
	}

	if !resp.Success() {
		return fmt.Errorf("创建工单失败: code=%d, msg=%s", resp.Code, resp.Msg)
	}

	log.Printf("✅ 工单创建成功: TicketID=%s", *resp.Data.TicketId)
	return nil
}
