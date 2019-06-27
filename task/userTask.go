package task

import (
	"strconv"
	"shici/util"
)

func TestTimedTask(){
	util.Logger.Info(strconv.FormatInt(util.UnixOfBeijingTime(), 10))
}