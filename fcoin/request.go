package fcoin

import (
	"strconv"
	"time"

	"github.com/gpmn/sheep/util"
)

// 进行签名后的HTTP GET请求, 参考官方Python Demo写的
// mapParams: map类型的请求参数, key:value
// strRequest: API路由路径
// return: 请求结果
func apiKeyGet(mapParams map[string]string, strRequestPath string, accessKey, secretKey string) string {
	strMethod := "GET"
	now := time.Now()
	timestamp := now.UnixNano() / 1000 / 1000

	var resParams = map[string]string{
		"FC-ACCESS-KEY":       accessKey,
		"FC-ACCESS-SIGNATURE": CreateSign(strMethod, strRequestPath, secretKey, mapParams, nil, timestamp),
		"FC-ACCESS-TIMESTAMP": strconv.FormatInt(timestamp, 10),
	}
	strUrl := FCoinHost + strRequestPath
	if len(mapParams) > 0 {
		strUrl = strUrl + "?" + util.Map2UrlQuery(mapParams)
	}
	return util.HttpGetRequest(strUrl, resParams)
}

// 进行签名后的HTTP POST请求, 参考官方Python Demo写的
// mapParams: map类型的请求参数, key:value
// strRequest: API路由路径
// return: 请求结果
func apiKeyPost(mapParams map[string]string, strRequestPath string, accessKey, secretKey string) string {
	strMethod := "POST"
	now := time.Now()
	timestamp := now.UnixNano() / 1000 / 1000

	var resParams = map[string]string{
		"FC-ACCESS-KEY":       accessKey,
		"FC-ACCESS-SIGNATURE": CreateSign(strMethod, strRequestPath, secretKey, nil, mapParams, timestamp),
		"FC-ACCESS-TIMESTAMP": strconv.FormatInt(timestamp, 10),
	}

	return util.HttpPostRequest(FCoinHost+strRequestPath, mapParams, resParams)
}
