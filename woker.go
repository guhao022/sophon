// Worker 是用于辅助Fetch处理超链接的一个队列
package sophon

import (
	"sync"
	"errors"
	"net/url"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"path"
	"log"
)

var (
	ErrEmptyHost = errors.New("sophon: invalid empty host")

	ErrQueueClosed = errors.New("sophon: send on a closed queue")
)

type Worker struct {
	// 需要处理的网站基本信息
	cmd chan Command

	count int64

	// channel 信号标记
	closed, cancelled, done chan byte

	wg sync.WaitGroup
}

// 添加任务
func (w *Worker) Send(c Command) error {
	if c == nil {
		return ErrEmptyHost
	}
	if u := c.URL(); u == nil || u.Host == "" {
		return ErrEmptyHost
	}

	select {
	case <-w.closed:
		return ErrQueueClosed
	default:
		w.cmd <- c
	}

	return nil
}

// 发送带method方法的任务
func (w *Worker) sendWithMethod(method string, rawurl []string) (int, error) {
	for i, v := range rawurl {
		u, err := url.Parse(v)
		if err != nil {
			return i, err
		}
		if err := w.Send(&Cmd{url: u, method: method}); err != nil {
			return i, err
		}
	}
	return len(rawurl), nil
}

func (w *Worker) SendStringGet(rawurl ...string) (int, error) {
	return w.sendWithMethod("GET", rawurl)
}

func (w *Worker) SendStringHead(rawurl ...string) (int, error) {
	return w.sendWithMethod("HEAD", rawurl)
}

// 关闭worker
func (w *Worker) Close() error {
	select {
	case <-w.closed:
		return nil
	default:
		close(w.closed)

		w.cmd <- nil
	// 等待
		w.wg.Wait()

	// 解除阻塞
		close(w.done)
		return nil
	}
}

// 阻塞当前goroutine直到worker关闭
func (w *Worker) Block() {
	<-w.done
}

// 取消
func (w *Worker) Cancel() error {
	select {
	case <-w.cancelled:
		return nil
	default:
		// 标记取消队列
		close(w.cancelled)

		return w.Close()
	}
}

func (w *Worker) processLinks(doc *goquery.Document) (result []*url.URL) {
	baseURL, _ := doc.Find("base[href]").Attr("href")
	urls := doc.Find("a[href]").Map(func(_ int, s *goquery.Selection) string {
		val, _ := s.Attr("href")
		if baseURL != "" {
			val = handleBaseTag(doc.Url, baseURL, val)
		}
		return val
	})
	for _, s := range urls {
		// If href starts with "#", then it points to this same exact URL, ignore (will fail to parse anyway)
		if len(s) > 0 && !strings.HasPrefix(s, "#") {
			if parsed, e := url.Parse(s); e == nil {
				parsed = doc.Url.ResolveReference(parsed)
				result = append(result, parsed)
			} else {
				log.Printf("ignore on unparsable policy %s: %s", s, e.Error())
			}
		}
	}
	return
}

func handleBaseTag(root *url.URL, baseHref string, aHref string) string {
	resolvedBase, err := root.Parse(baseHref)
	if err != nil {
		return ""
	}

	parsedURL, err := url.Parse(aHref)
	if err != nil {
		return ""
	}
	// If a[href] starts with a /, it overrides the base[href]
	if parsedURL.Host == "" && !strings.HasPrefix(aHref, "/") {
		aHref = path.Join(resolvedBase.Path, aHref)
	}

	resolvedURL, err := resolvedBase.Parse(aHref)
	if err != nil {
		return ""
	}
	return resolvedURL.String()
}
