package request

import (
	"crypto/tls"
	"errors"
	"github.com/flylan/apollo-config-lib/utils"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	METHOD_GET = "GET"

	DEFAULT_DIAL_KEEP_ALIVE         = 60 * time.Second
	DEFAULT_DIAL_TIMEOUT            = 1 * time.Second
	DEFAULT_MAX_IDLE_CONNS          = 512
	DEFAULT_MAX_IDLE_CONNS_PER_HOST = 512
)

var (
	transportInsecure     *http.Transport
	transportSecure       *http.Transport
	transportInsecureOnce sync.Once
	transportSecureOnce   sync.Once
)

func newHttpTransport() *http.Transport {
	return &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConns:        DEFAULT_MAX_IDLE_CONNS,
		MaxIdleConnsPerHost: DEFAULT_MAX_IDLE_CONNS_PER_HOST,
		DialContext: (&net.Dialer{
			KeepAlive: DEFAULT_DIAL_KEEP_ALIVE,
			Timeout:   DEFAULT_DIAL_TIMEOUT,
		}).DialContext,
	}
}

// 获取httpTransport
func getTransport(insecureSkipVerify bool) *http.Transport {
	if insecureSkipVerify {
		transportInsecureOnce.Do(
			func() {
				transportInsecure = newHttpTransport()
				transportInsecure.TLSClientConfig = &tls.Config{
					InsecureSkipVerify: insecureSkipVerify,
				}
			},
		)
		return transportInsecure
	}

	transportSecureOnce.Do(
		func() {
			transportSecure = newHttpTransport()
		},
	)
	return transportSecure
}

// 发送http GET请求
func SendGetRequest(requestUrl, appID, secret string, timeout time.Duration) (*Info, error) {
	info := &Info{RequestUrl: requestUrl}

	if requestUrl == "" {
		return info, errors.New("RequestUrl is empty")
	}

	req, err := http.NewRequest(METHOD_GET, requestUrl, nil)
	if err != nil {
		return info, err
	}
	info.RequestHeaders = req.Header

	//配置了秘钥就要生成相应的request headers
	if secret != "" {
		headers, err := buildHttpHeaders(requestUrl, appID, secret, timestamp())
		if err != nil {
			return info, err
		}
		for key, value := range headers {
			if key != "" {
				req.Header.Set(key, value)
			}
		}
	}

	resp, err := (&http.Client{Timeout: timeout, Transport: getTransport(req.URL.Scheme == utils.SCHEME_HTTPS)}).Do(req)
	if err != nil {
		return info, err
	}
	info.ResponseHeaders = resp.Header
	info.StatusCode = resp.StatusCode

	//读取Respone body
	body, err := ReadBody(resp)
	if err != nil {
		return info, err
	}
	info.ResponseBody = body

	return info, nil
}

// 读取正文内容
func ReadBody(resp *http.Response) (body []byte, err error) {
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	return body, err
}
