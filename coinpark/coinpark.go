package coinpark

import (
	"errors"

	"log"

	"github.com/leek-box/sheep/util"
)

const CoinParkHost = "https://api.coinpark.cc"

type CoinPark struct {
	accessKey string
	secretKey string
}

func NewCoinPark(accessKey, secretKey string) (*CoinPark, error) {
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("access key or secret key error")
	}
	f := &CoinPark{
		accessKey: accessKey,
		secretKey: secretKey,
	}

	return f, nil
}

func Ping() bool {
	path := "/v1/public"
	var req = map[string]string{
		"cmd": "ping",
	}

	ret := util.HttpGetRequest(CoinParkHost+path+"?"+util.Map2UrlQuery(req), req)

	log.Println(ret)

	return true

}
