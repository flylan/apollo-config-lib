package request

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"github.com/flylan/apollo-config-lib/utils"
	"time"
)

const (
	AUTHORIZATION_FORMAT      = "Apollo %s:%s"
	DELIMITER                 = "\n"
	HTTP_HEADER_AUTHORIZATION = "Authorization"
	HTTP_HEADER_TIMESTAMP     = "Timestamp"
)

func signature(timestamp, pathWithQuery, secret string) string {
	return hmacSha1Sign(timestamp+DELIMITER+pathWithQuery, secret)
}

func timestamp() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

func buildHttpHeaders(requestUrl, appID, secret, timestamp string) (map[string]string, error) {
	pathWithQuery, err := url2PathWithQuery(requestUrl)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		HTTP_HEADER_AUTHORIZATION: fmt.Sprintf(
			AUTHORIZATION_FORMAT,
			appID,
			signature(timestamp, pathWithQuery, secret),
		),
		HTTP_HEADER_TIMESTAMP: timestamp,
	}, nil
}

func url2PathWithQuery(urlString string) (string, error) {
	urlInfo, err := utils.ParseUrl(urlString)
	if err != nil {
		return "", err
	}
	return urlInfo.PathWithQuery, nil
}

// HmacSha1Sign 生成 HMAC-SHA1 签名并返回 Base64 编码的结果
func hmacSha1Sign(stringToSign, secret string) string {
	h := hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(stringToSign))
	signature := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(signature)
}
