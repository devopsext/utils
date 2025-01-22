package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"time"
)

func httpContentTypeAndAuthorizationHeaders(contentType string, authorization string) map[string]string {

	headers := make(map[string]string)
	if !IsEmpty(contentType) {
		headers["Content-Type"] = contentType
	}
	if !IsEmpty(authorization) {
		headers["Authorization"] = authorization
	}
	return headers
}

func HttpRequestGetHeader(client *http.Client, URL string) (header map[string][]string, err error) {

	resp, err := client.Head(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp.Header, nil
}

func HttpRequestRawWithHeadersOutCode(client *http.Client, method, URL string, headers map[string]string, raw []byte) (body []byte, code int, err error) {

	var reader io.Reader
	if raw != nil {
		reader = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, URL, reader)
	if err != nil {
		return nil, 0, err
	}

	for k, v := range headers {
		if IsEmpty(v) {
			continue
		}
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var b []byte
	if !IsEmpty(resp.Body) {
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, resp.StatusCode, err
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return b, resp.StatusCode, fmt.Errorf(resp.Status)
	}

	return b, resp.StatusCode, nil
}

func HttpRequestRawWithHeadersOutCodeSilent(client *http.Client, method, URL string, headers map[string]string, raw []byte) (body []byte, code int, err error) {

	var reader io.Reader
	if raw != nil {
		reader = bytes.NewReader(raw)
	}

	req, err := http.NewRequest(method, URL, reader)
	if err != nil {
		return nil, 0, err
	}

	for k, v := range headers {
		if IsEmpty(v) {
			continue
		}
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var b []byte
	if !IsEmpty(resp.Body) {
		b, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, resp.StatusCode, err
		}
	}
	if len(b) == 0 {
		return []byte(fmt.Sprintf(`{"code":%d}`, resp.StatusCode)), resp.StatusCode, nil
	}

	return b, resp.StatusCode, nil
}

func HttpRequestRawWithRetry(client *http.Client, method, URL string, headers map[string]string, raw []byte, maxRetries int, retryHeader string) (body []byte, code int, err error) {

	var reader io.Reader
	if raw != nil {
		reader = bytes.NewReader(raw)
	}

	var b []byte

	req, err := http.NewRequest(method, URL, reader)
	if err != nil {
		return nil, 0, err
	}

	for k, v := range headers {
		if IsEmpty(v) {
			continue
		}
		req.Header.Set(k, v)
	}

	for attempt := 0; attempt < maxRetries; attempt++ {

		resp, err := client.Do(req)
		if err != nil {
			return nil, resp.StatusCode, err
		}
		defer resp.Body.Close()

		if !IsEmpty(resp.Body) {
			b, err = io.ReadAll(resp.Body)
			if err != nil {
				return nil, resp.StatusCode, err
			}
		}

		if resp.StatusCode == http.StatusTooManyRequests {

			retryAfter := resp.Header.Get(retryHeader)

			if retryAfter != "" {
				duration, err := time.ParseDuration(retryAfter)
				if err == nil {
					time.Sleep(duration)
					continue
				}
			}
			backoff := time.Duration(math.Pow(2, float64(attempt))) * time.Second
			time.Sleep(backoff)
			continue
		}
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return b, resp.StatusCode, fmt.Errorf(resp.Status)
		}
		return b, resp.StatusCode, nil
	}
	return nil, 0, fmt.Errorf("max retries exceeded")
}

func HttpRequestRawWithHeaders(client *http.Client, method, URL string, headers map[string]string, raw []byte) ([]byte, error) {
	b, _, err := HttpRequestRawWithHeadersOutCode(client, method, URL, headers, raw)
	return b, err
}

func HttpRequestRawWithHeadersRetry(client *http.Client, method, URL string, headers map[string]string, raw []byte, maxRetries int, retryHeader string) ([]byte, error) {
	b, _, err := HttpRequestRawWithRetry(client, method, URL, headers, raw, maxRetries, retryHeader)
	return b, err
}

func HttpPostRawWithHeaders(client *http.Client, URL string, headers map[string]string, raw []byte) ([]byte, error) {
	return HttpRequestRawWithHeaders(client, "POST", URL, headers, raw)
}

func HttpPostRawWithHeadersRetry(client *http.Client, URL string, headers map[string]string, raw []byte, maxRetries int, retryHeader string) ([]byte, error) {
	return HttpRequestRawWithHeadersRetry(client, "POST", URL, headers, raw, maxRetries, retryHeader)
}

func HttpPostRaw(client *http.Client, URL, contentType string, authorization string, raw []byte) ([]byte, error) {

	headers := httpContentTypeAndAuthorizationHeaders(contentType, authorization)
	return HttpPostRawWithHeaders(client, URL, headers, raw)
}

func HttpPostRawRetry(client *http.Client, URL, contentType string, authorization string, raw []byte, maxRetries int, retryHeader string) ([]byte, error) {

	headers := httpContentTypeAndAuthorizationHeaders(contentType, authorization)
	return HttpPostRawWithHeadersRetry(client, URL, headers, raw, maxRetries, retryHeader)
}

func HttpPostRawWithHeadersOutCode(client *http.Client, URL string, headers map[string]string, raw []byte) (body []byte, code int, err error) {
	return HttpRequestRawWithHeadersOutCode(client, "POST", URL, headers, raw)
}

func HttpPostRawOutCode(client *http.Client, URL, contentType string, authorization string, raw []byte) (body []byte, code int, err error) {

	headers := httpContentTypeAndAuthorizationHeaders(contentType, authorization)
	return HttpPostRawWithHeadersOutCode(client, URL, headers, raw)
}

func HttpPutRawWithHeaders(client *http.Client, URL string, headers map[string]string, raw []byte) ([]byte, error) {
	return HttpRequestRawWithHeaders(client, "PUT", URL, headers, raw)
}

func HttpPutRaw(client *http.Client, URL, contentType string, authorization string, raw []byte) ([]byte, error) {

	headers := httpContentTypeAndAuthorizationHeaders(contentType, authorization)
	return HttpPutRawWithHeaders(client, URL, headers, raw)
}

func HttpDeleteRawWithHeaders(client *http.Client, URL string, headers map[string]string, raw []byte) ([]byte, error) {
	return HttpRequestRawWithHeaders(client, "DELETE", URL, headers, raw)
}

func HttpDeleteRaw(client *http.Client, URL, contentType string, authorization string, raw []byte) ([]byte, error) {

	headers := httpContentTypeAndAuthorizationHeaders(contentType, authorization)
	return HttpDeleteRawWithHeaders(client, URL, headers, raw)
}

func HttpGetRawWithHeaders(client *http.Client, URL string, headers map[string]string) ([]byte, error) {
	return HttpRequestRawWithHeaders(client, "GET", URL, headers, nil)
}

func HttpGetRawWithHeadersRetry(client *http.Client, URL string, headers map[string]string, maxRetries int, retryHeader string) ([]byte, error) {
	return HttpRequestRawWithHeadersRetry(client, "GET", URL, headers, nil, maxRetries, retryHeader)
}

func HttpGetRaw(client *http.Client, URL, contentType string, authorization string) ([]byte, error) {
	headers := httpContentTypeAndAuthorizationHeaders(contentType, authorization)
	return HttpGetRawWithHeaders(client, URL, headers)
}

func HttpGetRawRetry(client *http.Client, URL, contentType string, authorization string, maxRetries int, retryHeader string) ([]byte, error) {
	headers := httpContentTypeAndAuthorizationHeaders(contentType, authorization)
	return HttpGetRawWithHeadersRetry(client, URL, headers, maxRetries, retryHeader)
}

func HttpGetHeader(client *http.Client, URL string) (map[string][]string, error) {
	return HttpRequestGetHeader(client, URL)
}

func NewHttpClient(timeout int, insecure bool) *http.Client {

	var transport = &http.Transport{
		Dial:                (&net.Dialer{Timeout: time.Duration(timeout) * time.Second}).Dial,
		TLSHandshakeTimeout: time.Duration(timeout) * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: insecure},
	}

	var client = &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: transport,
	}
	return client
}

func NewHttpInsecureClient(timeout int) *http.Client {
	return NewHttpClient(timeout, true)
}

func NewHttpSecureClient(timeout int) *http.Client {
	return NewHttpClient(timeout, false)
}
