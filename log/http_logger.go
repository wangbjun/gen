package log

import (
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
	Logger.Sugar().Infof("Request_start method=%s url=%s", request.Method, request.URL.String())

	startTime := time.Now()

	response, err := c.rt.RoundTrip(request)

	duration := time.Since(startTime)
	duration /= time.Millisecond

	if err != nil {
		Logger.Sugar().Errorf("Response_error method=%s duration=%d url=%s error=%s",
			request.Method, duration, request.URL.String(), err.Error())
	} else {
		Logger.Sugar().Infof("Response_success method=%s status=%d duration=%d url=%s",
			request.Method, response.StatusCode, duration, request.URL.String())
	}

	return response, err
}
