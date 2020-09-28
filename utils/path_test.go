package utils

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"
)

func TestProjectPath(t *testing.T) {
	dirPath := filepath.Dir(".")
	fmt.Println(dirPath)
	_,dirPath1,_,_ := runtime.Caller(0)
	_,dirPath2,_,_ := runtime.Caller(1)
	fmt.Println(dirPath1)
	fmt.Println(dirPath2)

	_, filename, _, _ := runtime.Caller(0)
	utilsPath := filepath.Dir(filename)
	projectPath := filepath.Dir(utilsPath)
	projectName := filepath.Base(projectPath)
	fmt.Println(utilsPath)
	fmt.Println(projectPath)
	fmt.Println(projectName)
}

func TestGetCurAbPathDir (t *testing.T) {
	fmt.Println(GetCurAbPathDir())
}

func TestGetFileAbPath(t *testing.T) {
	fmt.Println(GetFileAbPath("utils", "path.go"))
}
