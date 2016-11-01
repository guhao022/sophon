// 负责抓取网页内容
package sophon

import (
	"net/http"
	"time"
	"sync"
	"io"
	"io/ioutil"
	"strings"
)

const (
	DefaultCrawlDelay = 5 * time.Second

	DefaultUserAgent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.71 Safari/537.36"

	DefaultWorkerIdleTTL = 30 * time.Second
)


type Fetcher struct {
	// 为每个请求调用处理程序. 所有成功排队请求产生一个处理程序。
	Handler Handler

	UserAgent string

	// 空闲中的goroutine等待时间
	WorkerIdleTTL time.Duration

	Client *http.Client

	CrawlDelay time.Duration

	AutoClose	bool

	w *Worker
	mu sync.Mutex

	hosts map[string]chan Command
}

func New(h Handler) *Fetcher {
	return &Fetcher{
		Handler:       h,
		CrawlDelay:    DefaultCrawlDelay,
		Client:    http.DefaultClient,
		UserAgent:     DefaultUserAgent,
		WorkerIdleTTL: DefaultWorkerIdleTTL,
	}
}

func (f *Fetcher) Start() *Worker {
	f.hosts = make(map[string]chan Command)

	f.w = &Worker{
		cmd:        make(chan Command, 1),
		closed:    make(chan byte),
		cancelled: make(chan byte),
		done:      make(chan byte),
	}

	// Start the one and only queue processing goroutine.
	f.w.wg.Add(1)
	go f.processWorker()

	return f.w
}

// 让所有的worker运行在自己的goroutine中
func (f *Fetcher) processWorker() {
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
		case <-f.w.cancelled:
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
				case <-f.w.cancelled:
					break loop
				default:
				// go on
				}

			res, err := f.doRequest(v)
			f.visit(v, res, err)
		// No delay on error - the remote host was not reached
			if err == nil {
				wait = time.After(delay)
			} else {
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
				go f.w.Close()
			}
			f.mu.Unlock()
			if ok {
				close(inch)
			}
			break loop
		}
	}

	f.w.wg.Done()
}

func (f *Fetcher) visit(cmd Command, res *http.Response, err error) {
	if res != nil && res.Body != nil {
		defer res.Body.Close()
	}
	// if the Command implements Handler, call that handler, otherwise
	// dispatch to the Fetcher's Handler.
	if h, ok := cmd.(Handler); ok {
		h.Handle(&Context{Cmd: cmd, W: f.w}, res, err)
		return
	}
	f.Handler.Handle(&Context{Cmd: cmd, W: f.w}, res, err)
}

//
func (f *Fetcher) doRequest(cmd Command) (*http.Response, error) {
	req, err := http.NewRequest(cmd.Method(), cmd.URL().String(), nil)
	if err != nil {
		return nil, err
	}
	// If the Command implements some other recognized interfaces, set
	// the request accordingly (see cmd.go for the list of interfaces).
	// First, the Header values.
	if hd, ok := cmd.(HeaderProvider); ok {
		for k, v := range hd.Header() {
			req.Header[k] = v
		}
	}
	// BasicAuth has higher priority than an Authorization header set by
	// a HeaderProvider.
	if ba, ok := cmd.(BasicAuthProvider); ok {
		req.SetBasicAuth(ba.BasicAuth())
	}
	// Cookies are added to the request, even if some cookies were set
	// by a HeaderProvider.
	if ck, ok := cmd.(CookiesProvider); ok {
		for _, c := range ck.Cookies() {
			req.AddCookie(c)
		}
	}
	// For the body of the request, ReaderProvider has higher priority
	// than ValuesProvider.
	if rd, ok := cmd.(ReaderProvider); ok {
		rdr := rd.Reader()
		rc, ok := rdr.(io.ReadCloser)
		if !ok {
			rc = ioutil.NopCloser(rdr)
		}
		req.Body = rc
	} else if val, ok := cmd.(ValuesProvider); ok {
		v := val.Values()
		req.Body = ioutil.NopCloser(strings.NewReader(v.Encode()))
		if req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		}
	}

	if req.Header.Get("User-Agent") == "" {
		req.Header.Set("User-Agent", f.UserAgent)
	}
	// Do the request.
	res, err := f.Client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

