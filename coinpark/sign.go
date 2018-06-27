package coinpark

import "github.com/leek-box/sheep/util"

func CreateSign(secret, cmds string) string {
	return util.ComputeHmacMd5(cmds, secret)
}
