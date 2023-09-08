package httpclient

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type HttpGetter interface {
	GetFileBytes(ctx context.Context, fileURL string) ([]byte, error)
	GetAndSaveFile(ctx context.Context, fileURL string, activityID string, counter int) (string, error)
}

type httpGetter struct {
	client *http.Client
}

func NewGetter() HttpGetter {
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		timeout = 30
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	return newGetter(client)
}

func newGetter(client *http.Client) HttpGetter {
	return &httpGetter{client: client}
}

func (g *httpGetter) GetFileBytes(ctx context.Context, fileURL string) ([]byte, error) {
	var respBytes []byte
	var err error

	resp, err := http.Get(fileURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		respBytes, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
	}

	return respBytes, nil
}

func (g *httpGetter) GetAndSaveFile(ctx context.Context, fileURL string, activityID string, counter int) (string, error) {
	var err error

	resp, err := http.Get(fileURL)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	parsedFileURL, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}

	path := parsedFileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	fullFileName := "/tmp/" + activityID + "-image" + strconv.Itoa(counter) + "-" + fileName

	file, err := os.Create(fullFileName)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	defer file.Close()

	return fullFileName, nil
}
