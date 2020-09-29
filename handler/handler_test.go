package handler

import (
	"file_store_net_http/utils"
	"fmt"
	"io/ioutil"
	"runtime"
	"testing"
)

func TestReadFile(t *testing.T) {
	data, err := ioutil.ReadFile(utils.GetFileAbPath("handler", "handler.go"))
	if err != nil {
		_, fn, line, _ := runtime.Caller(0)
		fmt.Println("filename: ", fn, "_", line, "error: ", err)
		return
	}
	fmt.Println(string(data))
}
