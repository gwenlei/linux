package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var dirmap = map[string]string{
	"json":   "template/json/",
	"cfg":    "template/cfg/",
	"script": "template/script/",
	"iso":    "template/iso/",
	"result": "result/",
}
var jsonmap = map[string]string{
	"CentOS6.6":    "centos6-6.json",
	"CentOS6.7":    "centos6-7.json",
	"CentOS7.1":    "centos7-1.json",
	"Ubuntu12.04":  "ubuntu12-04.json",
	"Ubuntu14.04":  "ubuntu14-04.json",
	"OpenSuse13.2": "opensuse13-2.json",
	"Windows7":     "windows7.json",
	"Windows2012":  "windows2012.json",
}
var cfgmap = map[string]string{
	"CentOS6.6":    "centos6-6.cfg",
	"CentOS6.7":    "centos6-7.cfg",
	"CentOS7.1":    "centos7-1.cfg",
	"Ubuntu12.04":  "ubuntu12-04.cfg",
	"Ubuntu14.04":  "ubuntu14-04.cfg",
	"OpenSuse13.2": "opensuse13-2.cfg",
	"Windows7":     "windows7.cfg",
	"Windows2012":  "windows2012.cfg",
}
var isomap = map[string]string{
	"CentOS6.6":    "centos6-6.iso",
	"CentOS6.7":    "centos6-7.iso",
	"CentOS7.1":    "centos7-1.iso",
	"Ubuntu12.04":  "ubuntu12-04.iso",
	"Ubuntu14.04":  "ubuntu14-04.iso",
	"OpenSuse13.2": "opensuse13-2.iso",
	"Windows7":     "windows7.iso",
	"Windows2012":  "windows2012.iso",
}
var scriptmap = map[string]string{
	"mysql":     "mysql.sh",
	"wordpress": "wordpress.sh",
}

func main() {
	http.HandleFunc("/build", build)         //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func build(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("build.html")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		for k, v := range r.Form {
			fmt.Println(k,":", strings.Join(v, " "))
		}
		json := buildjson(r)
		fmt.Println(json)
                //callpacker(json)
	}
}

func buildjson(r *http.Request) (result string) {
	resultdir := dirmap["result"] + "test3/"
	os.MkdirAll(resultdir, 0777)
	result = "test1.json"
	result = resultdir + result
	f, _ := os.Create(result)
	f.WriteString("just a test")
	jsondir := dirmap["json"]
	json := jsonmap[r.Form.Get("ostype")]
	json = jsondir + json
	cfgdir := dirmap["cfg"]
	cfg := cfgmap[r.Form.Get("ostype")]
	cfg = cfgdir + cfg
	isodir := dirmap["iso"]
	iso := isomap[r.Form.Get("ostype")]
	iso = isodir + iso
	scriptdir := dirmap["script"]
	var script = make([]string, 10)
	n := copy(script, r.Form["software"])
	fmt.Println("n=", n)
	for k, v := range script {
		fmt.Println(k, v)
		script[k] = scriptdir + scriptmap[v]
		n = n - 1
		if n == 0 {
			break
		}
	}
	ft, _ := os.Open(json)
	buf := bufio.NewReader(ft)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		f.WriteString(line)
	}

	fmt.Println(result)
	fmt.Println(json)
	fmt.Println(cfg)
	fmt.Println(script)
	fmt.Println(iso)
	defer ft.Close()
	defer f.Close()
	return result
}

func callpacker(json string) {
	fmt.Println("callpacker", json)
	attr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	p, err := os.StartProcess("/home/packerdir/packer", []string{"/home/packerdir/packer", "build", "/home/jsondir/centos66.json"}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
}
