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
	"text/template"

	"strings"

	. "github.com/mylxsw/git-web-hooks/gitlab"
	"github.com/mylxsw/git-web-hooks/libs"
)

var host = flag.String("host", "0.0.0.0", "服务监听地址")
var port = flag.Int("port", 61001, "服务监听端口")

// 请求参数：hosts, branch=master, webroot=/home/data, tmpl=./tmpl/demo.tmpl
func webHookHandler(resp http.ResponseWriter, req *http.Request) {

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

	outputs := make(chan libs.Output, 10)
	defer close(outputs)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		responseError(resp, err, http.StatusInternalServerError)
		return
	}

	// log.Printf("Request: %s", body)

	var params GitLabObj
	err = json.Unmarshal(body, &params)
	if err != nil {
		responseError(resp, err, http.StatusBadRequest)
		return
	}

	if params.ObjectKind == ObjectKindPush {
		tmpl, err := template.ParseFiles(tmplFilename)
		if err != nil {
			responseError(resp, err, http.StatusInternalServerError)
			return
		}

		tempFile, err := os.Create("/tmp/" + params.After)
		if err != nil {
			responseError(resp, err, http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(tempFile, struct {
			Git     GitLabObj
			WebRoot string
			Branch  string
		}{
			Git:     params,
			WebRoot: webroot,
			Branch:  branch,
		})
		if err != nil {
			responseError(resp, err, http.StatusInternalServerError)
			return
		}

		tempFile.Close()

		servers := ""
		for _, ser := range strings.Split(hosts, ",") {
			servers += " -o root@" + ser
		}

		go func() {
			for output := range outputs {
				log.Printf(
					"%s -> %s",
					output.Type,
					output.Content,
				)
			}
		}()

		command := fmt.Sprintf("orgalorg %s -i %s -C bash", servers, "/tmp/"+params.After)
		if err = libs.ExecShellCommand(command, outputs); err != nil {
			responseError(resp, err, http.StatusInternalServerError)
			return
		}

		log.Printf("操作完成.")
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

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR2, syscall.SIGQUIT)
	go func() {
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
