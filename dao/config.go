package dao

import (
	"database/sql"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strings"
)

const (
	DBUserName = "root"
	DBPassWord = "158931"
	DBName = "douyin"
	DBHost = "127.0.0.1"
	DBPort = "3306"
	RedisHost = DBHost
	RedisPort = "6379"
	RedisPassWord = ""
)

func InitDB() {
	LoadDB()
	LoadRedis()
}

var DB *sql.DB
func LoadDB() {
	//构建连接："用户名:密码@tcp(IP:端口)/数据库?charset=utf8"
	path := strings.Join([]string{DBUserName, ":", DBPassWord, "@tcp(", DBHost, ":", DBPort, ")/", DBName, "?charset=utf8"}, "")

	//打开数据库,前者是驱动名，所以要导入： _ "github.com/go-sql-driver/mysql"
	DB, _ = sql.Open("mysql", path)
	//设置数据库最大连接数
	DB.SetConnMaxLifetime(100)
	//设置上数据库最大闲置连接数
	DB.SetMaxIdleConns(10)
	//验证连接
	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}
}

func LoadRedis() {
	addr := strings.Join([]string{RedisHost,":",RedisPort},"")
	rdb = redis.NewClient(&redis.Options{
		Addr:addr,
		Password:RedisPassWord,
		DB: 0,	//0号数据库
	})
	_, err := rdb.Ping().Result()
	if err != nil {
		fmt.Println("ping err:", err)
		return
	}

}