package sophon

import (
	"net/url"
	"io"
	"net/http"
)

// Command 接口定义了访问需要请求资源的方法
type Command interface {
	URL() *url.URL
	Method() string
}

// BasicAuthProvider 接口获取凭证使用基本认证执行请求
type BasicAuthProvider interface {
	BasicAuth() (user string, pwd string)
}

// ReaderProvider 接口以io.Reader作为请求体，他比ValuesProvider接口的优先级更高，
// 所以，如果两个接口都实现了，则使用此接口
type ReaderProvider interface {
	Reader() io.Reader
}

// ValuesProvider 接口以url.Values作为请求主体，它比ReaderProvider接口的优先级低，
// 所以，如果两个接口都实现了，则使用ReaderProvider接口。
// 如果请求中没有明确设置 Content-Type , 则会自动设置为 "application/x-www-form-urlencoded".
type ValuesProvider interface {
	Values() url.Values
}

// CookiesProvider接口 获取cookies并在发送请求的时候一起发送
type CookiesProvider interface {
	Cookies() []*http.Cookie
}

// HeaderProvider 接口 获取头信息
type HeaderProvider interface {
	Header() http.Header
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


