package file

import (
	"log"
	"os"
)

func SaveFile(content string, name string) {
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Printf("[saveFile]打开文件%s失败：%s", name, err.Error())
	}
	_, err = file.WriteString(content)
	if err != nil {
		log.Printf("写入%s文件内容失败：%s", name, err.Error())
	}
	defer file.Close()
}

func ReadFile(name string) string {
	bytes, err := os.ReadFile(name)
	if err != nil {
		log.Printf("[readFile]打开文件%s失败：%s", name, err.Error())
	}
	return string(bytes)
}

func ExistFile(name string) bool {
	_, err := os.Stat(name)
	if os.IsExist(err) {
		return true
	}
	return false
}
