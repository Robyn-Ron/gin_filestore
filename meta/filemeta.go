package meta

import (
	"CloudWebOfGin/db"
	"fmt"
)

// 文件元信息结构
type FileMeta struct {
	FileSha1 string //签名加密
	FileName string //文件名
	FileSize int64 //文件大小
	Location string //文件路径
	UploadAt string //上传时间戳
}

var fileMetas map[string]FileMeta

func init()  {
	fileMetas = make(map[string]FileMeta)
}

//上传/更新 文件元信息
func UpdateFileMeta(fmeta FileMeta)  {
	//旧版本：保存在内存的map对象中
	fileMetas[fmeta.FileSha1] = fmeta
}

//获取文件的元信息
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

//删除文件元信息
func RemoveFileMeta(fileSha1 string){
	delete(fileMetas, fileSha1)

}

//新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.UploadAt)
}

//从mysql获取文件元信息
func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	tfile, err := db.GetFileMeta(fileSha1)
	if err != nil{
		fmt.Println(err.Error())
		return nil, err
	}

	var fileMeta = &FileMeta{
		FileSha1: tfile.FileHash,
		FileSize: tfile.FileSize.Int64,
		FileName: tfile.FileName.String,
		Location: tfile.FileAddr.String,
	}

	return fileMeta, nil
}