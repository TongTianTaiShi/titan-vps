package utils

import (
	"fmt"
	"github.com/opentracing/opentracing-go/log"
	"os"
	"strings"
)

func CreateFile(fileName string) (*os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		log.Error(err)
		return file, err
	}
	return file, nil

}
func ChmodFile(fileName string) {
	err := os.Chmod(fileName, 0600)
	if err != nil {
		fmt.Printf("os.Chmod err:%s \n", err.Error())
		log.Error(err)
		return
	}
	return

}

func CreateAndWriteFile(fileName, privateKey string) {
	// 读写|创建|追加的模式 模式打开文件
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Printf("os.OpenFile err:%s \n", err.Error())
		return
	}
	// write file
	write, err := file.Write([]byte(strings.Replace(privateKey, "\n", "", -1)))
	if err != nil {
		fmt.Printf("Write: %s \n", err.Error())
	}
	fmt.Printf("write bytes: %d \n", write)
	// 关闭文件
	_ = file.Close()
}

func RemoveFile(fileName string) {
	err := os.RemoveAll("fileName")
	if err != nil {
		fmt.Printf("os.OpenFile failed err:%s \n", err.Error())
		return
	}
}
