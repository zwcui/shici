package base

import (
	"github.com/astaxie/beego"
	"fmt"
	"github.com/go-xorm/xorm"
	"os"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/cache"
	_ "shici/models"
	"gopkg.in/mgo.v2"
	"github.com/garyburd/redigo/redis"
	"shici/util"
	"strconv"
)

//运行标识
const (
	RUN_MODE_DEV  = "dev"
	RUN_MODE_TEST = "test"
	RUN_MODE_PROD = "prod"
)

//数据库参数、redis参数、服务地址
var redisConn, serverURL string
var dbConfig databaseConfig
var rdConfig redisConfig
var mongoDBUrl, mongoDBName string

type databaseConfig struct {
	DbType   		string
	DbUser     		string
	DbPassword 		string
	DbName 			string
	DbCharset  		string
	DbHost     		string
	DbPort     		string
}

type redisConfig struct {
	RedisConn		string
	Auth			string
	Key				string
	DBNum			int
	//MaxIdle			int
	//MaxActive		int
	//IdleTimeout		time.Duration
}

//数据库引擎
var DBEngine *xorm.Engine

//Redis
var RedisCache cache.Cache
var redisPool *redis.Pool

//MongoDB
var MongoDBSession *mgo.Session

//系统初始化
func init(){

	if beego.BConfig.RunMode == RUN_MODE_DEV {
		serverURL = "http://106.14.202.179:8888"
		dbConfig.DbType = "mysql"
		dbConfig.DbHost = "106.14.202.179"
		dbConfig.DbPort = ":3306"
		dbConfig.DbUser = "shici"
		dbConfig.DbPassword = "shici"
		dbConfig.DbName = "shici"
		dbConfig.DbCharset = "utf8mb4"
		rdConfig.RedisConn = "106.14.202.179:6379"
		rdConfig.Auth = "baseapi"
		mongoDBUrl = "106.14.202.179:27017"
		mongoDBName = "shici"
	} else if beego.BConfig.RunMode == RUN_MODE_TEST {
		serverURL = "http://106.14.202.179:8888"
		dbConfig.DbType = "mysql"
		dbConfig.DbHost = "106.14.202.179"
		dbConfig.DbPort = ":3306"
		dbConfig.DbUser = "shici"
		dbConfig.DbPassword = "shici"
		dbConfig.DbName = "shici"
		dbConfig.DbCharset = "utf8mb4"
		rdConfig.RedisConn = "106.14.202.179:6379"
		rdConfig.Auth = "baseapi"
		mongoDBUrl = "106.14.202.179:27017"
		mongoDBName = "shici"
	} else if beego.BConfig.RunMode == RUN_MODE_PROD {
		serverURL = "http://106.14.202.179:8888"
		dbConfig.DbType = "mysql"
		dbConfig.DbHost = "106.14.202.179"
		dbConfig.DbPort = ":3306"
		dbConfig.DbUser = "shici"
		dbConfig.DbPassword = "shici"
		dbConfig.DbName = "shici"
		dbConfig.DbCharset = "utf8mb4"
		rdConfig.RedisConn = "106.14.202.179:6379"
		rdConfig.Auth = "baseapi"
		mongoDBUrl = "106.14.202.179:27017"
		mongoDBName = "shici"
	} else {
		panic("运行标识错误")
	}

	initDB(dbConfig)
	initRedis(rdConfig)
	initMongoDB()
}



//数据库初始化
func initDB(dbConfig databaseConfig){
	var err error
	//"root:123@/test?charset=utf8"
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s%s)/%s?charset=%s",
		dbConfig.DbUser, dbConfig.DbPassword, dbConfig.DbHost, dbConfig.DbPort, dbConfig.DbName, dbConfig.DbCharset)
	fmt.Println(dbUrl)
	DBEngine, err = xorm.NewEngine(dbConfig.DbType, dbUrl)
	if err != nil {
		panic("创建数据库连接Engine失败! err:"+err.Error())
	}
	DBEngine.ShowSQL(false)			//在控制台打印出生成的SQL
	DBEngine.SetMaxIdleConns(20)	//设置闲置的连接数
	DBEngine.SetMaxOpenConns(100)	//设置最大打开的连接数，默认值为0表示不限制
	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)	//启用一个全局的内存缓存，存放到内存中，缓存struct的记录数为1000条
	DBEngine.SetDefaultCacher(cacher)

	//SnakeMapper为默认值，结构体驼峰结构，表名转为下划线，可以不写。SameMapper为结构体与表名一致
	//表名前后缀 core.NewPrefixMapper(core.SnakeMapper{}, "prefix")  core.NewSufffixMapper(core.SnakeMapper{}, "suffix")
	//engine.SetMapper(core.SnakeMapper{})

	//engine.DBMetas()	//获取到数据库中所有的表，字段，索引的信息

	//engine.CreateTables()
	//engine.IsTableEmpty()
	//engine.IsTableExist()
	//engine.DropTables()
	//engine.CreateIndexes()
	//engine.CreateUniques()

	//engine.DumpAll()		//导出
	//engine.DumpAllToFile()
	//engine.Import()		//导入
	//engine.ImportFile()

	err = DBEngine.Ping()
	if err != nil {
		panic("数据库连接ping失败! err:"+err.Error())
	}

	//将sql写入到文件中
	f, err := os.Create("sql.log")
	if err != nil {
		panic("创建sql.log文件失败! err:"+err.Error())
	}
	 defer f.Close()
	DBEngine.SetLogger(xorm.NewSimpleLogger(f))

	//同步表结构
	//err = DBEngine.Sync2(new(models.User), new(models.Role), new(models.UserRole))
	//if err != nil {
	//	panic("同步表结构失败! err:"+err.Error())
	//}
}

//初始化redis
func initRedis(rdConfig redisConfig){
	var err error
	RedisCache, err = cache.NewCache("redis", `{"conn":"`+rdConfig.RedisConn+`", "key":"`+rdConfig.Key+`", "dbNum":"`+strconv.Itoa(rdConfig.DBNum)+`", "password":"`+rdConfig.Auth+`"}`)
	if err != nil {
		panic("redis初始化失败！err:"+err.Error())
	}
	RedisCache.Put("lastStartTime", strconv.FormatInt(util.UnixOfBeijingTime(), 10), 0)
}

//初始化MongoDB
func initMongoDB(){
	var err error
	session, err := mgo.Dial(mongoDBUrl)
	if err != nil {
		panic("MongoDB初始化失败！err:"+err.Error())
	}
	//session.SetMode(mgo.Monotonic, true)
	session.DB(mongoDBName)
}

//获取MondoDB的session
//调用完该方法需加 defer session.Close()
func MongoDB() (session *mgo.Session, database *mgo.Database){
	session, err := mgo.Dial(mongoDBUrl)
	if err != nil {
		panic("MongoDB session获取失败！err:"+err.Error())
	}
	session.SetMode(mgo.Monotonic, true)
	return session, session.DB(mongoDBName)
}