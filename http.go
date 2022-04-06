package utils

import (
	"crypto/tls"
	"net"
	"net/http"
	"time"
)

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
