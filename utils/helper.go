package utils

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"strconv"
)

const (
	QUESTION_MARK = "?"
	SCHEME_HTTP   = "http"
	SCHEME_HTTPS  = "https"

	SCHEME_HTTP_PORT  = "80"
	SCHEME_HTTPS_PORT = "443"

	NETWORK_TCP = "tcp"
	NETWORK_UDP = "udp"
)

type UrlInfo struct {
	Host          string
	Port          string
	Address       string
	Scheme        string
	PathWithQuery string
}

// 解析url（只支持解析http和https协议的）
func ParseUrl(u string) (*UrlInfo, error) {
	urlInfo, err := url.Parse(u)
	if err != nil {
		return nil, err
	}

	if urlInfo.Scheme != SCHEME_HTTP && urlInfo.Scheme != SCHEME_HTTPS {
		return nil, fmt.Errorf("Only supports HTTP or HTTPS protocols, url: %s", u)
	}

	// 检查hostname
	hostName := urlInfo.Hostname()
	if hostName == "" {
		return nil, fmt.Errorf("Unable to resolve hostname from url: %s", u)
	}

	// path和query部分
	pathWithQuery := urlInfo.Path
	if urlInfo.RawQuery != "" {
		pathWithQuery += QUESTION_MARK + urlInfo.RawQuery
	}

	// 提取端口号
	port := urlInfo.Port()
	if port == "" {
		switch urlInfo.Scheme {
		case SCHEME_HTTP:
			port = SCHEME_HTTP_PORT
		case SCHEME_HTTPS:
			port = SCHEME_HTTPS_PORT
		}
	}

	//检查端口合法性
	portInt, err := StrToInt64(port)
	if err != nil {
		return nil, fmt.Errorf("The host port is not an integer, current: %s", port)
	}
	if CompareInt64(portInt, 0) == -1 || CompareInt64(portInt, 65535) == 1 {
		return nil, fmt.Errorf("The host port must be between 0 and 65535, current: %s", port)
	}

	return &UrlInfo{
		Host:          hostName,
		Port:          port,
		Address:       fmt.Sprintf("%s:%s", hostName, port),
		Scheme:        urlInfo.Scheme,
		PathWithQuery: pathWithQuery,
	}, nil
}

// 获取出网（本地）ip
func GetOutboundIP(address string) (net.IP, error) {
	conn, err := net.Dial(NETWORK_UDP, address)
	if err != nil {
		return nil, err
	}
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	defer func() {
		_ = conn.Close()
	}()
	return localAddr.IP, nil
}

// 判断一个字节切片是否为空
func IsByteSliceEmpty(s []byte) bool {
	return s == nil || len(s) == 0
}

// CompareInt64 比较两个int64大小
// num1 < num2返回-1
// num1 > num2返回1
// num1 == num2返回0
func CompareInt64(num1, num2 int64) int {
	if num1 < num2 {
		return -1
	}
	if num1 > num2 {
		return 1
	}
	return 0
}

// 字符串转为int64类型
func StrToInt64(str string) (int64, error) {
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}

// 判断字节切片是否为合法的json字符串
func IsValidJSON(s []byte) bool {
	if IsByteSliceEmpty(s) {
		return false
	}
	var js json.RawMessage
	return json.Unmarshal(s, &js) == nil
}
