// 负责抓取网页内容
package sophon

import (
	"net/http"
	"time"
)

type Fetch interface {
	Do(*http.Request) (*http.Response, error)
}

type Fetcher struct {
	hosts map[string]chan Command

	UserAgent string

	// 空闲中的goroutine等待时间
	WorkerIdleTTL time.Duration
}

//type qu
