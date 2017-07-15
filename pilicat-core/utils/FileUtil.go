package utils

import (
	"io/ioutil"
	"os"
	"path"
)

func FileCreate(name string) (*os.File, error) {

	return os.Create(name)

}

func FileOpen(name string) (*os.File, error) {

	//文件读取可以使用 bufio、ioutil 库

	return os.Open(name)

}

func FileOpenMod(name string, flag int, perm os.FileMode) (*os.File, error) {
	return os.OpenFile(name, flag, perm)
}

func FileRead(name string) ([]byte, error) {
	return ioutil.ReadFile(name)
}

func MakeDir(path string, perm os.FileMode) error {
	//return os.Mkdir(path, perm)
	return os.MkdirAll(path, perm)
}

//  文件重命名
func FileReName(oldpath, newpath string) error {
	return os.Rename(oldpath, newpath)
}

//  删除文件
func FileDel(name string) error {
	return os.Remove(name)
}

//  删除整个目录
func DirDel(path string) error {
	return os.RemoveAll(path)
}

//  取文件的信息
func FileInfo(name string) (os.FileInfo, error) {
	//其他文件信息，可以通过FileOpen方法拿到os.File
	return os.Stat(name)

}

//  判断文件/或文件夹是否存在
func FileIsExist(name string) bool {
	var exist = true
	if _, err := os.Stat(name); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//  取文件所在的目录 例如：
//  services.FileDir("d://app//src//main.go") = d:/app/src
func FileDir(fullPath string) string {
	return path.Dir(fullPath)
}

//  取文件名，包含扩展名 例如：
//  services.FileFullName("d://main.go") = main.go
func FileFullName(fullPath string) string {
	return path.Base(fullPath)
}

//  取文件扩展名 例如：
//  services.FileExt("d://main.go") = .go
func FileExt(fullPath string) string {
	return path.Ext(fullPath)
}
