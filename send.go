package send

import "io"

// Message 邮件发送数据
type Message struct {
	Subject   string            // 标题
	Content   io.Reader         // 支持html的消息主体
	To        []string          // 邮箱地址
	Extension map[string]string // 发送邮件消息体扩展项
}

// Sender 提供邮件发送接口
type Sender interface {
	// Send 发送邮件
	// msg 邮件发送数据
	// isMass 是否是群发,默认为一对一发送
	Send(msg *Message, isMass bool) error

	// AsyncSend 异步发送邮件
	// msg 邮件发送数据
	// isMass 是否是群发,默认为一对一发送
	// handle 发送完成之后的处理函数，如果发送失败,则返回error
	AsyncSend(msg *Message, isMass bool, handle func(err error)) error
}
