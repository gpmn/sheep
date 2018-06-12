package fcoin

import (
	"encoding/base64"
	"github.com/leek-box/sheep/util"
	"strconv"
)

// 构造签名
// mapParams: 送进来参与签名的参数, Map类型
// strMethod: 请求的方法 GET, POST......
// strHostUrl: 请求的主机
// strRequestPath: 请求的路由路径
// strSecretKey: 进行签名的密钥
func CreateSign(strMethod, strRequestPath, secretKey string, urlParams, postParams map[string]string, unix int64) string {
	//1.HTTP_METHOD + HTTP_REQUEST_URI + TIMESTAMP + POST_BODY
	sortedUrlParams := util.MapSortByKey(urlParams)
	sortedPostParams := util.MapSortByKey(postParams)

	if len(urlParams) != 0 {
		strRequestPath = strRequestPath + "?" + util.Map2UrlQuery(sortedUrlParams)
	}
	// strconv.FormatInt(time.Now().Unix(), 10)

	strPayload := strMethod + FCoinHost + strRequestPath + strconv.FormatInt(unix, 10)

	if len(sortedPostParams) != 0 {
		strPayload = strPayload + util.Map2UrlQuery(sortedPostParams)
	}

	strPayloadBased := base64.StdEncoding.EncodeToString([]byte(strPayload))

	return util.ComputeHmacSha1(strPayloadBased, secretKey)

}
