package sophon

import (
	"sync"
	"errors"
	"net/url"
)

var (
	ErrEmptyHost = errors.New("sophon: invalid empty host")

	ErrQueueClosed = errors.New("sophon: send on a closed queue")
)

type Woker struct {
	// 需要处理的网站基本信息
	cmd chan Command

	// woker的数量
	count int64

	// channel 信号标记
	closed, canneled, done chan byte

	wg sync.WaitGroup
}

// 添加任务
func (w *Woker) Send(c Command) error {
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
func (w *Woker) sendWithMethod(method string, rawurl []string) (int, error) {
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

// 关闭队列
func (w *Woker) Close() error {
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

// 阻塞当前goroutine直到队列关闭
func (w *Woker) Block() {
	<-w.done
}

// 暂时取消
func (w *Woker) Cancel() error {
	select {
	case <-w.cancelled:
		return nil
	default:
	// 标记取消队列
		close(q.cancelled)

		return q.Close()
	}
}
