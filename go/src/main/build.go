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
		fmt.Println("ostype:", r.Form.Get("ostype"))
		fmt.Println("disksize:", r.Form.Get("disksize"))
		fmt.Println("software:", r.Form["software"])
		fmt.Println("service:", r.Form["service"])
		for _, v := range r.Form["software"] {
			fmt.Println("v:=", v)
		}
		for k, v := range r.Form {
			fmt.Println("key:", k)
			fmt.Println("val:", strings.Join(v, " "))
		}
		fmt.Println("cloudstackip:", r.Form.Get("cloudstackip"))
		json := buildjson(r)
		callpacker(json)
	}
}

func buildjson(r *http.Request) (json string) {
	jsondir := "result/test3/"
	os.MkdirAll(jsondir, 0777)
	json = "test1.json"
	json = jsondir + json
	f, _ := os.Create(json)
	f.WriteString("just a test")
	templatedir := "template/json/"
	template := jsonmap[r.Form.Get("ostype")]
	template = templatedir + template
	cfgdir := "template/cfg/"
	cfg := cfgmap[r.Form.Get("ostype")]
	cfg = cfgdir + cfg
	isodir := "template/iso/"
	iso := isomap[r.Form.Get("ostype")]
	iso = isodir + iso
	scriptdir := "template/script/"
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
	ft, _ := os.Open(template)
	buf := bufio.NewReader(ft)
	for {
		line, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		f.WriteString(line)
	}

	fmt.Println(json)
	fmt.Println(cfg)
	fmt.Println(script)
	fmt.Println(iso)
	defer ft.Close()
	defer f.Close()
	return json
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
