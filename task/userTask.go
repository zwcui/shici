package task

import (
	"strconv"
	"baseApi/util"
)

func TestTimedTask(){
	util.Logger.Info(strconv.FormatInt(util.UnixOfBeijingTime(), 10))
}