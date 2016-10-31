// 负责抓取网页内容
package sophon

import (
	"net/http"
	"time"
	"sync"
	"spider/core/downloader/surfer/agent"
)

type Fetcher struct {
	hosts map[string]chan Command

	UserAgent string

	// 空闲中的goroutine等待时间
	WorkerIdleTTL time.Duration

	Client *http.Client

	CrawlDelay time.Duration

	w *Worker
	mu sync.Mutex
}

func (f *Fetcher) Init() *Worker {
	f.hosts = make(map[string]chan Command)

	f.w = &Worker{
		cmd:        make(chan Command, 1),
		closed:    make(chan byte),
		cancelled: make(chan byte),
		done:      make(chan byte),
	}

	// Start the one and only queue processing goroutine.
	f.w.wg.Add(1)
	go f.processQueue()

	return f.w
}

func (f *Fetcher) processQueue() {
LOOP:
	for v := range f.w.cmd {
		if v == nil {
			select {
			case <-f.w.closed:
				break LOOP
			default:
			// Keep going
			}
		}
		select {
		case <-f.w.cancelled:
			continue
		default:
		// go on
		}

		u := v.URL()

		f.mu.Lock()
		in, ok := f.hosts[u.Host]
		if !ok {

			// Create the infinite queue: the in channel to send on, and the out channel
			// to read from in the host's goroutine, and add to the hosts map
			var out chan Command
			in, out = make(chan Command, 1), make(chan Command, 1)
			f.hosts[u.Host] = in
			f.mu.Unlock()
			f.w.wg.Add(1)
			// Start the infinite queue goroutine for this host
			go sliceIQ(in, out)
			// Start the working goroutine for this host
			go f.processChan(out, u.Host)

		} else {
			f.mu.Unlock()
		}
		// Send the request
		in <- v
	}

	// Close all host channels now that it is impossible to send on those. Those are the `in`
	// channels of the infinite queue. It will then drain any pending events, triggering the
	// handlers for each in the worker goro, and then the infinite queue goro will terminate
	// and close the `out` channel, which in turn will terminate the worker goro.
	f.mu.Lock()
	for _, ch := range f.hosts {
		close(ch)
	}
	f.hosts = make(map[string]chan Command)
	f.mu.Unlock()

	f.w.wg.Done()
}

// Goroutine for a host's worker, processing requests for all its URLs.
func (f *Fetcher) processChan(ch <-chan Command, hostKey string) {
	var (
		wait  <-chan time.Time
		ttl   <-chan time.Time
		delay = f.CrawlDelay
	)

	loop:
	for {
		select {
		case <-f.q.cancelled:
			break loop
		case v, ok := <-ch:
			if !ok {
				// Terminate this goroutine, channel is closed
				break loop
			}

		// Wait for the prescribed delay
			if wait != nil {
				<-wait
			}

		// was it cancelled during the wait? check again
				select {
				case <-f.q.cancelled:
					break loop
				default:
				// go on
				}

			switch r, ok := v.(robotCommand); {
			case ok:
				// This is the robots.txt request
				agent = f.getRobotAgent(r)
				// Initialize the crawl delay
				if agent != nil && agent.CrawlDelay > 0 {
					delay = agent.CrawlDelay
				}
				wait = time.After(delay)

			case agent == nil || agent.Test(v.URL().Path):
				// Path allowed, process the request
				res, err := f.doRequest(v)
				f.visit(v, res, err)
				// No delay on error - the remote host was not reached
				if err == nil {
					wait = time.After(delay)
				} else {
					wait = nil
				}

			default:
				// Path disallowed by robots.txt
				f.visit(v, nil, ErrDisallowed)
				wait = nil
			}
		// Every time a command is received, reset the ttl channel
			ttl = time.After(f.WorkerIdleTTL)

		case <-ttl:
		// Worker has been idle for WorkerIdleTTL, terminate it
			f.mu.Lock()
			inch, ok := f.hosts[hostKey]
			delete(f.hosts, hostKey)

		// Close the queue if AutoClose is set and there are no more hosts.
			if f.AutoClose && len(f.hosts) == 0 {
				go f.q.Close()
			}
			f.mu.Unlock()
			if ok {
				close(inch)
			}
			break loop
		}
	}

	f.q.wg.Done()
}
