package db

import (
	"database/sql"
	"errors"
	"file_store_net_http/db/mysql"
	"fmt"
	"log"
	"runtime"
)

func UserSignup(username, encrypt_pasword, phone string) error{
	s, err :=mysql.DBConn().Prepare(
		"insert into tbl_user(`user_name`,`user_pwd`,`phone`) values(?,?,?)")
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return err
	}
	s.Close()

	res, err := s.Exec(username, encrypt_pasword, phone)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return err
	}
	cnt, err := res.RowsAffected()
	if cnt <= 0 || err != nil {
		_, fn, line, _ := runtime.Caller(0)
		errors.New("插入用户信息失败")
		log.Println("filename:", fn, " line:", line, " error:", err)
		return err
	}
	return nil
}

func UserSignin(username, encrypt_password string) error {
	s, err := mysql.DBConn().Prepare(
		"SELECT * FROM tbl_user WHERE user_name=? AND user_pwd=?")
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return err
	}
	s.Close()

	row,err := s.Query(username, encrypt_password)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return err
	}
	if row == nil {
		_, fn, line, _ := runtime.Caller(0)
		err =  errors.New("查询用户信息失败")
		log.Println("filename:", fn, " line:", line, " error:",err)
		return err
	}
	return nil
}

func UpdateToken(username, token string) error {
	s, err := mysql.DBConn().Prepare(
		"REPLACE INTO tbl_user_token(`user_name`, `user_token`) VALUES(?,?)")
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return err
	}
	s.Close()
	res, err := s.Exec(username, token)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return err
	}

	rows,err := res.RowsAffected()
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return err
	}
	if rows <= 0 {
		_, fn, line, _ := runtime.Caller(0)
		err = errors.New(fmt.Sprintln("The db inserted zero line in error"))
		log.Println("filename:", fn, " line:", line, " error:", err)
		return err
	}
	return nil
}

//这个struct应该是MySQL中table存储的表结构;
//之后还应该对这个数据对象处理一下, 提取必要的信息, 给前端返回一个处理过的对象信息;
type User struct{
	UserName sql.NullString
	Email sql.NullString
	Phone sql.NullString
	Status sql.NullString
}

func GetUserInfo(username string)(*User, error){
	stmt, err := mysql.DBConn().Prepare(
		"SELECT user_name, phone FROM tbl_user WHERE user_name=?;")
	if err != nil{
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return nil, err
	}
	defer stmt.Close()

	tmpUser := User{}
	row := stmt.QueryRow(username).Scan(&tmpUser.UserName, &tmpUser.Phone)
	if row != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", row.Error())
		return nil, errors.New(row.Error())
	}
	return &tmpUser, nil
}
