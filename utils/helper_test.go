package utils

import (
	"fmt"
	"net"
	"testing"
)

func TestParseUrl(t *testing.T) {
	checkErrorUrl(t, "ht://81.68.181.139/a/b/c?d=123&e=456")
	checkErrorUrl(t, "tcp://81.68.181.139")
	checkErrorUrl(t, "udp://81.68.181.139")
	checkErrorUrl(t, "http:///81.68.181.139")
	checkErrorUrl(t, "http:/81.68.181.139")
	checkErrorUrl(t, "http:81.68.181.139")
	checkErrorUrl(t, "http//81.68.181.139")
	checkErrorUrl(t, "//81.68.181.139")
	checkErrorUrl(t, "http:://81.68.181.139")
	checkErrorUrl(t, "http://81.68.181.139:65536")
	checkErrorUrl(t, "https://81.68.181.139:-1")

	checkUrlInfo(t, "http", "www.baidu.com", "80", "/a/b/c?d=123&e=456")
	checkUrlInfo(t, "https", "www.baidu.com", "443", "/h/i?j=789")
	checkUrlInfo(t, "http", "81.68.181.139", "50080", "/haha")
	checkUrlInfo(t, "https", "81.68.181.139", "50443", "/hello/world?d=123&e=456")
}

func TestGetOutboundIP(t *testing.T) {
	testGetOutboundIP(t, "www.baidu.com:80")
	testGetOutboundIP(t, "www.google.com:443")
	testGetOutboundIP(t, "8.8.8.8:1234")
	testGetOutboundIP(t, "8.8.4.4:12345")
	testGetOutboundIP(t, "8.8.8.8:443")
	testGetOutboundIP(t, "8.8.4.4:443")
}

func TestIsByteSliceEmpty(t *testing.T) {
	if !IsByteSliceEmpty([]byte{}) {
		t.Fatal("IsByteSliceEmpty failed, test case 1")
	}
	if !IsByteSliceEmpty(nil) {
		t.Fatal("IsByteSliceEmpty failed, test case 2")
	}
	if IsByteSliceEmpty([]byte{'1'}) {
		t.Fatal("IsByteSliceEmpty failed, test case 3")
	}
}

func TestCompareInt64(t *testing.T) {
	if CompareInt64(10, 1000) != -1 {
		t.Fatal("CompareInt64 failed, test case 1")
	}
	if CompareInt64(1000, 1000) != 0 {
		t.Fatal("CompareInt64 failed, test case 2")
	}
	if CompareInt64(1000, 10) != 1 {
		t.Fatal("CompareInt64 failed, test case 3")
	}
}

func TestStrToInt64(t *testing.T) {
	testStrToInt64(t, "abc", false)
	testStrToInt64(t, "123123a", false)
	testStrToInt64(t, "54321", true)
	testStrToInt64(t, "123123123123123123", true)
	testStrToInt64(t, "哈哈", false)
}

func TestIsValidJSON(t *testing.T) {
	testIsValidJSON(t, []byte(`{"a":{"b": {"c": ["d"]}, "e": 123}}`), true)
	testIsValidJSON(t, []byte(`[]`), true)
	testIsValidJSON(t, []byte(`{"a":{"b": {"c": ["d"]}, "e": 123`), false)
	testIsValidJSON(t, []byte(`{"a":{"b": {"c": ["d"]}, "e": 123,}`), false)
	testIsValidJSON(t, []byte(`{"a":{"b": {"c": ["d",]}, "e": 123,}`), false)
	testIsValidJSON(t, []byte(`{"a":{"b": {"c": ["d",]}, "e": 123}`), false)
}

func testIsValidJSON(t *testing.T, s []byte, expectedRes bool) {
	res := IsValidJSON(s)
	if (expectedRes && !res) || (!expectedRes && res) {
		t.Fatal(fmt.Sprintf("StrToInt64 failed, s: %s, res: %v", string(s), res))
	}
}

func testStrToInt64(t *testing.T, str string, expectedErrIsNil bool) {
	res, err := StrToInt64(str)
	if (expectedErrIsNil && err != nil) || (!expectedErrIsNil && err == nil) {
		t.Fatal(fmt.Sprintf("StrToInt64 failed, str: %s, res: %d, err: %v", str, res, err))
	}
}

func testGetOutboundIP(t *testing.T, address string) {
	ip, err := GetOutboundIP(address)
	if err != nil {
		t.Fatal(err)
	}
	if !isPrivateIPv4(ip) {
		t.Fatal(fmt.Sprintf("%s is not private", ip))
	}
}

func isPrivateIPv4(ip net.IP) bool {
	if ip == nil {
		return false
	}
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range privateRanges {
		_, ipNet, _ := net.ParseCIDR(cidr)
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

func checkErrorUrl(t *testing.T, errURL string) {
	res, err := ParseUrl(errURL)
	if err == nil {
		t.Fatal(fmt.Sprintf("When the URL is invalid, an error message should be returned, url: %s, res: %v", errURL, res))
	}
}

func checkUrlInfo(t *testing.T, scheme, host, port, pathWithQuery string) {
	address := host + ":" + port
	h := address
	if port == "80" || port == "443" {
		h = host
	}
	var testURL1 = scheme + "://" + h + pathWithQuery
	res, err := ParseUrl(testURL1)
	if err != nil {
		t.Fatal(err)
	}
	if res.Host != host {
		t.Fatal(fmt.Sprintf("Host does not match, expected %s, got %s", host, res.Host))
	}
	if res.Port != port {
		t.Fatal(fmt.Sprintf("Port does not match, expected %s, got %s", port, res.Port))
	}
	if res.Address != address {
		t.Fatal(fmt.Sprintf("Address does not match, expected %s, got %s", address, res.Address))
	}
	if res.Scheme != scheme {
		t.Fatal(fmt.Sprintf("Scheme does not match, expected %s, got %s", scheme, res.Scheme))
	}
	if res.PathWithQuery != pathWithQuery {
		t.Fatal(fmt.Sprintf("PathWithQuery does not match, expected %s, got %s", pathWithQuery, res.PathWithQuery))
	}
}
