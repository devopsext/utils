package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
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

func HttpRequestRawWithHeaders(client *http.Client, method, URL string, headers map[string]string, raw []byte) ([]byte, error) {
	b, _, err := HttpRequestRawWithHeadersOutCode(client, method, URL, headers, raw)
	return b, err
}

func HttpPostRawWithHeaders(client *http.Client, URL string, headers map[string]string, raw []byte) ([]byte, error) {
	return HttpRequestRawWithHeaders(client, "POST", URL, headers, raw)
}

func HttpPostRaw(client *http.Client, URL, contentType string, authorization string, raw []byte) ([]byte, error) {

	headers := httpContentTypeAndAuthorizationHeaders(contentType, authorization)
	return HttpPostRawWithHeaders(client, URL, headers, raw)
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

func HttpGetRaw(client *http.Client, URL, contentType string, authorization string) ([]byte, error) {
	headers := httpContentTypeAndAuthorizationHeaders(contentType, authorization)
	return HttpGetRawWithHeaders(client, URL, headers)
}

func HttpGetHeader(client *http.Client, URL string) (map[string][]string, error) {
	return HttpRequestGetHeaderCode(client, URL)
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
