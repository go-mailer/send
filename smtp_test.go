package send

import (
	"bytes"
	"net/mail"
	"sync"
	"testing"
)

var (
	Host    = "smtp.exmail.qq.com:25"
	From    = mail.Address{Name: "Lyric", Address: "nianshou.tian@zjelite.com"}
	FromPwd = "xxx"
)

func TestSend(t *testing.T) {
	sender, err := NewSmtpSender(Host, From, FromPwd)
	if err != nil {
		t.Error(err)
		return
	}
	msg := &Message{
		Subject: "同步发送邮件测试",
		Content: bytes.NewBufferString("<h1>你好，同步测试邮件内容</h1>"),
		To:      []string{"tiannianshou@126.com"},
	}
	err = sender.Send(msg, false)
	if err != nil {
		t.Error(err)
	}
}

func TestAsyncSend(t *testing.T) {
	sender, err := NewSmtpSender(Host, From, FromPwd)
	if err != nil {
		t.Error(err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	msg := &Message{
		Subject: "异步发送邮件测试",
		Content: bytes.NewBufferString("<h1>你好，异步测试邮件内容</h1>"),
		To:      []string{"tiannianshou@126.com"},
	}
	err = sender.AsyncSend(msg, false, func(err error) {
		defer wg.Done()
		if err != nil {
			t.Error(err)
		}
	})
	if err != nil {
		t.Error(err)
	}
	wg.Wait()
}
