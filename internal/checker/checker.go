package checker

import (
	"context"
	"net/http"
	"sync"
	"time"
)

type Status string

const (
	StatusHealthy   Status = "healthy"
	StatusUnhealthy Status = "unhealthy"
)

type EndpointStatus struct {
	URL        string `json:"url"`
	Status     Status `json:"status"`
	StatusCode int    `json:"status_code"`
	Latency    int64  `json:"latency_ms"`
	Error      string `json:"error,omitempty"`
}

type Checker struct {
	endpoints []string
	client    *http.Client
	timeout   time.Duration
}

func NewChecker(endpoints []string, timeout time.Duration) *Checker {
	return &Checker{
		endpoints: endpoints,
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
	}
}

func (c *Checker) CheckAll(ctx context.Context) []EndpointStatus {
	results := make([]EndpointStatus, len(c.endpoints))
	var wg sync.WaitGroup

	for i, endpoint := range c.endpoints {
		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()
			results[idx] = c.Check(ctx, url)
		}(i, endpoint)
	}

	wg.Wait()
	return results
}

func (c *Checker) Check(ctx context.Context, url string) EndpointStatus {
	status := EndpointStatus{
		URL:    url,
		Status: StatusUnhealthy,
	}

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		status.Error = err.Error()
		return status
	}

	resp, err := c.client.Do(req)
	status.Latency = time.Since(start).Milliseconds()

	if err != nil {
		status.Error = err.Error()
		return status
	}
	defer resp.Body.Close()

	status.StatusCode = resp.StatusCode

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		status.Status = StatusHealthy
	}

	return status
}
