package client

import (
	"fmt"
	"github.com/flylan/apollo-config-lib/request"
	"testing"
)

func TestPrintNotifications(t *testing.T) {
	client1, notifications1, info1 := testNotificationsGet(t, []string{"application", "haha"})
	fmt.Println("client1", client1)
	fmt.Println("notifications1", notifications1)
	fmt.Println("info1", info1)
}

func TestNotificationsGet(t *testing.T) {
	checkNotifications(t, "application")
	checkNotifications(t, []string{"application", "haha"})
	checkNotifications(
		t,
		map[string]int64{
			"TEAM.test_case_1": DEFAULT_NOTIFICATION_ID,
			"test_case_2":      1,
			"test_case_3":      20,
			"TEAM.test_case_4": 300,
		},
	)
}

func checkNotifications(t *testing.T, a interface{}) {
	client, notifications, _ := testNotificationsGet(t, a)
	np := client.Notifications(a)
	if np == nil {
		t.Fatal("client.Notifications returned nil")
	}
	if np.NotificationsMap == nil || len(np.NotificationsMap) == 0 {
		t.Fatal("np.NotificationsMap is empty")
	}
	if len(np.NotificationsMap) != len(*notifications) {
		t.Fatal(fmt.Sprintf("client.Notifications returned wrong number of notifications, np.NotificationsMap: %v, *notifications: %v", np.NotificationsMap, *notifications))
	}
	for _, notification := range *notifications {
		notificationId, exists := np.NotificationsMap[notification.NamespaceName]
		if !exists {
			t.Fatal("client.Notifications did not contain notification")
		}
		if notificationId > notification.NotificationId {
			t.Fatal("client.Notifications notification id should be greater than notification id")
		}
	}
}

func testNotificationsGet(t *testing.T, a interface{}) (*Client, *Notifications, *request.Info) {
	client := testGetClient(t)
	notifications, info, err := client.Notifications(a).Get()
	if err != nil {
		t.Fatal(err)
	}
	if notifications == nil || len(*notifications) == 0 {
		t.Fatal("notifications empty")
	}
	return client, notifications, info
}
