package util

//隐藏昵称后几位，用*替代
func FormatNickname(nickname string, length int) (formatName string) {
	rs := []rune(nickname)
	rl := len(rs)
	end := length

	if end < 0 {
		end = 0
	}
	if end > rl {
		end = rl
	}
	formatName = string(rs[0:end]) + "*"
	return formatName
}

//手机号中间打*号
func FormatPhoneNo(phoneNo string) (formatString string) {
	if len(phoneNo) < 11{
		return phoneNo
	}
	rs := []rune(phoneNo)

	formatString = string(rs[0:3]) + "****" +  string(rs[7:11])
	return formatString
}


