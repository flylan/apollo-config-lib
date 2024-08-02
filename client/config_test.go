package client

import (
	"fmt"
	"github.com/flylan/apollo-config-lib/request"
	"net/http"
	"testing"
)

func TestPrintConfigs(t *testing.T) {
	client1, configs1, info1 := testConfigsGet(t, "TEAM.test_case_1", true)
	fmt.Println("client1", client1)
	fmt.Println("configs1", configs1)
	fmt.Println("info1", info1)
	fmt.Println()
	fmt.Println()

	client2, configs2, info2 := testConfigsGet(t, "TEAM.test_case_1", false)
	fmt.Println("client2", client2)
	fmt.Println("configs2", configs2)
	fmt.Println("info2", info2)
}

func TestConfigsGet(t *testing.T) {
	checkConfigs(t, "application", false)
	checkConfigs(t, "application", true)
	checkConfigs(t, "haha", false)
	checkConfigs(t, "haha", true)
	checkConfigs(t, "TEAM.test_case_1", false)
	checkConfigs(t, "TEAM.test_case_1", true)
	checkConfigs(t, "test_case_2", false)
	checkConfigs(t, "test_case_2", true)
	checkConfigs(t, "test_case_3", false)
	checkConfigs(t, "test_case_3", true)
	checkConfigs(t, "TEAM.test_case_4", false)
	checkConfigs(t, "TEAM.test_case_4", true)
}

func checkConfigs(t *testing.T, namespaceName string, noCache bool) {
	client, configs, info := testConfigsGet(t, namespaceName, noCache)
	if info.StatusCode != http.StatusOK {
		t.Fatal(fmt.Sprintf("configs.Cluster: %d not equal to 200", info.StatusCode))
	}
	if noCache && configs.ReleaseKey == "" {
		t.Fatal("When using an interface without cache, configs.ReleaseKey should not be empty")
	}
	if configs.AppID != client.AppID {
		t.Fatal(fmt.Sprintf("configs.AppID: %s not equal to client.AppID: %s", configs.AppID, client.AppID))
	}
	if configs.Cluster != client.ClusterName {
		t.Fatal(fmt.Sprintf("configs.Cluster: %s not equal to client.ClusterName: %s", configs.Cluster, client.ClusterName))
	}
	if configs.NamespaceName != namespaceName {
		t.Fatal(fmt.Sprintf("configs.Cluster: %s not equal to namespaceName: %s", configs.Cluster, namespaceName))
	}
	if configs.Configurations == nil || len(configs.Configurations) == 0 {
		t.Fatal("configs.Configurations is empty")
	}
}

func testConfigsGet(t *testing.T, namespaceName string, noCache bool) (*Client, *Configs, *request.Info) {
	client := testGetClient(t)
	cp := client.Configs(namespaceName)
	if !noCache {
		cp.UseNoCacheApi = false
	}
	configs, info, err := cp.Get()
	if err != nil {
		t.Fatal(err)
	}
	if configs == nil {
		t.Fatal("configs empty")
	}
	return client, configs, info
}
