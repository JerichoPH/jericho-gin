package tools

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"jericho-gin/wrongs"
	"os"
)

type FileSystem struct {
	path string
}

// NewFileSystem 构造函数
func NewFileSystem(path string) FileSystem {
	return FileSystem{path: path}
}

// IsExist 判断路径是否存在
func (receiver FileSystem) IsExist() bool {
	_, err := os.Stat(receiver.path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	return true
}

// GetPath 获取当前路径
func (receiver FileSystem) GetPath() string {
	return receiver.path
}

// SetPath 设置路径
func (receiver FileSystem) SetPath(path string) FileSystem {
	receiver.path = path
	return receiver
}

// ReadString 以文字读取文件内容
func (receiver FileSystem) ReadString() string {
	f, err := ioutil.ReadFile(receiver.path)
	if err != nil {
		wrongs.ThrowForbidden("文件读取失败:", err.Error())
	}
	return string(f)
}

// ReadJson 读取json文件
func (receiver FileSystem) ReadJson() interface{} {
	f, err := ioutil.ReadFile(receiver.path)
	if err != nil {
		wrongs.ThrowForbidden("文件读取失败:", err.Error())
	}

	var t interface{}
	err = json.Unmarshal(f, &t)
	if err != nil {
		wrongs.ThrowForbidden("JSON解析失败: %s", err.Error())
	}

	return t
}

// WriteJson 写入json文件
func (receiver FileSystem) WriteJson(value any) FileSystem {
	if bytes, err := json.Marshal(value); err != nil {
		wrongs.ThrowForbidden("JSON解析失败: %s", err.Error())
	} else {
		receiver.WriteString(string(bytes))
	}
	return receiver
}

// WriteString 写入string文件
func (receiver FileSystem) WriteString(value string) FileSystem {
	file, err := os.OpenFile(receiver.path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		wrongs.ThrowForbidden("文件打开失败: %s", err.Error())
	}
	// 及时关闭file句柄
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			wrongs.ThrowForbidden("文件关闭失败: %s", err.Error())
		}
	}(file)

	// 写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	_, err = write.WriteString(value)
	if err != nil {
		wrongs.ThrowForbidden("文件写入失败: %s", err.Error())
	}

	// Flush将缓存的文件真正写入到文件中
	err = write.Flush()
	if err != nil {
		wrongs.ThrowForbidden("文件写入失败: %s", err.Error())
	}

	return receiver
}
