package db

import (
	"file_store_net_http/db/mysql"
	"log"
	"runtime"
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

func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) error {
	stmt, err := mysql.DBConn().Prepare(
		"INSERT INTO tbl_user_file(`user_name`, `file_sha1`, `file_name`, `file_size`, `upload_at`) VALUES(?,?,?,?,?)")
	if err != nil{
		_, fn, line, _ := runtime.Caller(0)
		log.Println("file: ", fn, "; line: ", line, "; error: ", err)
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil{
		_, fn, line, _ := runtime.Caller(0)
		log.Println("file: ", fn, "; line: ", line, "; error: ", err)
		return err
	}
	row, err := res.RowsAffected()
	if err != nil{
		_, fn, line, _ := runtime.Caller(0)
		log.Println("file: ", fn, "; line: ", line, "; error: ", err)
		return err
	}
	if row <= 0 {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("file: ", fn, "; line: ", line, "; error: ", err)
		return err
	}
	return nil
}

func QueryUserFileMetas(username string, limit int)([]UserFile, error) {
	stmt, err := mysql.DBConn().Prepare(
		"SELECT file_sha1, file_name, file_size, upload_at, last_update FROM tbl_user_file WHERE user_name=? limit ?")
	if err != nil{
		_, fn, line, _ := runtime.Caller(0)
		log.Println("file: ", fn, "; line: ", line, "; error: ", err)
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(username, limit)
	if err != nil{
		_, fn, line, _ := runtime.Caller(0)
		log.Println("file: ", fn, "; line: ", line, "; error: ", err)
		return nil, err
	}

	var userFiles []UserFile

	for rows.Next(){
		tmpufile := UserFile{}
		err = rows.Scan(&tmpufile.FileHash, &tmpufile.FileName, &tmpufile.FileSize, &tmpufile.UploadAt, &tmpufile.LastUpdated)
		if err != nil{
			_, fn, line, _ := runtime.Caller(0)
			log.Println("file: ", fn, "; line: ", line, "; error: ", err)
			return nil, err
		}

		userFiles = append(userFiles, tmpufile)
	}

	return userFiles, nil
}