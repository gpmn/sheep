package bibox

import (
	"github.com/leek-box/sheep/util"
)

func CreateSign(secret string, cmds string) string {
	return util.ComputeHmacMd5(string(cmds), secret)
}
