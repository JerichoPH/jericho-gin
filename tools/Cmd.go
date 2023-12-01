package tools

import (
	"io"
	"io/ioutil"
	"log"
	"os/exec"
)

type Cmd struct{}

// Process 执行终端命令
func (receiver Cmd) Process(name string, params ...string) (bool, string) {
	cmd := exec.Command(name, params...)

	if stdout, err := cmd.StdoutPipe(); err != nil { // 获取输出对象，可以从该对象中读取输出结果
		log.Fatal(err)
		return false, err.Error()
	} else {
		if err := cmd.Start(); err != nil { // 运行命令
			log.Fatal(err)
			return false, err.Error()
		}
		defer func(stdout io.ReadCloser) (bool, string) {
			err := stdout.Close()
			if err != nil {
				log.Fatal(err)
				return false, err.Error()
			}
			return true, ""
		}(stdout)

		if opBytes, err := ioutil.ReadAll(stdout); err != nil { // 读取输出结果
			log.Fatal(err)
			return false, err.Error()
		} else {
			return true, string(opBytes)
		}
	}
}
