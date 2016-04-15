package send

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
)

// NewSmtpSender 创建基于smtp的邮件发送实例(使用PlainAuth)
// addr 服务器地址
// from 发送者
// authPwd 验证密码
// 如果创建实例发生异常，则返回错误
func NewSmtpSender(addr string, from mail.Address, authPwd string) (Sender, error) {
	smtpCli := &SmtpClient{
		addr: addr,
		from: from,
	}
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	smtpCli.auth = smtp.PlainAuth("", from.Address, authPwd, host)
	return smtpCli, nil
}

// SmtpClient 使用smtp发送邮件
type SmtpClient struct {
	addr string
	from mail.Address
	auth smtp.Auth
}

// Send 发送邮件
func (this *SmtpClient) Send(msg *Message, isMass bool) (err error) {
	if isMass {
		err = this.massSend(msg)
	} else {
		err = this.oneSend(msg)
	}
	return
}

// AsyncSend 异步发送邮件
func (this *SmtpClient) AsyncSend(msg *Message, isMass bool, handle func(err error)) error {
	go func() {
		err := this.Send(msg, isMass)
		handle(err)
	}()
	return nil
}

// oneSend 一对一按顺序发送
func (this *SmtpClient) oneSend(msg *Message) error {
	for _, addr := range msg.To {
		header := this.getHeader(msg.Subject)
		header["To"] = addr
		if msg.Extension != nil {
			for k, v := range msg.Extension {
				header[k] = v
			}
		}
		data := this.getData(header, msg.Content)
		err := smtp.SendMail(this.addr, this.auth, this.from.Address, []string{addr}, data)
		if err != nil {
			return err
		}
	}
	return nil
}

// massSend 群发邮件
func (this *SmtpClient) massSend(msg *Message) error {
	header := this.getHeader(msg.Subject)
	if msg.Extension != nil {
		for k, v := range msg.Extension {
			header[k] = v
		}
	}
	data := this.getData(header, msg.Content)
	return smtp.SendMail(this.addr, this.auth, this.from.Address, msg.To, data)
}

func (this *SmtpClient) getHeader(subject string) map[string]string {
	header := make(map[string]string)
	header["From"] = this.from.String()
	header["Subject"] = mime.QEncoding.Encode("utf-8", subject)
	header["Mime-Version"] = "1.0"
	header["Content-Type"] = "text/html;charset=utf-8"
	header["Content-Transfer-Encoding"] = "Quoted-Printable"
	return header
}

func (this *SmtpClient) getData(header map[string]string, body io.Reader) []byte {
	buf := new(bytes.Buffer)
	for k, v := range header {
		fmt.Fprintf(buf, "%s: %s\r\n", k, v)
	}
	fmt.Fprintf(buf, "\r\n")
	io.Copy(buf, body)
	return buf.Bytes()
}
