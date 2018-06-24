package bibox

import "errors"

const BiboxHost = "https://api.fcoin.com/v2/"

type Bibox struct {
	accessKey string
	secretKey string
}

func NewBibox(accessKey, secretKey string) (*Bibox, error) {
	if accessKey == "" || secretKey == "" {
		return nil, errors.New("access key or secret key error")
	}
	f := &Bibox{
		accessKey: accessKey,
		secretKey: secretKey,
	}

	return f, nil
}
