package request

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ResponseData struct {
	Size   int64
	Timing int64
	Err    error
}

type HttpClient struct {
	Client     *http.Client
	RetryCount int
}

func New(timeout int64, retry int) HttpClient {
	return HttpClient{
		Client: &http.Client{
			Timeout: time.Duration(timeout * int64(time.Millisecond)),
			Transport: &http.Transport{
				TLSHandshakeTimeout:   time.Second,
				ResponseHeaderTimeout: time.Second,
			},
		},
		RetryCount: retry,
	}
}

func (hc HttpClient) GetRequestInfo(checkUrl string) ResponseData {
	_, err := url.ParseRequestURI(checkUrl)
	if err != nil {
		return ResponseData{
			Err: err,
		}
	}

	response, timing, err := hc.request(checkUrl)
	if err != nil {
		return ResponseData{
			Err: err,
		}
	}

	size, err := io.Copy(io.Discard, response.Body)
	if err != nil {
		return ResponseData{
			Err: err,
		}
	}

	return ResponseData{
		Size:   size,
		Timing: timing.Milliseconds(),
	}
}

func (hc HttpClient) request(url string) (*http.Response, time.Duration, error) {
	for i := 0; i < hc.RetryCount; i++ {
		timeStart := time.Now()
		response, err := hc.Client.Get(url)
		timing := time.Since(timeStart)

		if err != nil {
			if os.IsTimeout(err) {
				continue
			}

			return nil, time.Duration(0), fmt.Errorf("Request error: %v", err)
		}

		return response, timing, nil
	}

	return nil, time.Duration(0), fmt.Errorf("Request timeout")
}
