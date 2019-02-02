package coinpark

import "github.com/gpmn/sheep/util"

func CreateSign(secret, cmds string) string {
	return util.ComputeHmacMd5(cmds, secret)
}
