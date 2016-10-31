package sophon

import (
	"sync"
	"errors"
)

var (
	ErrEmptyHost = errors.New("sophon: invalid empty host")
)

type Queue struct {
	// 队列中的任务
	task chan Command

	// channel处理信号
	closed, cancelled, done chan struct{}

	// 计数
	count int64

	wg sync.WaitGroup
}

// 关闭队列
func (q *Queue) Close() error {
	select {
	case <-q.closed:
		return nil
	default:
		close(q.closed)

		q.task <- nil
		// 等待
		q.wg.Wait()

		// 解除阻塞
		close(q.done)
		return nil
	}
}

// 阻塞当前goroutine直到队列关闭
func (q *Queue) Block() {
	<-q.done
}

// 暂时取消
func (q *Queue) Cancel() error {
	select {
	case <-q.cancelled:
		return nil
	default:
		// 标记取消队列
		close(q.cancelled)

		return q.Close()
	}
}

