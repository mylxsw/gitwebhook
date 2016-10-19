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
	"syscall"

	. "aicode.cc/GitWebHook/gitlab"
	"aicode.cc/GitWebHook/libs"
)

var host = flag.String("host", "0.0.0.0", "服务监听地址")
var port = flag.Int("port", 61001, "服务监听端口")

func webHookHandler(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responseError(resp, err, http.StatusInternalServerError)
		return
	}

	log.Printf("Request: %s", body)

	var params GitLabObj
	err = json.Unmarshal(body, &params)
	if err != nil {
		responseError(resp, err, http.StatusBadRequest)
		return
	}

	if params.ObjectKind == ObjectKindPush {
		log.Printf("检测到push请求")
		command := "ls"
		result, err := libs.ExecShellCommand(command)
		if err != nil {
			responseError(resp, err, http.StatusInternalServerError)
			return
		}

		log.Printf("Result: %s", result)
		log.Printf("Hello, world")
	}

	responseSuccess(resp)
}

func responseSuccess(resp http.ResponseWriter) {
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(`{"status":"ok"}`))
}

func responseError(resp http.ResponseWriter, err error, code int) {
	log.Printf("Error: %s", err)
	http.Error(resp, err.Error(), http.StatusBadRequest)
}

func main() {
	flag.Parse()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGUSR2, syscall.SIGQUIT)
		for {
			s := <-c
			fmt.Println("get signal:", s)
		}

	}()

	fmt.Printf("Listening to %s:%d...\n", *host, *port)

	http.HandleFunc("/", webHookHandler)
	err := http.ListenAndServe(*host+":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatalf("创建Webhook服务失败: %s", err)
	}
}
