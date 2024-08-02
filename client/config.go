package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flylan/apollo-config-lib/request"
	"github.com/flylan/apollo-config-lib/utils"
	"net"
	"net/url"
	"sync"
)

var ipCache sync.Map

type ConfigsParam struct {
	Client                                     *Client
	Ip                                         net.IP
	UseNoCacheApi                              bool
	NamespaceName, ReleaseKey, Messages, Label string
}

type Configurations map[string]string

type Configs struct {
	AppID          string         `json:"appId"`
	Cluster        string         `json:"cluster"`
	NamespaceName  string         `json:"namespaceName"`
	Configurations Configurations `json:"configurations"`
	ReleaseKey     string         `json:"releaseKey"`
}

// 构建一个获取配置实例
func (c *Client) Configs(namespaceName string) *ConfigsParam {
	//应用部署的机器ip
	var ip net.IP
	if cache, ok := ipCache.Load(c.address); ok {
		ip = cache.(net.IP)
	} else {
		ip, _ = utils.GetOutboundIP(c.address)
		ipCache.Store(c.address, ip)
	}
	return &ConfigsParam{
		Client:        c,
		Ip:            ip,
		UseNoCacheApi: true,
		NamespaceName: namespaceName,
	}
}

// 从Apollo读取配置
func (cp *ConfigsParam) Get() (*Configs, *request.Info, error) {
	//初始化info
	info := &request.Info{}
	//必须传入NamespaceName
	if cp.NamespaceName == "" {
		return nil, info, errors.New("NamespaceName is empty")
	}
	if cp.UseNoCacheApi {
		return cp.noCacheGet(info)
	}
	return cp.get(info)
}

// 构造基础请求链接
func (cp *ConfigsParam) buildBaseURL(format string) string {
	return fmt.Sprintf(
		format,
		cp.Client.ConfigServerUrl,
		cp.Client.AppID,
		cp.Client.ClusterName,
		cp.NamespaceName,
	)
}

// 发起获取配置请求
func (cp *ConfigsParam) sendGetRequest(requestUrl string) (*request.Info, error) {
	return request.SendGetRequest(
		requestUrl,
		cp.Client.AppID,
		cp.Client.Secret,
		cp.Client.RequestTimeout.GetConfigs,
	)
}

// 通过带缓存的Http接口从Apollo读取配置
func (cp *ConfigsParam) get(info *request.Info) (*Configs, *request.Info, error) {
	//构建请求链接
	requestUrl := cp.buildBaseURL("%s/configfiles/json/%s/%s/%s")
	if !utils.IsByteSliceEmpty(cp.Ip) {
		requestUrl = fmt.Sprintf("%s?ip=%s", requestUrl, cp.Ip)
	}

	//发送get请求
	info, err := cp.sendGetRequest(requestUrl)
	if err != nil {
		return nil, info, err
	}

	//带缓存接口只需要判断200状态码
	if !info.IsGetDataSuccess() {
		return nil, info, errors.New(fmt.Sprintf("%s returns HTTP status code: %d", requestUrl, info.StatusCode))
	}

	//转换json字符串为结构体
	var configurations = new(Configurations)
	if !utils.IsByteSliceEmpty(info.ResponseBody) {
		err = json.Unmarshal(info.ResponseBody, configurations)
		if err != nil {
			return nil, info, err
		}
	}

	return &Configs{
		AppID:          cp.Client.AppID,
		Cluster:        cp.Client.ClusterName,
		NamespaceName:  cp.NamespaceName,
		Configurations: *configurations,
	}, info, err
}

// 通过不带缓存的Http接口从Apollo读取配置
func (cp *ConfigsParam) noCacheGet(info *request.Info) (*Configs, *request.Info, error) {
	requestUrl := cp.buildBaseURL("%s/configs/%s/%s/%s")
	params := url.Values{}

	//上一次的releaseKey
	if cp.ReleaseKey != "" {
		params.Add("releaseKey", cp.ReleaseKey)
	}

	//最新的 notificationId
	if cp.Messages != "" {
		params.Add("messages", cp.Messages)
	}

	//灰度配置的标签
	if cp.Label != "" {
		params.Add("label", cp.Label)
	}

	//应用部署的机器ip
	if !utils.IsByteSliceEmpty(cp.Ip) {
		params.Add("ip", cp.Ip.String())
	}

	//构建最终请求链接
	queryStr := params.Encode()
	if queryStr != "" {
		requestUrl = fmt.Sprintf("%s?%s", requestUrl, queryStr)
	}

	//发送get请求
	info, err := cp.sendGetRequest(requestUrl)
	if err != nil {
		return nil, info, err
	}

	//不带缓存接口，可能返回200或者304状态码
	if !info.IsGetDataSuccess() && !info.IsDataNotModified() {
		return nil, info, errors.New(fmt.Sprintf("%s returns HTTP status code: %d", requestUrl, info.StatusCode))
	}

	//转换json字符串为结构体
	var configs = new(Configs)
	if !utils.IsByteSliceEmpty(info.ResponseBody) {
		err = json.Unmarshal(info.ResponseBody, configs)
		if err != nil {
			return nil, info, err
		}
	}

	return configs, info, nil
}
