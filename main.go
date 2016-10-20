package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/mylxsw/git-web-hooks/gitlab"
	"github.com/mylxsw/git-web-hooks/libs"
	"github.com/mylxsw/git-web-hooks/pidfile"
)

var host = flag.String("host", "0.0.0.0", "服务监听地址")
var port = flag.Int("port", 61001, "服务监听端口")
var secretKey = flag.String("secret", "", "接口秘钥，用于验证来源是否合法")
var pidFile = flag.String("pidfile", "/tmp/gitwebhook.pid", "pid文件路径")
var concurrent = flag.Int("concurrent", 5, "并发执行线程数")

var stopRunning bool = false
var stopRunningChan chan struct{}
var taskParamChan chan libs.TaskParam = make(chan libs.TaskParam, 5)

// 请求参数：hosts, branch=master, webroot=/home/data, tmpl=./tmpl/demo.tmpl
func webHookHandler(resp http.ResponseWriter, req *http.Request) {
	// check the secret key
	if *secretKey != "" && req.FormValue("secret") != *secretKey {
		responseError(resp, fmt.Errorf("秘钥错误"), http.StatusUnauthorized)
		return
	}

	tmplFilename := req.FormValue("tmpl")
	if tmplFilename == "" {
		tmplFilename = "./tmpl/demo.tmpl"
	}
	webroot := req.FormValue("webroot")
	if webroot == "" {
		webroot = "/home/data"
	}
	branch := req.FormValue("branch")
	if branch == "" {
		branch = "master"
	}
	hosts := req.FormValue("hosts")
	if hosts == "" {
		responseError(resp, fmt.Errorf("缺少hosts"), http.StatusBadRequest)
		return
	}

	// 读取请求body，解析json数据为GitLabObj类型
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responseError(resp, err, http.StatusInternalServerError)
		return
	}
	log.Printf("Request: %s", body)

	var params gitlab.GitLabObj
	err = json.Unmarshal(body, &params)
	if err != nil {
		responseError(resp, err, http.StatusBadRequest)
		return
	}

	taskParamChan <- libs.TaskParam{
		Git:          params,
		TmplFilename: tmplFilename,
		WebRoot:      webroot,
		Branch:       branch,
		Servers:      strings.Split(hosts, ","),
	}

	responseSuccess(resp)
}

// success response
func responseSuccess(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(`{"status":200}`))
}

// error response
func responseError(resp http.ResponseWriter, err error, code int) {
	log.Printf("Error: %s", err)
	resp.Header().Set("Content-Type", "application/json")
	http.Error(resp, fmt.Sprintf(`{"status":500, "message": "%s"}`, err.Error()), http.StatusBadRequest)
}

func main() {
	flag.Parse()

	stopRunningChan = make(chan struct{}, *concurrent)

	// 创建进程pid文件
	pid, err := pidfile.New(*pidFile)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer pid.Remove()

	// 信号处理程序，接收退出信号，平滑退出进程
	initSignalReceiver()

	// 启动任务消费队列
	go func() {
		// 没有消息队列了就自动退出吧
		defer os.Exit(0)

		var wg sync.WaitGroup
		for i := 0; i < *concurrent; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()
				worker(i)
			}(i)
		}

		wg.Wait()
	}()
	log.Printf("Listening to %s:%d...\n", *host, *port)

	http.HandleFunc("/", webHookHandler)
	if err = http.ListenAndServe(*host+":"+strconv.Itoa(*port), nil); err != nil {
		log.Fatalf("创建Webhook服务失败: %s", err)
	}

}

func worker(i int) {
	defer func() {
		log.Printf("Task customer [%d] stopped.", i)
	}()

	log.Printf("Task customer [%d] started.", i)

	for {
		// worker exit
		if stopRunning {
			return
		}

		select {
		case taskParam := <-taskParamChan:
			if err := libs.ExecuteDeployTask(taskParam); err != nil {
				log.Printf("Error: %s", err)
			}
		case <-stopRunningChan:
			return
		}
	}
}

// 初始化信号接受处理程序
func initSignalReceiver() {
	signalChan := make(chan os.Signal)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGUSR2,
	)
	go func() {
		for {
			sig := <-signalChan
			switch sig {
			case syscall.SIGUSR2, syscall.SIGHUP:
				stopRunning = true
				//close(command)
				for i := 0; i < *concurrent; i++ {
					stopRunningChan <- struct{}{}
				}
				log.Print("Received exit signal.")
			}
		}
	}()

}
