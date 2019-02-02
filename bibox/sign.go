package bibox

import (
	"github.com/gpmn/sheep/util"
)

func CreateSign(secret string, cmds string) string {
	return util.ComputeHmacMd5(string(cmds), secret)
}
