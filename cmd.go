package sophon

import (
	"net/url"
)

// Command 接口定义了访问需要请求资源的方法
type Command interface {
	URL() *url.URL
	Method() string
}

// 实现一个基本的command
type Cmd struct {
	url    *url.URL
	method string
}

func (cmd *Cmd) URL() *url.URL {
	return cmd.url
}

func (cmd *Cmd) Method() string {
	return cmd.method
}
