package models

//角色
type Role struct {
	RoleId					int64			`description:"角色id" json:"uId" xorm:"pk autoincr"`
	RoleName				string			`description:"角色名称" json:"uId"`
	RoleType				int64			`description:"角色类型" json:"uId"`
}

//用户角色关联
type UserRole struct {
	UId						int64			`description:"用户id" json:"uId"`
	RoleId					int64			`description:"角色id" json:"roleId"`
}