package request

import (
	"fmt"
	"github.com/flylan/apollo-config-lib/utils"
	"testing"
	"time"
)

func TestHttp(t *testing.T) {
	requestUrl := "http://81.68.181.139:8080/notifications/v2?appId=apollo-client-test&cluster=default&notifications=%5B%7B%22namespaceName%22%3A%22%22%2C%22notificationId%22%3A0%7D%2C%7B%22namespaceName%22%3A%22application%22%2C%22notificationId%22%3A-1%7D%5D"
	appID := "apollo-client-test"
	secret := "4081edabfe4e4ba097cc16defc526c2f"
	timeout := 10 * time.Second
	info, err := SendGetRequest(requestUrl, appID, secret, timeout)
	if err != nil {
		t.Fatal(err)
	}
	if !utils.IsValidJSON(info.ResponseBody) {
		t.Fatal(fmt.Sprintf("invalid json: %s", string(info.ResponseBody)))
	}
	if !info.IsGetDataSuccess() {
		t.Fatal(fmt.Sprintf("error response: %v", info))
	}
	if info.IsDataNotModified() {
		t.Fatal(fmt.Sprintf("error response: %v", info))
	}
}
