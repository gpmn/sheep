package okex

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"log"

	"github.com/gpmn/sheep/util"
)

const (
	apiURL     = "https://www.okex.com/api/"
	apiVersion = "v1/"
)

// 将map格式的请求参数转换为字符串格式的
// mapParams: map格式的参数键值对
// return: 查询字符串
func map2UrlQuery(mapParams map[string]string) string {
	var keys []string
	for key := range mapParams {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	var strParams string
	for _, key := range keys {
		strParams += (key + "=" + mapParams[key] + "&")
	}

	if 0 < len(strParams) {
		strParams = string([]rune(strParams)[:len(strParams)-1])
	}

	return strParams
}

func httpGetRequest(strUrl string, mapParams map[string]string) string {
	httpClient := &http.Client{}

	var strRequestUrl string
	if nil == mapParams {
		strRequestUrl = strUrl
	} else {
		strParams := map2UrlQuery(mapParams)
		strRequestUrl = strUrl + "?" + strParams
	}

	// 构建Request, 并且按官方要求添加Http Header
	request, err := http.NewRequest("GET", strRequestUrl, nil)
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")

	// 发出请求
	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	// 解析响应内容
	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}

func httpPostRequest(strUrl string, values url.Values, accessKey, secretKey string) string {
	httpClient := &http.Client{}

	values.Set("api_key", accessKey)

	hasher := util.MD5([]byte(values.Encode() + "&secret_key=" + secretKey))
	values.Set("sign", strings.ToUpper(util.HexEncodeToString(hasher)))

	request, err := http.NewRequest("POST", strUrl, strings.NewReader(values.Encode()))
	if nil != err {
		return err.Error()
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Accept-Language", "zh-cn")

	response, err := httpClient.Do(request)
	if nil != err {
		return err.Error()
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if nil != err {
		return err.Error()
	}

	return string(body)
}

func (o *OKEX) apiKeyPost(values url.Values, strRequestPath string, dst interface{}) error {
	strUrl := apiURL + apiVersion + strRequestPath
	resp := httpPostRequest(strUrl, values, o.accessKey, o.secretKey)
	log.Println(resp)
	return json.Unmarshal([]byte(resp), dst)
}
