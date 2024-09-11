package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flylan/apollo-config-lib/request"
	"github.com/flylan/apollo-config-lib/utils"
	"net/url"
)

const DEFAULT_NOTIFICATION_ID = -1

type Notification struct {
	NamespaceName  string `json:"namespaceName"`
	NotificationId int64  `json:"notificationId"`
}

type Notifications []Notification

type NotificationsParam struct {
	Client           *Client
	NotificationsMap map[string]int64
}

// 构建一个应用感知配置实例
func (c *Client) Notifications(a interface{}) *NotificationsParam {
	nm := map[string]int64{}
	switch v := a.(type) {
	case string:
		nm[v] = DEFAULT_NOTIFICATION_ID
	case []string:
		for _, namespaceName := range v {
			nm[namespaceName] = DEFAULT_NOTIFICATION_ID
		}
	case map[string]int64:
		nm = v
	}
	return &NotificationsParam{Client: c, NotificationsMap: nm}
}

// 应用感知配置更新
func (np *NotificationsParam) Get() (*Notifications, *request.Info, error) {
	//初始化
	info := &request.Info{}

	if np.NotificationsMap == nil || len(np.NotificationsMap) == 0 {
		return nil, info, errors.New("NotificationsMap is empty")
	}

	// 将map转换为JSON字符串
	notifications := make(Notifications, len(np.NotificationsMap))
	for namespaceName, notificationId := range np.NotificationsMap {
		notifications = append(notifications, Notification{NamespaceName: namespaceName, NotificationId: notificationId})
	}
	nj, err := json.Marshal(notifications)
	if err != nil {
		return nil, info, err
	}

	// 清空切片内容
	notifications = notifications[:0]

	//构建请求链接
	requestUrl := fmt.Sprintf(
		"%s/notifications/v2?appId=%s&cluster=%s&notifications=%s",
		np.Client.ConfigServerUrl,
		np.Client.AppID,
		np.Client.ClusterName,
		url.QueryEscape(string(nj)),
	)

	//发送get请求
	info, err = request.SendGetRequest(
		requestUrl,
		np.Client.AppID,
		np.Client.Secret,
		np.Client.RequestTimeout.GetNotifications,
		info,
	)
	if err != nil {
		return nil, info, err
	}

	//请求失败
	if !info.IsGetDataSuccess() {
		return nil, info, errors.New(string(info.ResponseBody))
	}

	//转换json字符串为结构体
	if !utils.IsByteSliceEmpty(info.ResponseBody) {
		err = json.Unmarshal(info.ResponseBody, &notifications)
		if err != nil {
			return nil, info, err
		}
	}

	return &notifications, info, nil
}
