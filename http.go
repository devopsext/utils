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

func HttpRequestRawWithHeaders(client *http.Client, method, URL string, headers map[string]string, raw []byte) ([]byte, error) {

	reader := bytes.NewReader(raw)

	req, err := http.NewRequest(method, URL, reader)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		if IsEmpty(v) {
			continue
		}
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("response status code: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func HttpPostRawWithHeaders(client *http.Client, URL string, headers map[string]string, raw []byte) ([]byte, error) {
	return HttpRequestRawWithHeaders(client, "POST", URL, headers, raw)
}

func HttpPostRaw(client *http.Client, URL, contentType string, authorization string, raw []byte) ([]byte, error) {

	headers := make(map[string]string)
	if !IsEmpty(contentType) {
		headers["Content-Type"] = contentType
	}
	if !IsEmpty(authorization) {
		headers["Authorization"] = authorization
	}
	return HttpPostRawWithHeaders(client, URL, headers, raw)
}

func HttpRequestRawWithHeadersOutCode(client *http.Client, method, URL string, headers map[string]string, raw []byte) (body []byte, code int, err error) {

	reader := bytes.NewReader(raw)

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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, 0, fmt.Errorf(resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return b, resp.StatusCode, nil
}

func HttpPostRawWithHeadersOutCode(client *http.Client, URL string, headers map[string]string, raw []byte) (body []byte, code int, err error) {
	return HttpRequestRawWithHeadersOutCode(client, "POST", URL, headers, raw)
}

func HttpPostRawOutCode(client *http.Client, URL, contentType string, authorization string, raw []byte) (body []byte, code int, err error) {

	headers := make(map[string]string)
	if !IsEmpty(contentType) {
		headers["Content-Type"] = contentType
	}
	if !IsEmpty(authorization) {
		headers["Authorization"] = authorization
	}
	return HttpPostRawWithHeadersOutCode(client, URL, headers, raw)
}

func HttpPutRawWithHeaders(client *http.Client, URL string, headers map[string]string, raw []byte) ([]byte, error) {
	return HttpRequestRawWithHeaders(client, "PUT", URL, headers, raw)
}

func HttpPutRaw(client *http.Client, URL, contentType string, authorization string, raw []byte) ([]byte, error) {

	headers := make(map[string]string)
	if !IsEmpty(contentType) {
		headers["Content-Type"] = contentType
	}
	if !IsEmpty(authorization) {
		headers["Authorization"] = authorization
	}
	return HttpPutRawWithHeaders(client, URL, headers, raw)
}

func HttpGetRawWithHeaders(client *http.Client, URL string, headers map[string]string) ([]byte, error) {

	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		if IsEmpty(v) {
			continue
		}
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf(resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func HttpGetRaw(client *http.Client, URL, contentType string, authorization string) ([]byte, error) {

	headers := make(map[string]string)
	if !IsEmpty(contentType) {
		headers["Content-Type"] = contentType
	}
	if !IsEmpty(authorization) {
		headers["Authorization"] = authorization
	}

	return HttpGetRawWithHeaders(client, URL, headers)
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
