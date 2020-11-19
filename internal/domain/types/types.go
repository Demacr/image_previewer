package types

import (
	"io"
	"net/http"
)

type DownloadedImage struct {
	Buffer io.ReadCloser
	Header http.Header
}
