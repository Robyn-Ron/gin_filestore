package db

import (
	"CloudWebOfGin/db/mysql"
	"fmt"
	"time"
)

type UserFile struct {
	UserName string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdated string
}

//上传用户文件接口
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt, err := mysql.DBConn().Prepare("INSERT INTO tbl_user_file(`user_name`, `file_sha1`, `file_name`, `file_size`, `upload_at`) VALUES(?,?,?,?,?)")
	if err != nil{
		fmt.Println(err.Error())
		return false
	}

	defer stmt.Close()

	result, err := stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil{
		fmt.Println(err.Error())
		return false
	}

	if count, err := result.RowsAffected(); err != nil{
		fmt.Println(err.Error())
		return false
	}else if count <= 0{
		return false
	}

	return true
}

//批量获取用户文件信息
func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mysql.DBConn().Prepare("SELECT file_sha1, file_name, file_size, upload_at, last_update FROM tbl_user_file WHERE user_name=? limit ?")
	if err != nil{
		fmt.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil{
		fmt.Println(err.Error())
		return nil, err
	}

	var userFiles []UserFile

	for rows.Next(){
		tmpufile := UserFile{}
		err = rows.Scan(&tmpufile.FileHash, &tmpufile.FileName, &tmpufile.FileSize, &tmpufile.UploadAt, &tmpufile.LastUpdated)
		if err != nil{
			fmt.Println(err.Error())
			return nil, err
		}

		userFiles = append(userFiles, tmpufile)
	}

	return userFiles, nil
}