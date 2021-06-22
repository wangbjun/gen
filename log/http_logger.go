package log

import (
	"fmt"
	"net/http"
	"time"
)

type loggedRoundTripper struct {
	rt http.RoundTripper
}

// GetHttpLoggerTransport 在默认http client基础上增加日志功能
func GetHttpLoggerTransport() *loggedRoundTripper {
	return &loggedRoundTripper{http.DefaultTransport}
}

func (c *loggedRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	Info(fmt.Sprintf("Request_start method=%s url=%s", request.Method, request.URL.String()))

	startTime := time.Now()
	response, err := c.rt.RoundTrip(request)

	duration := time.Since(startTime)
	duration /= time.Millisecond
	if err != nil {
		Error(fmt.Sprintf("Response_error method=%s duration=%d url=%s error=%s",
			request.Method, duration, request.URL.String(), err.Error()))
	} else {
		Info(fmt.Sprintf("Response_success method=%s status=%d duration=%d url=%s",
			request.Method, response.StatusCode, duration, request.URL.String()))
	}
	return response, err
}
