package huobi

import (
	"time"

	"github.com/gpmn/sheep/util"
)

const host = "https://api.huobi.pro"

// 进行签名后的HTTP GET请求, 参考官方Python Demo写的
// mapParams: map类型的请求参数, key:value
// strRequest: API路由路径
// return: 请求结果
func apiKeyGet(mapParams map[string]string, strRequestPath string, accessKey, secretKey string) string {
	strMethod := "GET"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	mapParams["AccessKeyId"] = accessKey
	mapParams["SignatureMethod"] = "HmacSHA256"
	mapParams["SignatureVersion"] = "2"
	mapParams["Timestamp"] = timestamp

	hostName := "api.huobi.pro"
	mapParams["Signature"] = createSign(mapParams, strMethod, hostName, strRequestPath, secretKey)

	strUrl := host + strRequestPath
	return util.HttpGetRequest(strUrl, util.MapValueEncodeURI(mapParams))
}

// 进行签名后的HTTP POST请求, 参考官方Python Demo写的
// mapParams: map类型的请求参数, key:value
// strRequest: API路由路径
// return: 请求结果
func apiKeyPost(mapParams map[string]string, strRequestPath string, accessKey, secretKey string) string {
	strMethod := "POST"
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05")

	mapParams2Sign := make(map[string]string)
	mapParams2Sign["AccessKeyId"] = accessKey
	mapParams2Sign["SignatureMethod"] = "HmacSHA256"
	mapParams2Sign["SignatureVersion"] = "2"
	mapParams2Sign["Timestamp"] = timestamp

	hostName := "api.huobi.pro"

	mapParams2Sign["Signature"] = createSign(mapParams2Sign, strMethod, hostName, strRequestPath, secretKey)
	strUrl := host + strRequestPath + "?" + util.Map2UrlQuery(util.MapValueEncodeURI(mapParams2Sign))

	return util.HttpPostRequest(strUrl, mapParams, nil)
}
