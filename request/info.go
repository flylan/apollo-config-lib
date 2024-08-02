package request

import "net/http"

type Info struct {
	RequestUrl      string
	RequestHeaders  http.Header
	StatusCode      int
	ResponseHeaders http.Header
	ResponseBody    []byte
}

// 判断是否获取数据成功
func (i *Info) IsGetDataSuccess() bool {
	return i.StatusCode == http.StatusOK
}

// 判断数据是否无变更
func (i *Info) IsDataNotModified() bool {
	return i.StatusCode == http.StatusNotModified
}
