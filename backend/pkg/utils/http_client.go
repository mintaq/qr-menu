package utils

import (
	"net/http"
	"time"
)

var HttpClient *http.Client

func CreateHttpClient() {
	HttpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
}
