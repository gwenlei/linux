package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
        "io/ioutil"
)

var dirmap = map[string]string{
	"json":   "template/json/",
	"cfg":    "template/cfg/",
	"script": "template/script/",
	"iso":    "/home/html/iso/",
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
	"CentOS6.6":    "CentOS-6.6-x86_64-bin-DVD1.iso",
	"CentOS6.7":    "centos6-7.iso",
	"CentOS7.1":    "centos7-1.iso",
	"Ubuntu12.04":  "ubuntu12-04.iso",
	"Ubuntu14.04":  "ubuntu14-04.iso",
	"OpenSuse13.2": "opensuse13-2.iso",
	"Windows7":     "windows7.iso",
	"Windows2012":  "windows2012.iso",
}
var md5map = map[string]string{
	"CentOS-6.6-x86_64-bin-DVD1.iso": "7b1fb1a11499b31271ded79da6af8584",
	"centos6-7.iso":                  "12345",
	"centos7-1.iso":                  "12345",
	"ubuntu12-04.iso":                "12345",
	"ubuntu14-04.iso":                "12345",
	"opensuse13-2.iso":               "12345",
	"windows7.iso":                   "12345",
	"windows2012.iso":                "12345",
}
var scriptmap = map[string]string{
	"mysql":     "mysql.sh",
	"wordpress": "wordpress.sh",
}
var httpmap = map[string]string{
	"http://192.168.0.82/cfg/": "/home/html/cfg/",
        "http://192.168.0.82:9090/static/cfg/": "static/cfg/",
}
var resultmap = map[string]string{
	"http":      "http://192.168.0.82:9090/static/cfg/",
	"outputdir":    "static/result/test3/output/",
        "outputimage":    "centos66.qcow2",
	"jsondir":   "static/result/test3/",
	"scriptdir": "static/result/test3/script/",
}
var logmap = map[string]string{
"1":"7574",
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
        http.HandleFunc("/", sayHelloName)
	http.HandleFunc("/build", build)         //设置访问的路由
        http.HandleFunc("/setting", setting)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func sayHelloName(w http.ResponseWriter, r *http.Request) {

        fmt.Fprintf(w, "hello, world")
}

func setting(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("setting.html")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		for k, v := range r.Form {
			fmt.Println(k, ":", strings.Join(v, " "))
		}
	}
}

func build(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("build2.html")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		for k, v := range r.Form {
			fmt.Println(k, ":", strings.Join(v, " "))
		}
		json := buildjson(r)
		fmt.Println("buildjson end", json)
                p:=callpacker(json)
		fmt.Fprintf(w, strconv.Itoa(p.Pid)+" packer running")
                if r.Form.Get("compat")=="0.1"{
                  output:=resultmap["outputdir"]+resultmap["outputimage"]
                  newoutput:=resultmap["outputdir"]+"tr"+resultmap["outputimage"]
                  go calltransform(p,output,newoutput)
                }
         }
}

func buildjson(r *http.Request) (result string) {
	os.MkdirAll(resultmap["jsondir"], 0777)
	os.MkdirAll(resultmap["scriptdir"], 0777)
        os.MkdirAll(httpmap[resultmap["http"]], 0777)
	jsondir := dirmap["json"]
	json := jsonmap[r.Form.Get("ostype")]
	newjson := resultmap["jsondir"] + json
	json = jsondir + json
	cfgdir := dirmap["cfg"]
	cfg := cfgmap[r.Form.Get("ostype")]
	newcfg := httpmap[resultmap["http"]] + cfg
	newcfgs := resultmap["http"] + cfg
	cfg = cfgdir + cfg
	isodir := dirmap["iso"]
	iso := isomap[r.Form.Get("ostype")]
	iso = isodir + iso
	disksizen, _ := strconv.Atoi(r.Form.Get("disksize"))
	disksizens := strconv.Itoa(disksizen * 1024)
	//new json file
        os.Create(newjson)
	buf,_:= ioutil.ReadFile(json)
          line:= string(buf)
		line = strings.Replace(line, "DISK_SIZE", disksizens, -1)
		line = strings.Replace(line, "SSH_USERNAME", r.Form.Get("user"), -1)
		line = strings.Replace(line, "SSH_PASSWORD", r.Form.Get("password"), -1)
		line = strings.Replace(line, "VM_NAME", resultmap["outputimage"], -1)
		line = strings.Replace(line, "OUTPUT_DIRECTORY", resultmap["outputdir"], -1)
		line = strings.Replace(line, "ISO_CHECKSUM", md5map[isomap[r.Form.Get("ostype")]], -1)
		line = strings.Replace(line, "ISO_URL", iso, -1)
		line = strings.Replace(line, "KS_CFG", newcfgs, -1)
		line = strings.Replace(line, "HEADLESS", r.Form.Get("headless"), -1)
	scriptdir := dirmap["script"]
	var script = make([]string, 10)
	var newscript = make([]string, 10)
	n := copy(script, r.Form["software"])
	copy(newscript, script)
	fmt.Println("n=", n)
	if n > 0 {
		var scriptfiles = "\"" + resultmap["scriptdir"] + "base.sh" + "\",\n"
		// copy script
		newbasescriptf, _ := os.Create(resultmap["scriptdir"] + "base.sh")
		basescriptf, _ := os.Open(scriptdir + "base.sh")
                io.Copy(newbasescriptf, basescriptf)
		defer basescriptf.Close()
		defer newbasescriptf.Close()
		for k, v := range script {
			fmt.Println(k, v)
			script[k] = scriptdir + scriptmap[v]
			newscript[k] = resultmap["scriptdir"] + scriptmap[v]
			scriptfiles = scriptfiles + "\"" + newscript[k] + "\""
			n = n - 1
			// copy script
			newscriptf, _ := os.Create(newscript[k])
			scriptf, _ := os.Open(script[k])
			io.Copy(newscriptf, scriptf)
			defer scriptf.Close()
			defer newscriptf.Close()

			if n == 0 {
				break
			}
			scriptfiles = scriptfiles + ",\n"
		}
		fmt.Println("scriptfiles=", scriptfiles)
		buf, _ := ioutil.ReadFile(jsondir + "provisioners.json")
                line =line+ ",\n"+string(buf)
			line = strings.Replace(line, "SCRIPTFILES", scriptfiles, -1)
			line = strings.Replace(line, "SSH_PASSWORD", r.Form.Get("password"), -1)
	}
        line =line+ "}"
        ioutil.WriteFile(newjson, []byte(line), 0)

	// new cfg file part
	var partitions string
	for k, v := range r.Form["part"] {
		sizen, _ := strconv.Atoi(r.Form["size"][k])
		sizens := strconv.Itoa(sizen * 1024)
		partitions = partitions + "part " + v + " --fstype=ext4 --size=" + sizens + "\n"
	}
        os.Create(newcfg)
	buf, _ = ioutil.ReadFile(cfg)
        line = string(buf)
		line = strings.Replace(line, "SSH_USERNAME", r.Form.Get("user"), -1)
		line = strings.Replace(line, "SSH_PASSWORD", r.Form.Get("password"), -1)
		line = strings.Replace(line, "PARTITIONS", partitions, -1)
		ioutil.WriteFile(newcfg, []byte(line), 0)
	

	fmt.Println(json)
	fmt.Println(cfg)
	fmt.Println(script)
	fmt.Println(iso)
	fmt.Println(newjson)
	return newjson
}
func callpacker(json string) *os.Process {
	fmt.Println("callpacker", json)
        inf, _ := os.Create(resultmap["jsondir"]+"inf.log")
        outf, _ := os.Create(resultmap["jsondir"]+"outf.log")
        errf, _ := os.Create(resultmap["jsondir"]+"errf.log")
	attr := &os.ProcAttr{
		//Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
                Files: []*os.File{inf, outf, errf},
	}
	p, err := os.StartProcess("/home/packerdir/packer", []string{"/home/packerdir/packer", "build", json}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("p=[",p,"]")
        return p
}
func calltransform(p *os.Process,output string,newoutput string) *os.Process {
         p.Wait()
 	fmt.Println("calltransform", output,newoutput)
        inf, _ := os.Create(resultmap["jsondir"]+"inf2.log")
        outf, _ := os.Create(resultmap["jsondir"]+"outf2.log")
        errf, _ := os.Create(resultmap["jsondir"]+"errf2.log")
	attr := &os.ProcAttr{
		//Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
                Files: []*os.File{inf, outf, errf},
	}
	p2, err := os.StartProcess("/bin/qemu-img", []string{"/bin/qemu-img", "convert","-f", "qcow2",output,"-O", "qcow2", "-o", "compat=0.10",newoutput}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("p2=[",p2,"]")
        return p2
}
