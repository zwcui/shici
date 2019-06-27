package models

//app弹出提示
type AlertMessage struct {
	AlertCode		string		`description:"提示信息码，forward开头表示跳转actionurl" json:"alertCode"`
	AlertMessage	string		`description:"提示信息" json:"alertMessage"`
}
