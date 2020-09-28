package main

import (
	"CloudWebOfGin/store/ceph"
	"fmt"
	"gopkg.in/amz.v1/s3"
)

func main()  {
	bucket := ceph.GetCephBucket("testbucket")

	//创建一个新的bucket
	err := bucket.PutBucket(s3.PublicRead)
	fmt.Printf("create bucket err: %v", err)
	//查询这个bucket下面指定条件的object keys
	res, err := bucket.List("", "", "", 100)
	fmt.Printf("%+v", res)

	//新上传一个对象
	err = bucket.Put("/testupload/a.txt", []byte("just for test"), "octet-stream", s3.PublicRead)
	fmt.Printf("upload err:%+v\n", err)

	//查询这个bucket下面指定条件的 object keys
	res, err = bucket.List("", "", "", 100)
	fmt.Printf("%+v", res)
}