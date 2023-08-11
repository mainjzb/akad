package main

import (
	"crypto/tls"
	"io"
	"net/http"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

func Test(t *testing.T) {
	client := &http.Client{
		Transport: &http.Transport{
			DialContext: nil,
			// DialTLSContext:
			TLSClientConfig:   &tls.Config{ServerName: "http2.akamai.com"},
			ForceAttemptHTTP2: true,
		},
	}
	req, err := http.NewRequest("GET", "https://http2.akamai.com/demo/tile-16.png", nil)
	if err != nil {
		log.Error(err)
	}
	req.Host = "http2.akamai.com"

	response, err := client.Do(req)
	_ = response

	io.ReadAll(response.Body)
	time.Sleep(1 * time.Second)
}
