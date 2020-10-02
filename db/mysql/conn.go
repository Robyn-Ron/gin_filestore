package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
	"runtime"
)

var (
	db *sql.DB
)

func init() {
	db,_ = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	db.SetConnMaxLifetime(1000)
	err := db.Ping()
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Fatalln("filename:", fn, " line:", line, " error:", err)
		os.Exit(1)
	}
}

func DBConn() *sql.DB{
	return db
}
