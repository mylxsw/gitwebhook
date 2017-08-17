package libs

import (
	"fmt"
	"log"

	"github.com/mylxsw/git-web-hooks/gitlab"
)

// 部署任务参数
type TaskParam struct {
	Git          gitlab.GitLabObj
	Actions      []string
	TmplFilename string
	WebRoot      string
	Branch       string
	Servers      []string
}

// 执行部署命令
func ExecuteDeployTask(param TaskParam) error {
	outputs := make(chan Output, 10)
	defer close(outputs)

	// 输出channel，用于控制命令的输出
	go func() {
		for output := range outputs {
			log.Printf(
				"%s -> %s",
				output.Type,
				output.Content,
			)
		}
	}()

	// 解析模板，生成临时文件
	tempFilename := "/tmp/" + param.Git.After
	ParseTemplate(param.TmplFilename, tempFilename, param)

	servers := ""
	for _, server := range param.Servers {
		servers += " -o root@" + server
	}

	command := fmt.Sprintf("orgalorg --lock-file %s %s -i %s -C bash", tempFilename, servers, tempFilename)
	if err := ExecShellCommand(command, outputs); err != nil {
		return err
	}

	log.Printf("操作完成.")

	return nil
}
