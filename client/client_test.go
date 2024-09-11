package client

import (
	"sync"
	"testing"
)

var (
	client   *Client
	syncOnce sync.Once
)

func TestNewClient(t *testing.T) {
	configServerUrl := "http:///81.68.181.139:8080"
	appId := "apollo-client-test"
	var err error
	_, err = NewClient("", appId)
	if err == nil {
		t.Fatal("NewClient should return error when configServerUrl is empty")
	}
	_, err = NewClient(configServerUrl, "")
	if err == nil {
		t.Fatal("NewClient should return error when appId is empty")
	}
	_, err = NewClient("tcp://81.68.181.139:8080", "apollo-client-test")
	if err == nil {
		t.Fatal("NewClient should return error when url is invalid")
	}
	_, err = NewClient("://81.68.181.139:8080", "apollo-client-test")
	if err == nil {
		t.Fatal("NewClient should return error when url is invalid")
	}
	_, err = NewClient("http://81.68.181.139:50080", "apollo-client-test")
	if err == nil {
		t.Fatal("NewClient should return error when port is not reachable")
	}
}

func testNewClient(t *testing.T) *Client {
	client, err := NewClient("http://81.68.181.139:8080", "apollo-client-test")
	if err != nil {
		t.Fatal(err)
	}
	client.Secret = "4081edabfe4e4ba097cc16defc526c2f"
	return client
}

func testGetClient(t *testing.T) *Client {
	syncOnce.Do(func() {
		client = testNewClient(t)
	})
	return client
}
