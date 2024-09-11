package client

import (
	"errors"
	"fmt"
	"github.com/flylan/apollo-config-lib/utils"
	"net"
	"time"
)

const DEFAULT_CLUSTER_NAME = "default"

type RequestTimeout struct {
	GetConfigs       time.Duration
	GetNotifications time.Duration
}

type Client struct {
	ConfigServerUrl string
	AppId           string
	ClusterName     string
	Secret          string
	RequestTimeout  RequestTimeout
	address         string
}

func NewClient(configServerUrl, appId string) (*Client, error) {
	if configServerUrl == "" {
		return nil, errors.New("ConfigServerUrl is empty")
	}
	if appId == "" {
		return nil, errors.New("AppId is empty")
	}

	//解析url
	urlInfo, err := utils.ParseUrl(configServerUrl)
	if err != nil {
		return nil, err
	}

	//检测主机端口是否可以访问
	conn, err := net.DialTimeout(utils.NETWORK_TCP, urlInfo.Address, 3*time.Second)
	if err != nil {
		return nil, fmt.Errorf("Port %s on %s is closed or not reachable", urlInfo.Port, urlInfo.Host)
	}
	defer func() { _ = conn.Close() }()

	return &Client{
		ConfigServerUrl: configServerUrl,
		AppId:           appId,
		ClusterName:     DEFAULT_CLUSTER_NAME,
		address:         urlInfo.Address,
		RequestTimeout: RequestTimeout{
			GetConfigs:       10 * time.Second,
			GetNotifications: 60 * time.Second,
		},
	}, nil
}
