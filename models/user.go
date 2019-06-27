package models

/*	1. 默认float32和float64映射到数据库中为float,real,double这几种类型，这几种数据库类型数据库的实现一般都是非精确的
		如果一定要作为查询条件，请将数据库中的类型定义为Numeric或者Decimal  `xorm:"Numeric"`
	2. 复合主键  Id(xorm.PK{1, 2})

*/

type User struct {
	UId       			int64			`description:"注册时间" json:"uId" xorm:"pk autoincr"`
	PhoneNumber			string			`description:"手机号" json:"phoneNumber"`
	NickName 			string			`description:"昵称" json:"nickName" xorm:"notnull "`		//string类型默认映射为varchar(255)
	Password 			string			`description:"密码" json:"password" xorm:"notnull"`
	Salt	 			string			`description:"密码" json:"salt" xorm:"notnull"`
	Gender        		int    			`description:"性别,1 男, 2 女" json:"gender" xorm:"notnull default 0"`
	Created           	int64  			`description:"注册时间" json:"created" xorm:"created"`
	Updated           	int64  			`description:"修改时间" json:"updated" xorm:"updated"`
	DeletedAt         	int64  			`description:"删除时间" json:"deleted" xorm:"deleted"`
}

