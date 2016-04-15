# 基于golang的smtp包实现的邮件发送

[![GoDoc](https://godoc.org/github.com/go-mailer/send?status.svg)](https://godoc.org/github.com/go-mailer/send)

> 支持异步发送邮件、一对一发送、群发

## 获取

``` bash
$ go get -v github.com/go-mailer/send
```

## 使用范例(异步发送)

``` go
package main

import (
	"bytes"
	"fmt"
	"net/mail"
	"sync"

	"github.com/go-mailer/send"
)

func main() {
	from := mail.Address{Name: "Lyric", Address: "nianshou.tian@zjelite.com"}
	sender, err := send.NewSmtpSender("smtp.exmail.qq.com:25", from, "xxx")
	if err != nil {
		panic(err)
	}
	var wg sync.WaitGroup
	wg.Add(1)
	msg := &send.Message{
		Subject: "异步发送邮件测试",
		Content: bytes.NewBufferString("<h1>你好，异步测试邮件内容</h1>"),
		To:      []string{"tiannianshou@126.com"},
	}
	err = sender.AsyncSend(msg, false, func(err error) {
		defer wg.Done()
		if err != nil {
			fmt.Println("发送邮件出现错误：", err)
		}
	})
	if err != nil {
		panic(err)
	}
	wg.Wait()
	fmt.Println("邮件发送完成")
}
```

## License

	Copyright 2016.All rights reserved.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

