package utils

import (
	"os"
)

func GetFileName(filePath string) string {
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			return filePath[i+1:]
		}
	}
	return filePath
}

func WriteFile(filePath, fileContent string) error {
	// 获取文件的目录部分
	dir := filePath[:len(filePath)-len("/"+GetFileName(filePath))]
	// 使用 os.MkdirAll 创建目录，如果目录已经存在不会报错
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	// 写入文件内容
	err = os.WriteFile(filePath, []byte(fileContent), 0644)
	if err != nil {
		return err
	}
	return nil
}
