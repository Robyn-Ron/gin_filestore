package db

import (
	"CloudWebOfGin/db/mysql"
	"database/sql"
	"fmt"
)

//注册用户
func UserSignup(username, password, phone string) bool {
	stmt, err := mysql.DBConn().Prepare("INSERT INTO tbl_user(`user_name`, `user_pwd`, `phone`) VALUES(?,?,?);")
	if err != nil{
		fmt.Println("1"+err.Error())
		return false
	}

	defer stmt.Close()


	result, err := stmt.Exec(username, password, phone)
	if err != nil{
		fmt.Println("2"+err.Error())
		return false
	}

	if count, err := result.RowsAffected(); err != nil{
		fmt.Println("3"+err.Error())
		return false
	}else {
		if count > 0{
			return true
		}else {
			return false
		}
	}
}

//登录
func UserSignin(username, password string) bool {
	stmt, err := mysql.DBConn().Prepare("SELECT * FROM tbl_user WHERE user_name=? AND user_pwd=? ;")
	if err != nil{
		fmt.Println(err.Error())
		return false
	}

	rows, err := stmt.Query(username, password)

	if err != nil{
		fmt.Println(err.Error())
		return false
	}else if rows == nil{
		//没有查找到对应数据
		fmt.Println("用户不存在或者密码错误")
		return false
	}

	return true
}

//保存token信息
func UpdateToken(username, token string) bool {
	stmt, err := mysql.DBConn().Prepare("REPLACE INTO tbl_user_token(`user_name`, `user_token`) VALUES(?,?);")
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, token)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}

	if count, err := result.RowsAffected(); err != nil{
		fmt.Println(err.Error())
		return false
	}else if count > 0{
		return true
	}else if count <= 0{
		return false
	}

	return true
}

type User struct {
	UserName sql.NullString
	Email sql.NullString
	Phone sql.NullString
	Status sql.NullInt32
}

//查询用户信息
func GetUserInfo(username string) (*User, error) {
	stmt, err := mysql.DBConn().Prepare("SELECT user_name, phone FROM tbl_user WHERE user_name=?;")
	if err != nil{
		fmt.Println(err.Error())
		return nil, err
	}
	defer stmt.Close()

	tmpUser := User{}

	row := stmt.QueryRow(username).Scan(&tmpUser.UserName, &tmpUser.Phone)
	if row != nil{
		fmt.Println(row.Error())
		return nil, row
	}

	return &tmpUser, nil
}