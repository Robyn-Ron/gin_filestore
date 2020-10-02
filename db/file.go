package db

import (
	"database/sql"
	"errors"
	"file_store_net_http/db/mysql"
	"fmt"
	"log"
	"runtime"
)

//函数是否需要返回error? 一般都需要, 一般返回error的func, 都不需要print error
//	一般不需要返回值的, 都可以只返回一个error, 而需要返回值的, 则返回需要的返回值+error对象;
func OnFileUploadFinished(filehash, filename string, filesize int64, fileaddr string) (bool,error) {
	s,err := mysql.DBConn().Prepare(
		"insert into tbl_file(`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) values(?,?,?,?,1)")
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return false, err
	}
	defer s.Close()

	res, err := s.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return false, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return false, err
	}
	if rows <= 0 {
		_, fn, line, _ := runtime.Caller(0)
		err = errors.New(fmt.Sprintln("The db inserted zero line in error"))
		log.Println("filename:", fn, " line:", line, " error:", err)
		return false, err
	}
	return true, nil
}

type TableFileMeta struct{
	FileHash sql.NullString
	FileName sql.NullString//可以为nil
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func GetFileMeta(filehash string) (*TableFileMeta, error){
	stat,err := mysql.DBConn().Prepare("SELECT file_sha1, file_name, file_size, file_addr from tbl_file where file_sha1 = ? and status = 1 limit 1")
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return nil, err
	}
	defer stat.Close()

	row := stat.QueryRow(filehash)
	res := TableFileMeta{}
	err = row.Scan(&res.FileHash, &res.FileName, &res.FileSize, &res.FileAddr)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return nil, err
	}
	return &res, nil
}
