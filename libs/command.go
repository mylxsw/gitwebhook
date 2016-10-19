package libs

import (
	"bufio"
	"io"
	"os/exec"
	"sync"
)

// OutputType 是输出类型
type OutputType string

const (
	Stdout OutputType = "STDOUT"
	Stderr OutputType = "STDERR"
)

type Output struct {
	Type    OutputType
	Content string
}

// 执行shell命令
func ExecShellCommand(command string, output chan Output) (err error) {
	cmd := exec.Command("/bin/sh", "-x", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		bindOutput(output, &stdout, Stdout)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bindOutput(output, &stderr, Stderr)
	}()

	if err = cmd.Start(); err != nil {
		return err
	}

	if err = cmd.Wait(); err != nil {
		return err
	}

	wg.Wait()

	return nil
}

func bindOutput(output chan Output, input *io.ReadCloser, outputType OutputType) error {
	reader := bufio.NewReader(*input)
	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			if err != io.EOF {
				return err
			}
			break
		}

		output <- Output{
			Type:    outputType,
			Content: line,
		}
	}

	return nil
}
