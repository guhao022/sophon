package sophon

import "sync"

type Queue struct {
	// 队列中的任务
	task chan interface{}

	// channel处理信号
	closed, cancelled, done chan struct{}

	// 计数
	count int64

	wg sync.WaitGroup
}

func NewQueue(task interface{}) *Queue {
	q := new(Queue)
	q.task <- task
	q.closed = make(chan struct{})
	q.cancelled = make(chan struct{})
	q.done = make(chan struct{})

	return q
}

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


