package utils

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

const defaultHTTPTimeout = 20 * time.Second

func Request(url, method string, headers http.Header, body io.Reader) ([]byte, error) {
	client := http.Client{
		Timeout: defaultHTTPTimeout,
	}

	request, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, errors.Wrap(err, "request creation failed")
	}

	if headers != nil {
		request.Header = headers
	}
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response status code is not 200, code: %d", resp.StatusCode)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func Get(url string, headers http.Header, body io.Reader) ([]byte, error) {
	return Request(url, http.MethodGet, headers, body)
}

func Post(url string, headers http.Header, body io.Reader) ([]byte, error) {
	return Request(url, http.MethodPost, headers, body)
}
