package httpclient

import (
	"bytes"
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"go-poc/utils/log"
)

type RetriableError struct {
	message string
}

type HttpRequestPayload struct {
	Method string `json:"method"`
	URL    string `json:"url"`
	Body   string `json:"body,omitempty"`
}

type HttpResponsePayload struct {
	StatusCode int    `json:"status_code"`
	Body       string `json:"body,omitempty"`
}

type HttpDoer interface {
	Do(ctx context.Context, req *http.Request) ([]byte, int, error)
	SetTimeout(duration string) error
}

type ProxiedHttpDoer interface {
	HttpDoer
	WithProxyAuthHeader(value string) (ProxiedHttpDoer, error)
}

type httpDoer struct {
	client *http.Client
}

type proxiedHttpDoer struct {
	httpDoer
}

func NewDoer() HttpDoer {
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		timeout = 30
	}

	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	return newDoer(client)
}

func NewDoerWithProxy(proxyUrl *url.URL) HttpDoer {
	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		timeout = 30
	}

	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)},
	}
	return newProxiedDoer(client)
}

func newDoer(client *http.Client) HttpDoer {
	return &httpDoer{client: client}
}

func newProxiedDoer(client *http.Client) ProxiedHttpDoer {
	return &proxiedHttpDoer{httpDoer{client: client}}
}

func (d *proxiedHttpDoer) WithProxyAuthHeader(value string) (ProxiedHttpDoer, error) {
	oldTransport, _ := d.httpDoer.client.Transport.(*http.Transport)
	// copy the transport object so it doesn't use same header
	newTransport := *oldTransport

	header := http.Header{}
	header.Add("Proxy-Authorization", "Basic "+value)
	newTransport.ProxyConnectHeader = header

	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		timeout = 30
	}

	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: &newTransport,
	}

	return newProxiedDoer(client), nil
}

// SetTimeout sets timeout tolerance for the next http client call via Do()
func (d *httpDoer) SetTimeout(duration string) error {
	timeoutDuration, err := time.ParseDuration(duration)
	if err != nil {
		return err
	}
	d.client.Timeout = timeoutDuration
	return nil
}

// The main logic of *httpDoer.Do()
func (d *httpDoer) do(ctx context.Context, req *http.Request) ([]byte, int, error, HttpRequestPayload, HttpResponsePayload) {
	reqObj := HttpRequestPayload{
		Method: req.Method,
		URL:    req.URL.String(),
	}
	resObj := HttpResponsePayload{}

	truncatedBodyBytes := truncateBytes(copyBodyBytes(req))
	if truncatedBodyBytes != nil && len(truncatedBodyBytes) > 1 {
		reqObj.Body = string(truncatedBodyBytes)
		log.WithContext(ctx).Debugf("%s %s\n%s", req.Method, req.URL, string(truncatedBodyBytes))
	} else {
		log.WithContext(ctx).Debugf("%s %s", req.Method, req.URL)
	}

	resp, err := d.client.Do(req)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			resObj.StatusCode = 500
			log.WithContext(ctx).Errorf("%s %s %s", req.Method, req.URL, newRetriableError(err.Error()))
			return nil, 500, err, reqObj, resObj
		}

		log.WithContext(ctx).Errorf("%s %s %s", req.Method, req.URL, newRetriableError(err.Error()))
		return nil, 0, err, reqObj, resObj
	}
	defer resp.Body.Close()

	resObj.StatusCode = resp.StatusCode
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithContext(ctx).Errorf("%s %s %s", req.Method, req.URL, newRetriableError(err.Error()))
		return nil, resp.StatusCode, err, reqObj, resObj
	}

	truncatedRespBytes := truncateBytes(respBytes)
	if len(truncatedRespBytes) > 1 {
		resObj.Body = string(truncatedRespBytes)
		log.WithContext(ctx).Debugf("%s\n%s", resp.Status, string(truncatedRespBytes))
	} else {
		log.WithContext(ctx).Debugf("%s", resp.Status)
	}

	req.Body = ioutil.NopCloser(bytes.NewBuffer(truncatedBodyBytes))
	return respBytes, resp.StatusCode, err, reqObj, resObj
}

func (d *httpDoer) Do(ctx context.Context, req *http.Request) ([]byte, int, error) {
	respBytes, statusCode, err, _, _ := d.do(ctx, req)

	return respBytes, statusCode, err
}

func copyBodyBytes(req *http.Request) []byte {
	if req.Body == nil {
		return nil
	}
	bodyBytes, _ := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	return bodyBytes
}

func truncateBytes(bytes []byte) []byte {
	if bytes == nil || len(bytes) < log.MAX_LOG_ENTRY_SIZE {
		return bytes
	}
	return bytes[0:log.MAX_LOG_ENTRY_SIZE]
}

func newRetriableError(message string) *RetriableError {
	return &RetriableError{message: message}
}
