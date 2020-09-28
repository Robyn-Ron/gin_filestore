package utils

import (
	"path/filepath"
	"runtime"
)

var (
	projectAbPath string
	projectName string
)

func init() {
	initProjectPath()
}

func initProjectPath() {
	_, filename, _, _ := runtime.Caller(0)
	utilsPath := filepath.Dir(filename)
	projectAbPath = filepath.Dir(utilsPath)
	projectName = filepath.Base(projectAbPath)
}

func GetCurAbPathDir() string {
	_, file, _, _ := runtime.Caller(1)
	abDir := filepath.Dir(file)
	return abDir
}

func GetFileAbPath(filepaths ... string) string {
	var res string = projectAbPath
	for _, filename := range filepaths {
		res = filepath.Join(res, filename)
	}
	return res
}
