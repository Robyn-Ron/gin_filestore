package utils

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"hash"
	"io"
	"os"
	"path/filepath"
)

//ctrypo加密:
//	方式1: 直接使用: sha1.Sum(data []byte), 返回加密后的数据;
//	方式2: 创建sha1.New(), 使用带有buffer缓存的对象缓存data数据, 再使用sha1.Sum(nil)加密数据;

type Sha1Stream struct{
	_sha1 hash.Hash
}

//使用sha1加密数据, 同时使用16进制编码密文, 减少发送数据大小;
//使用内建struct, 可以类似于实现了strBuilder, 可以不断在_sha1的buffer中添加数据;
func Sha1(data []byte)string {
	//创建对象
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func MD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum(nil))
}

func FileMD5(file *os.File) string {
	_md5 := md5.New()
	io.Copy(_md5, file)
	return hex.EncodeToString(_md5.Sum(nil))
}

func (this *Sha1Stream)Update(data []byte) {
	if this._sha1 == nil {
		this._sha1 = sha1.New()
	}
	this._sha1.Write(data)
}

func (this *Sha1Stream)Sum() string {
	return hex.EncodeToString(this._sha1.Sum(nil))
}

func PathExists(path string) (bool,error) {
	//获取Stat()有2中方式: 1. File对象的Stat()方法; 2.os.Stat(filepath string)方法;
	_, err := os.Stat(path)
	if err == nil {
		return true,nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, nil
}

//先判断文件是否存在, 再得到文件大小;
func GetFileSize(path string)int64{
	var res int64
	filepath.Walk(path, func(root string, f os.FileInfo, err error)error{
		res = f.Size()
		return nil
	})
	return res
}
