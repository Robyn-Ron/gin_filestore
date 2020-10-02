package meta

import (
	"file_store_net_http/db"
	"log"
	"runtime"
)

//类似于文件描述符的结构体定义: 1.index; 2.fdNumber; 3.name; 分别提供给OS, User;
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string //上传时间戳
}

//存储所有文件信息(缺点: 没有保存在磁盘中, 可能会丢失数据)
var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

func GetFileMeta(filesha1 string) FileMeta {
	return fileMetas[filesha1]
}

func RemoveFileMeta(fileSha1 string){
	delete(fileMetas, fileSha1)
}

func UpdateFileMetaDB(fmeta FileMeta) bool {
	flag, err := db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
	}
	return flag
}

//查询DB中对象的属性值时, 主要的查询逻辑:
//	1. 构建查询table得到的对象, 这个对象的属性用sql module中的数据类型来定义;
//	2. 凡是对象, 都使用*指针作为返回值;
func GetFileMetaDB(filehash string) (*FileMeta, error) {
	tableFileMeta, err := db.GetFileMeta(filehash)
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		log.Println("filename:", fn, " line:", line, " error:", err)
		return nil, err
	}
	res := FileMeta{
		FileSha1: tableFileMeta.FileHash.String,
		FileName: tableFileMeta.FileName.String,
		FileSize: tableFileMeta.FileSize.Int64,
		Location: tableFileMeta.FileAddr.String,
	}
	return &res, nil
	
}
