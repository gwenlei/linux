package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
        "time"
"encoding/json"
)
var dat map[string](map[string]string)

func main() {
        buf, _ := ioutil.ReadFile("static/data/data.json")
        if err := json.Unmarshal(buf, &dat); err != nil {
                panic(err)
        }
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/build", build) //设置访问的路由
	http.HandleFunc("/setdat", setdat)
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("index.html")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		for k, v := range r.Form {
			fmt.Println(k, ":", strings.Join(v, " "))
		}
	}
}

func setdat(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("dat.html")
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
		p := callpacker(json)
		//fmt.Fprintf(w, strconv.Itoa(p.Pid)+" packer running")
		if r.Form.Get("compat") == "0.1" {
			output := dat["resultmap"]["outputdir"] + dat["resultmap"]["outputimage"]
			newoutput := dat["resultmap"]["outputdir"] + "tr" + dat["resultmap"]["outputimage"]
			go calltransform(p, output, newoutput)
		}
                http.Redirect(w, r, "/build", 302)
	}
}

func buildjson(r *http.Request) (result string) {
        dat["resultmap"]["jsondir"]="static/result/"+time.Now().Format("20060102150405")+"/"
        dat["resultmap"]["outputdir"]=dat["resultmap"]["jsondir"]+"output/"
        dat["resultmap"]["scriptdir"]=dat["resultmap"]["jsondir"]+"script/"
        dat["resultmap"]["cfgdir"]=dat["resultmap"]["jsondir"]+"cfg/"
	os.MkdirAll(dat["resultmap"]["jsondir"], 0777)
	os.MkdirAll(dat["resultmap"]["scriptdir"], 0777)
	os.MkdirAll(dat["resultmap"]["cfgdir"], 0777)
	jsondir := dat["dirmap"]["json"]
	json := dat["jsonmap"][r.Form.Get("ostype")]
	newjson := dat["resultmap"]["jsondir"] + json
	json = jsondir + json
	cfgdir := dat["dirmap"]["cfg"]
	cfg := dat["cfgmap"][r.Form.Get("ostype")]
	newcfg := dat["resultmap"]["cfgdir"] + cfg
	newcfgs := dat["resultmap"]["http"] + newcfg
	cfg = cfgdir + cfg
	isodir := dat["dirmap"]["iso"]
	iso := dat["isomap"][r.Form.Get("ostype")]
	iso = isodir + iso
	disksizen, _ := strconv.Atoi(r.Form.Get("disksize"))
	disksizens := strconv.Itoa(disksizen * 1024)
	//new json file
	os.Create(newjson)
	buf, _ := ioutil.ReadFile(json)
	line := string(buf)
	line = strings.Replace(line, "DISK_SIZE", disksizens, -1)
	line = strings.Replace(line, "SSH_USERNAME", r.Form.Get("user"), -1)
	line = strings.Replace(line, "SSH_PASSWORD", r.Form.Get("password"), -1)
	line = strings.Replace(line, "VM_NAME", dat["resultmap"]["outputimage"], -1)
	line = strings.Replace(line, "OUTPUT_DIRECTORY", dat["resultmap"]["outputdir"], -1)
	line = strings.Replace(line, "ISO_CHECKSUM", dat["md5map"][dat["isomap"][r.Form.Get("ostype")]], -1)
	line = strings.Replace(line, "ISO_URL", iso, -1)
	line = strings.Replace(line, "KS_CFG", newcfgs, -1)
	line = strings.Replace(line, "HEADLESS", r.Form.Get("headless"), -1)
	scriptdir := dat["dirmap"]["script"]
	var script = make([]string, 10)
	var newscript = make([]string, 10)
	n := copy(script, r.Form["software"])
	copy(newscript, script)
	fmt.Println("n=", n)
	if n > 0 {
		var scriptfiles = "\"" + dat["resultmap"]["scriptdir"] + "base.sh" + "\",\n"
		// copy script
		newbasescriptf, _ := os.Create(dat["resultmap"]["scriptdir"] + "base.sh")
		basescriptf, _ := os.Open(scriptdir + "base.sh")
		io.Copy(newbasescriptf, basescriptf)
		defer basescriptf.Close()
		defer newbasescriptf.Close()
		for k, v := range script {
			fmt.Println(k, v)
			script[k] = scriptdir + dat["scriptmap"][v]
			newscript[k] = dat["resultmap"]["scriptdir"] + dat["scriptmap"][v]
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
		line = line + ",\n" + string(buf)
		line = strings.Replace(line, "SCRIPTFILES", scriptfiles, -1)
		line = strings.Replace(line, "SSH_PASSWORD", r.Form.Get("password"), -1)
	}
	line = line + "}"
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
	inf, _ := os.Create(dat["resultmap"]["jsondir"] + "inf.log")
	outf, _ := os.Create(dat["resultmap"]["jsondir"] + "outf.log")
	errf, _ := os.Create(dat["resultmap"]["jsondir"] + "errf.log")
	attr := &os.ProcAttr{
		//Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Files: []*os.File{inf, outf, errf},
	}
	p, err := os.StartProcess("/home/packerdir/packer", []string{"/home/packerdir/packer", "build", json}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("p=[", p, "]")
	return p
}
func calltransform(p *os.Process, output string, newoutput string) *os.Process {
	p.Wait()
	fmt.Println("calltransform", output, newoutput)
	inf, _ := os.Create(dat["resultmap"]["jsondir"] + "inf2.log")
	outf, _ := os.Create(dat["resultmap"]["jsondir"] + "outf2.log")
	errf, _ := os.Create(dat["resultmap"]["jsondir"] + "errf2.log")
	attr := &os.ProcAttr{
		//Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Files: []*os.File{inf, outf, errf},
	}
	p2, err := os.StartProcess("/bin/qemu-img", []string{"/bin/qemu-img", "convert", "-f", "qcow2", output, "-O", "qcow2", "-o", "compat=0.10", newoutput}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("p2=[", p2, "]")
	return p2
}
