package meta

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

