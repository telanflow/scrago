package downloader

import (
	"net/http"
)

type Option func(downloader Downloader)

type Downloader interface {
	Do(request *http.Request) (*http.Response, error)
	UseOptions(options ...Option)
}

