package libs

import (
	"bytes"
	"os/exec"
	"log"
)

// 执行shell命令
func ExecShellCommand(command string) (result string, err error) {
	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Printf("执行命令 %s 失败: %s", command, panicErr)
		}
	}()

	var stdOut, stdErr bytes.Buffer

	cmd := exec.Command("/bin/sh", "-c", command)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err = cmd.Run(); err != nil {
		panic(err.Error())
	}

	if stdErr.Len() != 0 {
		panic(stdErr.String())
	}

	result = stdOut.String()
	return
}
