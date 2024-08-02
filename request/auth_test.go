package request

import (
	"fmt"
	"github.com/flylan/apollo-config-lib/utils"
	"testing"
	"time"
)

func TestSignature(t *testing.T) {
	if signature(
		"1722245200513",
		"/configfiles/json/apollo-client-test/default/application?ip=10.4.123.251",
		"4081edabfe4e4ba097cc16defc526c2f",
	) != "fdgkCEFZMa1quB9XAuJq2Cu5PAY=" {
		t.Fatal("test case 1 signature error")
	}

	if signature(
		"1722245333026",
		"/notifications/v2?appId=apollo-client-test&cluster=default&notifications=%5B%7B%22namespaceName%22%3A%22%22%2C%22notificationId%22%3A0%7D%2C%7B%22namespaceName%22%3A%22application%22%2C%22notificationId%22%3A-1%7D%5D",
		"2024edabfe4e4ba097cc16defc526c2f",
	) != "7q+rTjdNr9r83EmhNhWXmBh69KA=" {
		t.Fatal("test case 2 signature error")
	}
}

func TestUrl2PathWithQuery(t *testing.T) {
	var res string
	var err error
	pathWithQuery := "/a/b/c?d=123&e=456"

	//检测错误判断
	errURL1 := "htps:/www.baidu.com" + pathWithQuery
	_, err = url2PathWithQuery(errURL1)
	if err == nil {
		t.Fatal(fmt.Sprintf("When the URL is invalid, an error message should be returned, url: %s", errURL1))
	}

	res, err = url2PathWithQuery("https://www.google.com" + pathWithQuery)
	if err != nil {
		t.Fatal(err)
	}
	if res != pathWithQuery {
		t.Fatal(fmt.Sprintf("url2PathWithQuery error, expect: %s, but: %s", pathWithQuery, res))
	}
}

func TestHmacSha1Sign(t *testing.T) {
	stringToSign := "1722246984182\n/notifications/v2?appId=apollo-client-test&cluster=default&notifications=%5B%7B%22namespaceName%22%3A%22%22%2C%22notificationId%22%3A0%7D%2C%7B%22namespaceName%22%3A%22application%22%2C%22notificationId%22%3A-1%7D%5D"
	secret := "4081edabfe4e4ba097cc16defc526c2f"
	expect := "pGULDgQ31tCkZU+GQANYKtvhp9E="
	but := hmacSha1Sign(stringToSign, secret)
	if but != expect {
		t.Fatal(fmt.Sprintf("hmacSha1Sign error, expect: %s, but: %s", expect, but))
	}
}

func TestTimestamp(t *testing.T) {
	ts1 := time.Now().UnixNano() / int64(time.Millisecond)
	ts2, err := utils.StrToInt64(timestamp())
	if err != nil {
		t.Fatal(err)
	}
	if utils.CompareInt64(ts1-1000, ts2) == 1 || utils.CompareInt64(ts1+1000, ts2) == -1 {
		t.Fatal(fmt.Sprintf("timestamp error, ts1: %d, ts2: %d", ts1, ts2))
	}
}

func TestBuildHttpHeaders(t *testing.T) {
	requestUrl := "http://81.68.181.139:8080/notifications/v2?appId=apollo-client-test&cluster=default&notifications=%5B%7B%22namespaceName%22%3A%22%22%2C%22notificationId%22%3A0%7D%2C%7B%22namespaceName%22%3A%22application%22%2C%22notificationId%22%3A-1%7D%5D"
	appID := "apollo-client-test"
	secret := "4081edabfe4e4ba097cc16defc526c2f"
	ts1 := "1722309152123"
	headers, err := buildHttpHeaders(requestUrl, appID, secret, ts1)
	if err != nil {
		t.Fatal(err)
	}
	authorization, exists := headers["Authorization"]
	if !exists {
		t.Fatal(fmt.Sprintf("authorization error, authorization header not found"))
	}
	if authorization != "Apollo "+appID+":zPGup4xc51rjWAfcxXfAFn4PCE4=" {
		t.Fatal(fmt.Sprintf("authorization error, authorization header is wrong, authorization header: %s", authorization))
	}

	ts2, exists := headers["Timestamp"]
	if !exists {
		t.Fatal(fmt.Sprintf("timestamp error, timestamp header not found"))
	}
	if ts2 != ts1 {
		t.Fatal(fmt.Sprintf("timestamp error, ts1: %s, ts2: %s", ts1, ts2))
	}
}
