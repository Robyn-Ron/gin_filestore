package db

import (
	"CloudWebOfGin/db/mysql"
	"database/sql"
	"fmt"
)

//上传文件
func OnFileUploadFinished(filehash, filename string, filesize int64, fileaddr string) bool {
	//获取到数据库连接
	stmt, err := mysql.DBConn().Prepare("INSERT INTO tbl_file(`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) VALUES (?,?,?,?,1);")
	defer stmt.Close()

	if err != nil{
		fmt.Println("Failed to prepare statement, err:" + err.Error())
		return false
	}

	result, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}

	if rf, err := result.RowsAffected(); err == nil {
		if rf <= 0{
			//没有新增数据进去
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}

		return true
	}

	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString //允许为null
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

//从mysql获取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mysql.DBConn().Prepare("SELECT file_sha1, file_addr, file_name, file_size FROM tbl_file WHERE file_sha1=? AND status=1 limit 1")
	if err != nil{
		fmt.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	tfile := TableFile{}

	row := stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)
	if row != nil{
		fmt.Println(row.Error())
		return nil, row
	}

	return &tfile, nil
}