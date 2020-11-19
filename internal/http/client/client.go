package client

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	domain "github.com/Demacr/image_previewer/internal/domain/types"
	"github.com/pkg/errors"
)

func GetImage(url string) (result *domain.DownloadedImage, err error) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "GET", "http://"+url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error during creating request")
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error during getting image")
	}
	defer resp.Body.Close()

	buffer := &bytes.Buffer{}
	if _, err := io.Copy(buffer, resp.Body); err != nil {
		return nil, errors.Wrap(err, "error during reading body message")
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "image/") {
		return nil, errors.New("wrong Content-Type")
	}

	result = &domain.DownloadedImage{
		Buffer: ioutil.NopCloser(buffer),
		Header: resp.Header,
	}

	return result, nil
}
