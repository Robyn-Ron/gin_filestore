package mysql

import (
	"database/sql"
	"fmt"
	"os"
)
import _ "github.com/go-sql-driver/mysql" //调用包的时候，就会初始化，并且将自己注册到database/sql里面去，然后再通过调用data/sql包里面提供的方法去操作mysql数据库

var db *sql.DB

//初始化数据库连接对象
func init()  {
	db, _ = sql.Open("mysql", "root:123456@tcp(192.168.243.181:3306)/fileserver?charset=utf8") //用户名：密码@传输方式（数据库IP：端口号）/数据库名称？字符集
	db.SetConnMaxLifetime(1000) //设置最大连接数
	err := db.Ping()
	if err != nil{
		fmt.Println("Failed to connect to mysql, err:" + err.Error())
		os.Exit(1)
	}
}

func DBConn() *sql.DB {
	return db
}