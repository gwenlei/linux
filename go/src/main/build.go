package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
        "bufio"
)

var conf map[string](map[string]string)
var dat map[string](map[string]string)
var reportlog map[string](map[string]string)
var ansibleroles map[string](map[string]string)
var ansiblelog map[string](map[string]string)

func main() {
	conf = make(map[string](map[string]string))
	dat = make(map[string](map[string]string))
	reportlog = make(map[string](map[string]string))
	buf, _ := ioutil.ReadFile("static/data/data.json")
	if len(buf) > 0 {
		if err := json.Unmarshal(buf, &conf); err != nil {
			panic(err)
		}
	}
	//初始化设定
	buf, _ = ioutil.ReadFile("static/data/dataqemu.json")
	if len(buf) > 0 {
		if err := json.Unmarshal(buf, &dat); err != nil {
			panic(err)
		}
	}
	buf, _ = ioutil.ReadFile("static/data/ansibleroles.json")
	if len(buf) > 0 {
		if err := json.Unmarshal(buf, &ansibleroles); err != nil {
			panic(err)
		}
	}
	buf, _ = ioutil.ReadFile("static/data/reportlog.json")
	if len(buf) > 0 {
		if err := json.Unmarshal(buf, &reportlog); err != nil {
			panic(err)
		}
	}
	buf, _ = ioutil.ReadFile("static/data/ansiblelog.json")
	if len(buf) > 0 {
		if err := json.Unmarshal(buf, &ansiblelog); err != nil {
			panic(err)
		}
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/build", build) //设置访问的路由
	http.HandleFunc("/setdat", setdat)
	http.HandleFunc("/setansible", setansible)
	http.HandleFunc("/report", report)
	http.HandleFunc("/upload", UploadServer)
        http.HandleFunc("/searchroles",searchroles)
        http.HandleFunc("/deleteroles",deleteroles)
	//err := http.ListenAndServe(conf["servermap"]["server"], nil) //设置监听的端口
	err := http.ListenAndServeTLS(conf["servermap"]["server"], "server.crt", "server.key", nil) //设置监听的端口
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
		//clear dat map
		for _, v := range dat {
			for j, _ := range v {
				delete(v, j)
			}
		}
		//reset dat map
		tmp := [...]string{"jsonmap", "cfgmap", "isomap", "md5map", "scriptmap", "resultmap", "servermap", "floppymap"}
		for _, vt := range tmp {
			for k, v := range r.Form[vt+"+fieldid"] {
				dat[vt][v] = r.Form[vt+"+fieldvalue"][k]
			}
		}
		filename := r.Form.Get("settingfile")
		backfile := "static/data/log/" + filename[strings.LastIndex(filename, "/")+1:] + time.Now().Format("20060102150405")
		newdataf, _ := os.Create(backfile)
		fmt.Println("filename:", filename)
		fmt.Println("backfile:", backfile)
		//dataf, _ := os.Open("static/data/data.json")
		dataf, _ := os.Open(filename)
		io.Copy(newdataf, dataf)
		defer newdataf.Close()
		defer dataf.Close()
		line, _ := json.Marshal(dat)
		ioutil.WriteFile(filename, line, 0)
		http.Redirect(w, r, "/setdat", 302)
	}
}

func setansible(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("ansible.html")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		for k, v := range r.Form {
			fmt.Println(k, ":", strings.Join(v, " "))
		}
		//clear ansibleroles map
		for _, v := range ansibleroles {
			for j, _ := range v {
				delete(v, j)
			}
		}
		ansibleroles["Ubuntu16.04"] = make(map[string]string)
		//reset dat map
		tmp := [...]string{"Ubuntu16.04"}
		for _, vt := range tmp {
			for k, v := range r.Form[vt+"+fieldid"] {
				ansibleroles[vt][v] = r.Form[vt+"+fieldvalue"][k]
			}
		}
		filename := "static/data/ansibleroles.json"
		backfile := "static/data/log/" + filename[strings.LastIndex(filename, "/")+1:] + time.Now().Format("20060102150405")
		newdataf, _ := os.Create(backfile)
		fmt.Println("filename:", filename)
		fmt.Println("backfile:", backfile)
		//dataf, _ := os.Open("static/data/data.json")
		dataf, _ := os.Open(filename)
		io.Copy(newdataf, dataf)
		defer newdataf.Close()
		defer dataf.Close()
		line, _ := json.Marshal(ansibleroles)
		ioutil.WriteFile(filename, line, 0)
		http.Redirect(w, r, "/setansible", 302)
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
			fmt.Println(k, ":", strings.Join(v, " "))
		}

		timest := buildjson(r)
		fmt.Println("buildjson end", timest)
		go callpacker(timest)
		go calltransform(timest)
		go callbzip2(timest)
		http.Redirect(w, r, "/build", 302)
	}
}

func report(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("report.html")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		for k, v := range r.Form {
			fmt.Println(k, ":", strings.Join(v, " "))
		}
		for _, v := range r.Form["clickt"] {
			delete(reportlog, v)
			if err := exec.Command("rm", "-rf", dat["resultmap"]["resultdir"]+v).Run(); err != nil {
				fmt.Printf("Error removing build directory: %s %s\n", dat["resultmap"]["resultdir"]+v, err)
			}

		}
		liner, _ := json.Marshal(reportlog)
		ioutil.WriteFile("static/data/reportlog.json", liner, 0)

		http.Redirect(w, r, "/report", 302)
	}
}

func buildjson(r *http.Request) (timest string) {
	timest = time.Now().Format("20060102150405")

	dat = make(map[string](map[string]string))
	bufst, _ := ioutil.ReadFile(r.Form.Get("settingfile"))
	if len(bufst) > 0 {
		if err := json.Unmarshal(bufst, &dat); err != nil {
			panic(err)
		}
	} else {
		fmt.Println("settingfile err")
		os.Exit(1)
	}
	//report
	reportlog[timest] = make(map[string]string)
	reportlog[timest]["resultdir"] = dat["resultmap"]["resultdir"] + timest + "/"
	reportlog[timest]["outputdir"] = reportlog[timest]["resultdir"] + "output/"
	reportlog[timest]["timestamp"] = timest
	reportlog[timest]["buildtype"] = r.Form.Get("buildtype")
	reportlog[timest]["ostype"] = r.Form.Get("ostype")
	reportlog[timest]["vmname"] = r.Form.Get("vmname")
	reportlog[timest]["user"] = r.Form.Get("user")
	reportlog[timest]["password"] = r.Form.Get("password")
	reportlog[timest]["disksize"] = r.Form.Get("disksize")
	reportlog[timest]["compat"] = r.Form.Get("compat")
	reportlog[timest]["headless"] = r.Form.Get("headless")
	reportlog[timest]["bzip2"] = r.Form.Get("bzip2")
	reportlog[timest]["settingfile"] = r.Form.Get("settingfile")
	reportlog[timest]["status"] = "waiting"
	if reportlog[timest]["buildtype"] == "qemu" {
		reportlog[timest]["downloadlink"] = reportlog[timest]["outputdir"] + reportlog[timest]["vmname"]
	} else {
		reportlog[timest]["downloadlink"] = reportlog[timest]["outputdir"] + reportlog[timest]["vmname"] + ".vhd"
	}
	if reportlog[timest]["bzip2"] == "Yes" {
		reportlog[timest]["downloadlink"] = reportlog[timest]["downloadlink"] + ".bz2"
	}
	fmt.Println("downloadlink=", reportlog[timest]["downloadlink"])
	for k, v := range r.Form["part"] {
		reportlog[timest]["part"] = reportlog[timest]["part"] + v + ":" + r.Form["size"][k] + " "
	}
	for _, v := range r.Form["software"] {
		reportlog[timest]["software"] = reportlog[timest]["software"] + v + "\n"
	}
	for _, v := range r.Form["ansible"] {
		reportlog[timest]["ansible"] = reportlog[timest]["ansible"] + v + " "
	}
	liner, _ := json.Marshal(reportlog)
	ioutil.WriteFile("static/data/reportlog.json", liner, 0)

	os.MkdirAll(reportlog[timest]["resultdir"], 0777)
	os.MkdirAll(reportlog[timest]["resultdir"]+"json/", 0777)
	os.MkdirAll(reportlog[timest]["resultdir"]+"script/", 0777)
	os.MkdirAll(reportlog[timest]["resultdir"]+"cfg/", 0777)
	tmplog := reportlog[timest]["resultdir"] + "form.log"
	os.Create(tmplog)
	var tmplogs string
	for k, v := range r.Form {
		tmplogs = tmplogs + k + ":"
		for _, v1 := range v {
			tmplogs = tmplogs + v1 + " "
		}
		tmplogs = tmplogs + "\n"
	}
	ioutil.WriteFile(tmplog, []byte(tmplogs), 0)

	fmt.Println("settingfile:" + reportlog[timest]["settingfile"])
	fmt.Println("settingfileback:" + reportlog[timest]["resultdir"] + reportlog[timest]["settingfile"][strings.LastIndex(reportlog[timest]["settingfile"], "/")+1:])
	CopyFile(reportlog[timest]["settingfile"], reportlog[timest]["resultdir"]+reportlog[timest]["settingfile"][strings.LastIndex(reportlog[timest]["settingfile"], "/")+1:])

	json := dat["jsonmap"][r.Form.Get("ostype")]
	reportlog[timest]["newjson"] = reportlog[timest]["resultdir"] + "json/" + json[strings.LastIndex(json, "/")+1:]
	cfg := dat["cfgmap"][r.Form.Get("ostype")]
	reportlog[timest]["newcfg"] = reportlog[timest]["resultdir"] + "cfg/" + cfg[strings.LastIndex(cfg, "/")+1:]
	//reportlog[timest]["newcfgs"] = "https://" + dat["servermap"]["server"] + "/" + reportlog[timest]["newcfg"]
	if index := strings.LastIndex(r.Form.Get("ostype"), "CentOS7"); index >= 0 {
		reportlog[timest]["newcfgs"] = reportlog[timest]["newcfg"][strings.LastIndex(reportlog[timest]["newcfg"], "/")+1:]
	} else if index := strings.LastIndex(r.Form.Get("ostype"), "Ubuntu16"); index >= 0 {
		reportlog[timest]["newcfgs"] = reportlog[timest]["newcfg"][strings.LastIndex(reportlog[timest]["newcfg"], "/")+1:]
	} else if index := strings.LastIndex(r.Form.Get("ostype"), "CentOS"); index >= 0 {
		reportlog[timest]["newcfgs"] = "floppy:/" + reportlog[timest]["newcfg"][strings.LastIndex(reportlog[timest]["newcfg"], "/")+1:]
	} else if index := strings.LastIndex(r.Form.Get("ostype"), "Ubuntu"); index >= 0 {
		reportlog[timest]["newcfgs"] = "/floppy/" + reportlog[timest]["newcfg"][strings.LastIndex(reportlog[timest]["newcfg"], "/")+1:]
	}
	reportlog[timest]["iso"] = dat["isomap"][r.Form.Get("ostype")]
	disksizen, _ := strconv.Atoi(r.Form.Get("disksize"))
	disksizens := strconv.Itoa(disksizen * 1024)
	//new json file
	os.Create(reportlog[timest]["newjson"])
	buf, _ := ioutil.ReadFile(json)
	line := string(buf)
	line = strings.Replace(line, "DISK_SIZE", disksizens, -1)
	line = strings.Replace(line, "SSH_USERNAME", r.Form.Get("user"), -1)
	line = strings.Replace(line, "SSH_PASSWORD", r.Form.Get("password"), -1)
	line = strings.Replace(line, "VM_NAME", r.Form.Get("vmname"), -1)
	line = strings.Replace(line, "OUTPUT_DIRECTORY", reportlog[timest]["resultdir"]+"output/", -1)
	line = strings.Replace(line, "ISO_CHECKSUM", dat["md5map"][dat["isomap"][r.Form.Get("ostype")]], -1)
	line = strings.Replace(line, "ISO_URL", reportlog[timest]["iso"], -1)
	line = strings.Replace(line, "FLOPPY_CFG", reportlog[timest]["newcfg"], -1)
	line = strings.Replace(line, "KS_CFG", reportlog[timest]["newcfgs"], -1)
	line = strings.Replace(line, "WIN_CFG", reportlog[timest]["newcfg"], -1)
	line = strings.Replace(line, "FLOPPYDIR", reportlog[timest]["resultdir"]+dat["floppymap"][r.Form.Get("ostype")][strings.LastIndex(dat["floppymap"][r.Form.Get("ostype")], "/")+1:], -1)
	line = strings.Replace(line, "CFGDIR", reportlog[timest]["resultdir"]+"cfg", -1)
	line = strings.Replace(line, "HEADLESS", r.Form.Get("headless"), -1)
	line = strings.Replace(line, "VHDDIR", reportlog[timest]["resultdir"]+"vhd/", -1)
	var script = make([]string, 10)
	var newscript = make([]string, 10)
	n := copy(script, r.Form["software"])
	copy(newscript, script)
	fmt.Println("n=", n)
	fmt.Println("len=", len(r.Form["software"]))
	var scriptfiles string
	if n > 0 || len(r.Form["ansible"]) > 0 {
		scriptfiles = ",\"provisioners\": [\n" +
			"{\n" +
			"\"type\": \"shell\",\n" +
			"\"execute_command\": \"echo 'SSH_PASSWORD' | {{.Vars}} sudo -S -E bash '{{.Path}}'\",\n" +
			"\"scripts\": [\n"
		for k, v := range r.Form["software"] {
			fmt.Println(k, v)

	f, err := os.Open(v)
	if err != nil {
		fmt.Println("open err")
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
		}
                newscript[k] = reportlog[timest]["resultdir"] + "script/" + line[:len(line)-1]
                scriptfiles = scriptfiles + "\"" + newscript[k] + "\","
			// copy script
			newscriptf, _ := os.Create(newscript[k])
                        fmt.Println("vscripts=",v[:strings.LastIndex(v,"/")+1]+line[:len(line)-1])
			scriptf, _ := os.Open(v[:strings.LastIndex(v,"/")+1]+line[:len(line)-1])
			io.Copy(newscriptf, scriptf)
			defer scriptf.Close()
			defer newscriptf.Close()
	}

		}
		if len(r.Form["ansible"]) > 0 {
			if len(r.Form["software"]) > 0 {
				scriptfiles = scriptfiles + ",\n"
			}
			scriptfiles = scriptfiles + "\"" + reportlog[timest]["resultdir"] + "script/ansible.sh" + "\""
			// copy script
			newscriptf, _ := os.Create(reportlog[timest]["resultdir"] + "script/ansible.sh")
			scriptf, _ := os.Open("template/script/ubuntu1604/ansible.sh")
			io.Copy(newscriptf, scriptf)
			defer scriptf.Close()
			defer newscriptf.Close()
		}
                // delete ,
                if strings.LastIndex(scriptfiles,",")==len(scriptfiles)-1 {
                   scriptfiles=scriptfiles[:len(scriptfiles)-1]
                }
		scriptfiles = scriptfiles + "]}\n"

		if len(r.Form["ansible"]) > 0 {
			ansiblefiles := ",{\n" +
				"\"type\": \"ansible-local\",\n" +
				"\"playbook_file\": \"" + reportlog[timest]["resultdir"] + "ansible/main.yml\",\n" +
				"\"role_paths\": [\n"
			ansiblemain := "---\n" + "- hosts: all\n" + "  roles:\n"
			n := len(r.Form["ansible"])
			for _, v := range r.Form["ansible"] {
				ansiblefiles = ansiblefiles + "\"/etc/ansible/roles/" + ansibleroles[r.Form.Get("ostype")][v] + "\""
				ansiblemain = ansiblemain + "    - " + ansibleroles[r.Form.Get("ostype")][v] + "\n"
				n = n - 1
				if n > 0 {
					ansiblefiles = ansiblefiles + ",\n"
					ansiblemain = ansiblemain + "\n"
				}
			}
			ansiblefiles = ansiblefiles + "]}"
			os.MkdirAll(reportlog[timest]["resultdir"]+"ansible/", 0777)
			os.Create(reportlog[timest]["resultdir"] + "ansible/main.yml")
			ioutil.WriteFile(reportlog[timest]["resultdir"]+"ansible/main.yml", []byte(ansiblemain), 0)
			scriptfiles = scriptfiles + ansiblefiles
		}
		scriptfiles = scriptfiles + "]"
		fmt.Println("scriptfiles=", scriptfiles)
	}
	line = strings.Replace(line, "SCRIPTFILES", scriptfiles, -1)

	ioutil.WriteFile(reportlog[timest]["newjson"], []byte(line), 0)

	// new cfg file part
	var partitions string
	if index := strings.LastIndex(r.Form.Get("ostype"), "CentOS"); index >= 0 {
		for k, v := range r.Form["part"] {
			sizen, _ := strconv.Atoi(r.Form["size"][k])
			sizens := strconv.Itoa(sizen * 1024)
			if v == "swap" {
				partitions = partitions + "part swap --size=" + sizens + "\n"
			} else {
				partitions = partitions + "part " + v + " --fstype=ext4 --size=" + sizens + "\n"
			}
		}
	} else if index := strings.LastIndex(r.Form.Get("ostype"), "Ubuntu"); index >= 0 {
		for k, v := range r.Form["part"] {
			sizen, _ := strconv.Atoi(r.Form["size"][k])
			sizens := strconv.Itoa(sizen * 1024)
			if k == 0 {
				partitions = partitions + "d-i partman-auto/method string regular\nd-i partman-auto/expert_recipe string boot-root :: "
				if v == "swap" {
					partitions = partitions + "64 " + sizens + " 300% $primary{ } linux-swap method{ swap } format{ } "
				} else if v == "/boot" {
					partitions = partitions + "64 " + sizens + " 200 ext4 $primary{ } $bootable{ } method{ format } format{ } use_filesystem{ } filesystem{ ext4 } mountpoint{ /boot } "
				} else {
					partitions = partitions + sizens + " 4000 -1 ext4 $primary{ } method{ format } format{ } use_filesystem{ } filesystem{ ext4 } mountpoint{ " + v + " } . "
				}
			} else {
				if v == "swap" {
					partitions = partitions + "64 " + sizens + " 300% linux-swap method{ swap } format{ } "
				} else if v == "/boot" {
					partitions = partitions + "64 " + sizens + " 200 ext4 $bootable{ } method{ format } format{ } use_filesystem{ } filesystem{ ext4 } mountpoint{ /boot } "
				} else {
					partitions = partitions + sizens + " 4000 -1 ext4 method{ format } format{ } use_filesystem{ } filesystem{ ext4 } mountpoint{ " + v + " } . "
				}
			}
		}
	} else if index := strings.LastIndex(r.Form.Get("ostype"), "OpenSuse"); index >= 0 {
		for k, v := range r.Form["part"] {
			partitions = partitions + "<partition><create config:type=\"boolean\">true</create><crypt_fs config:type=\"boolean\">false</crypt_fs><filesystem config:type=\"symbol\">btrfs</filesystem><format config:type=\"boolean\">true</format><loop_fs config:type=\"boolean\">false</loop_fs><mount>" + v + "</mount><mountby config:type=\"symbol\">device</mountby><partition_id config:type=\"integer\">" + strconv.Itoa(130+k+1) + "</partition_id><partition_nr config:type=\"integer\">" + strconv.Itoa(1+k) + "</partition_nr><raid_options/><resize config:type=\"boolean\">false</resize><size>" + r.Form["size"][k] + "G</size></partition>\n"
		}
	}
	var partitionadd string
	var partitionmodify string
	if index := strings.LastIndex(r.Form.Get("ostype"), "Windows"); index >= 0 {
		for k, v := range r.Form["part"] {
			sizen, _ := strconv.Atoi(r.Form["size"][k])
			sizens := strconv.Itoa(sizen * 1024)
			partitionadd = partitionadd + "<CreatePartition wcm:action=\"add\"><Order>" + strconv.Itoa(k+1) + "</Order><Type>Primary</Type><Extend>false</Extend><Size>" + sizens + "</Size></CreatePartition>\n"
			partitionmodify = partitionmodify + "<ModifyPartition wcm:action=\"add\"><Format>NTFS</Format><Label>" + r.Form.Get("ostype") + "</Label><Letter>" + v + "</Letter><Order>" + strconv.Itoa(k+1) + "</Order><PartitionID>" + strconv.Itoa(k+1) + "</PartitionID></ModifyPartition>\n"
		}
	}
	os.Create(reportlog[timest]["newcfg"])
	buf, _ = ioutil.ReadFile(cfg)
	line = string(buf)
	line = strings.Replace(line, "SSH_USERNAME", r.Form.Get("user"), -1)
	line = strings.Replace(line, "SSH_PASSWORD", r.Form.Get("password"), -1)
	line = strings.Replace(line, "PARTITIONS", partitions, -1)
	line = strings.Replace(line, "PARTITONADD", partitionadd, -1)
	line = strings.Replace(line, "PARTITONMODIFY", partitionmodify, -1)
	ioutil.WriteFile(reportlog[timest]["newcfg"], []byte(line), 0)
	if index := strings.LastIndex(r.Form.Get("ostype"), "Windows"); index >= 0 {
		CopyDir(dat["floppymap"][r.Form.Get("ostype")], reportlog[timest]["resultdir"]+dat["floppymap"][r.Form.Get("ostype")][strings.LastIndex(dat["floppymap"][r.Form.Get("ostype")], "/")+1:])
	}

	fmt.Println(reportlog[timest]["newjson"])
	return timest
}
func callpacker(timest string) *os.Process {
	fmt.Println("callpacker", timest)
	outf, _ := os.Create(reportlog[timest]["resultdir"] + "packer.log")
	attr := &os.ProcAttr{
		//Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Files: []*os.File{outf, outf, outf},
	}
	p, err := os.StartProcess(dat["servermap"]["packer"], []string{dat["servermap"]["packer"], "build", reportlog[timest]["newjson"]}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("p=[", p, "]")
	reportlog[timest]["packerpid"] = strconv.Itoa(p.Pid)
	go checkstatus(p, "packer", timest)
	return p
}
func calltransform(timest string) {
	if reportlog[timest]["compat"] != "0.1" {
		fmt.Printf("compat:No\n")
		return
	}
	for {
		if reportlog[timest]["status"] == "packer failed" {
			fmt.Printf("compat:packer failed\n")
			return
		} else if reportlog[timest]["status"] == "packer success" {
			break
		} else {
			fmt.Printf("compat sleep 2m\n")
			time.Sleep(120 * time.Second)
		}
	}

	fmt.Println("calltransform")
	outf, _ := os.Create(reportlog[timest]["resultdir"] + "convert.log")
	attr := &os.ProcAttr{
		//Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Files: []*os.File{outf, outf, outf},
	}
	output := reportlog[timest]["resultdir"] + "output/" + reportlog[timest]["vmname"]
	newoutput := reportlog[timest]["resultdir"] + "output/tr" + reportlog[timest]["vmname"]
	fmt.Println("output=[", output, "]")
	fmt.Println("newoutput=[", newoutput, "]")
	p2, err := os.StartProcess("/bin/qemu-img", []string{"/bin/qemu-img", "convert", "-f", "qcow2", output, "-O", "qcow2", "-o", "compat=0.10", newoutput}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("p2=[", p2, "]")
	reportlog[timest]["transformpid"] = strconv.Itoa(p2.Pid)
	go checkstatus(p2, "transform", timest)

}

func callbzip2(timest string) {
	if reportlog[timest]["bzip2"] != "Yes" {
		fmt.Printf("bzip2:No\n")
		return
	}
	for {
		if reportlog[timest]["status"] == "packer failed" {
			fmt.Printf("compat:packer failed\n")
			return
		} else if reportlog[timest]["status"] == "transform failed" {
			fmt.Printf("compat:transform failed\n")
			return
		} else if reportlog[timest]["compat"] == "0.1" && reportlog[timest]["status"] == "transform success" {
			break
		} else if reportlog[timest]["compat"] != "0.1" && reportlog[timest]["status"] == "packer success" {
			break
		} else {
			fmt.Printf("bzip2 sleep 2m\n")
			time.Sleep(120 * time.Second)
		}
	}
	var vmstr string
	if reportlog[timest]["buildtype"] == "qemu" {
		vmstr = reportlog[timest]["outputdir"] + reportlog[timest]["vmname"]
	} else {
		vmstr = reportlog[timest]["outputdir"] + reportlog[timest]["vmname"] + ".vhd"
	}
	fmt.Printf("bzip2:" + vmstr + "\n")
	reportlog[timest]["status"] = "bzip2 running"
	liner, _ := json.Marshal(reportlog)
	ioutil.WriteFile("static/data/reportlog.json", liner, 0)
	cmd := exec.Command("bzip2", "-z", vmstr)
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("bzip2 end:" + vmstr + "\n")
	reportlog[timest]["status"] = "bzip2 success"
	liner, _ = json.Marshal(reportlog)
	ioutil.WriteFile("static/data/reportlog.json", liner, 0)
}
func checkstatus(p *os.Process, pname string, timest string) bool {
	fmt.Println("checkstatus", pname, p)
	reportlog[timest]["status"] = pname + " running"
	reportlog[timest][pname+"start"] = time.Now().Format("20060102150405")
	liner, _ := json.Marshal(reportlog)
	ioutil.WriteFile("static/data/reportlog.json", liner, 0)
	pw, _ := p.Wait()
	fmt.Println("checkstatus over", p)
	fmt.Println("timest=", timest)
	reportlog[timest][pname+"stop"] = time.Now().Format("20060102150405")
	t1, _ := time.Parse("20060102150405", reportlog[timest][pname+"stop"])
	t2, _ := time.Parse("20060102150405", reportlog[timest][pname+"start"])
	reportlog[timest][pname+"time"] = strconv.Itoa(int(t1.Sub(t2)) / 1e9)
	fmt.Println("t1=", t1)
	fmt.Println("t2=", t2)
	fmt.Println("cost=", t1.Sub(t2))
	status := pw.Success()
	if status == true {
		reportlog[timest]["status"] = pname + " success"
		fmt.Println("checkstatus over success ", pname, p)
	} else {
		reportlog[timest]["status"] = pname + " failed"
		fmt.Println("checkstatus over failed ", pname, p)
	}
	liner, _ = json.Marshal(reportlog)
	ioutil.WriteFile("static/data/reportlog.json", liner, 0)
	return status

}

func CopyFile(source string, dest string) (err error) {
	sourcefile, err := os.Open(source)
	if err != nil {
		return err
	}

	defer sourcefile.Close()

	destfile, err := os.Create(dest)
	if err != nil {
		return err
	}

	defer destfile.Close()

	_, err = io.Copy(destfile, sourcefile)
	if err == nil {
		sourceinfo, err := os.Stat(source)
		if err != nil {
			err = os.Chmod(dest, sourceinfo.Mode())
		}

	}

	return
}

func CopyDir(source string, dest string) (err error) {

	// get properties of source dir
	sourceinfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	// create dest dir

	err = os.MkdirAll(dest, sourceinfo.Mode())
	if err != nil {
		return err
	}

	directory, _ := os.Open(source)

	objects, err := directory.Readdir(-1)

	for _, obj := range objects {

		sourcefilepointer := source + "/" + obj.Name()

		destinationfilepointer := dest + "/" + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			err = CopyDir(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		} else {
			// perform copy
			err = CopyFile(sourcefilepointer, destinationfilepointer)
			if err != nil {
				fmt.Println(err)
			}
		}

	}
	return
}

func UploadServer(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("文件上传异常")
		}
	}()
	if r.Method == "GET" {
		t, _ := template.ParseFiles("upload.html")
		t.Execute(w, nil)
	} else {
		r.ParseMultipartForm(32 << 20) //在使用r.MultipartForm前必须先调用ParseMultipartForm方法，参数为最大缓存
		// fmt.Println(r.MultipartForm)
		// fmt.Println(r.MultipartReader())
		if r.MultipartForm != nil && r.MultipartForm.File != nil {
			fhs := r.MultipartForm.File["userfile"] //获取所有上传文件信息
			num := len(fhs)

			fmt.Printf("total：%d files", num)

			//循环对每个文件进行处理
			for n, fheader := range fhs {
				//获取文件名
				filename := fheader.Filename

				//结束文件
				file, err := fheader.Open()
				if err != nil {
					fmt.Println(err)
				}

				//保存文件
				defer file.Close()
				f, err := os.Create("static/upload/" + filename)
				defer f.Close()
				io.Copy(f, file)

				//获取文件状态信息
				fstat, _ := f.Stat()

				//打印接收信息
				fmt.Fprintf(w, "%s  NO.: %d  Size: %d KB  Name：%s\n", time.Now().Format("2006-01-02 15:04:05"), n, fstat.Size()/1024, filename)
				fmt.Printf("%s  NO.: %d  Size: %d KB  Name：%s\n", time.Now().Format("2006-01-02 15:04:05"), n, fstat.Size()/1024, filename)

			}
		}

		return
	}

}

func searchroles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("search.html")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		for k, v := range r.Form {
			fmt.Println(k, ":", strings.Join(v, " "))
		}

		timest := time.Now().Format("20060102150405")
                ansiblelog[timest] = make(map[string]string)
                ansiblelog[timest]["resultdir"] = "static/result/ansible/"+timest+"/"
                os.MkdirAll(ansiblelog[timest]["resultdir"], 0777)
                ansiblelog[timest]["ostype"]=r.Form.Get("ostype")
                ansiblelog[timest]["backingfile"]=r.Form.Get("backingfile")
                ansiblelog[timest]["times"]=r.Form.Get("times")
                ansiblelog[timest]["status"]="running"
                ansiblelog[timest]["rolename"]=r.Form.Get("rolename")
                ansiblelog[timest]["timestamp"]=timest
                ansiblelog[timest]["log"]=ansiblelog[timest]["resultdir"]+"ansible.log"
                liner, _ := json.Marshal(ansiblelog)
	        ioutil.WriteFile("static/data/ansiblelog.json", liner, 0)

  
		go callsearch(timest)
		go callfilter(timest)
		http.Redirect(w, r, "/searchroles", 302)
	}
}

func callsearch(timest string) {
	fmt.Println("callsearch", timest)
	outf, _ := os.Create(ansiblelog[timest]["resultdir"] + "search.log")
	attr := &os.ProcAttr{
		//Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Files: []*os.File{outf, outf, outf},
	}
	p, err := os.StartProcess("/usr/bin/ansible-galaxy",[]string{"/usr/bin/ansible-galaxy", "search", ansiblelog[timest]["rolename"]}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("p=[", p, "]")
	ansiblelog[timest]["searchpid"] = strconv.Itoa(p.Pid)
	go checkansiblestatus(p, "search", timest)
}

func checkansiblestatus(p *os.Process, pname string, timest string) {
	fmt.Println("checkstatus", pname, p)
	ansiblelog[timest]["status"] = pname + " running"
	ansiblelog[timest][pname+"start"] = time.Now().Format("20060102150405")
	liner, _ := json.Marshal(ansiblelog)
	ioutil.WriteFile("static/data/ansiblelog.json", liner, 0)
	pw, _ := p.Wait()
	fmt.Println("checkstatus over", p)
	fmt.Println("timest=", timest)
	ansiblelog[timest][pname+"stop"] = time.Now().Format("20060102150405")
	t1, _ := time.Parse("20060102150405", ansiblelog[timest][pname+"stop"])
	t2, _ := time.Parse("20060102150405", ansiblelog[timest][pname+"start"])
	ansiblelog[timest][pname+"time"] = strconv.Itoa(int(t1.Sub(t2)) / 1e9)
	fmt.Println("t1=", t1)
	fmt.Println("t2=", t2)
	fmt.Println("cost=", t1.Sub(t2))
	status := pw.Success()
	if status == true {
		ansiblelog[timest]["status"] = pname + " success"
		fmt.Println("checkstatus over success ", pname, p)
	} else {
		ansiblelog[timest]["status"] = pname + " failed"
		fmt.Println("checkstatus over failed ", pname, p)
	}
	liner, _ = json.Marshal(ansiblelog)
	ioutil.WriteFile("static/data/ansiblelog.json", liner, 0)
}

func callfilter(timest string) {
	for {
		if ansiblelog[timest]["status"] == "search failed" {
			fmt.Printf("search failed\n")
			return
		} else if ansiblelog[timest]["status"] == "search success" {
			break
		}else {
			fmt.Printf("callfilter sleep 2m\n")
			time.Sleep(120 * time.Second)
		}
	}
        ansiblelog[timest]["status"] = "filter running"
        ansiblelog[timest]["filterstart"] = time.Now().Format("20060102150405")
	liner, _ := json.Marshal(ansiblelog)
	ioutil.WriteFile("static/data/ansiblelog.json", liner, 0)

	f, err := os.Open(ansiblelog[timest]["resultdir"] + "search.log")
	if err != nil {
		fmt.Println("open err")
	}
	buf := bufio.NewReader(f)
	n := 0
	rolelist := []string{}
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		n = n + 1
		if n < 6 {
			continue
		}
	        //fmt.Println("line=", line)
                l :=len(strings.Split(line, " "))
                //fmt.Println("linelength=",l)
                if l>1{
		   rolelist = append(rolelist, strings.Split(line, " ")[1])
                }
	}
	fmt.Println("rolelist sum", len(rolelist))


	os.Create(ansiblelog[timest]["resultdir"]+"vm.xml")
        ansiblelog[timest]["sourcefile"]="/home/code/mycode/go/src/main/static/result/ansible/"+timest+"/vm.qcow2"
        ansiblelog[timest]["vmname"]="filtervm"+timest
	buf2, _ := ioutil.ReadFile("template/vm/vm.xml")
	line := string(buf2)
	line = strings.Replace(line, "VMNAME", ansiblelog[timest]["vmname"], -1)
	line = strings.Replace(line, "SOURCEFILE", ansiblelog[timest]["sourcefile"], -1)
	ioutil.WriteFile(ansiblelog[timest]["resultdir"]+"vm.xml", []byte(line), 0)

        outf, _ := os.Create(ansiblelog[timest]["resultdir"]+"ansible.log")
	attr := &os.ProcAttr{
		Files: []*os.File{outf, outf, outf},
	}
	p, err := os.StartProcess("/usr/bin/virsh", []string{"/usr/bin/virsh", "shutdown", ansiblelog[timest]["vmname"]}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ := p.Wait()
	fmt.Println("virsh shutdown ", ansiblelog[timest]["vmname"],pw.Success())
        time.Sleep(60 * time.Second)

	p, err = os.StartProcess("/usr/bin/qemu-img", []string{"/usr/bin/qemu-img", "create", "-f", "qcow2", ansiblelog[timest]["sourcefile"], "-b", ansiblelog[timest]["backingfile"]}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ = p.Wait()
	fmt.Println("qemu-img end ",pw.Success())
        time.Sleep(30 * time.Second)

	p, err = os.StartProcess("/usr/bin/virsh", []string{"/usr/bin/virsh", "define", ansiblelog[timest]["resultdir"]+"vm.xml"}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ = p.Wait()
	fmt.Println("virsh define ", ansiblelog[timest]["vmname"],pw.Success())

	p, err = os.StartProcess("/usr/bin/virsh", []string{"/usr/bin/virsh", "start", ansiblelog[timest]["vmname"]}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ = p.Wait()
	fmt.Println("virsh start ",pw.Success())
        time.Sleep(60 * time.Second)

        outf2, _ := os.Create(ansiblelog[timest]["resultdir"]+"mac.log")
	attr2 := &os.ProcAttr{
		Files: []*os.File{outf2, outf2, outf2},
	}
	p, err = os.StartProcess("/usr/bin/virsh", []string{"/usr/bin/virsh", "domiflist", ansiblelog[timest]["vmname"]}, attr2)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ = p.Wait()
	fmt.Println("virsh domiflist ",ansiblelog[timest]["vmname"],pw.Success())

	f, err = os.Open(ansiblelog[timest]["resultdir"] + "mac.log")
	if err != nil {
		fmt.Println("open err")
	}
	buf = bufio.NewReader(f)
	n = 0
        var mac string
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		n = n + 1
		if index := strings.LastIndex(line, "default"); index >= 0 {
			mac = line[strings.LastIndex(line, " ")+1 : len(line)-1]
                        mac = strings.Replace(mac, " ", "", -1)
                        break
		}
		
	}
	fmt.Println("mac=", mac)


        outf3, _ := os.Create(ansiblelog[timest]["resultdir"]+"ip.log")
	attr3 := &os.ProcAttr{
		Files: []*os.File{outf3, outf3, outf3},
	}
	p, err = os.StartProcess("/usr/bin/virsh", []string{"/usr/bin/virsh", "net-dhcp-leases", "default"}, attr3)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ = p.Wait()
	fmt.Println("virsh net-dhcp-leases default ",pw.Success())

	f, err = os.Open(ansiblelog[timest]["resultdir"] + "ip.log")
	if err != nil {
		fmt.Println("open err")
	}
	buf = bufio.NewReader(f)
        var ip string
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		if index := strings.LastIndex(line, mac); index >= 0 {
                        ip = line[strings.LastIndex(line, "ipv4")+4 : strings.LastIndex(line, "/")] 
                        ip = strings.Replace(ip, " ", "", -1) 
                        break
		}
		
	}
	fmt.Println("ip=", ip)

        p, err = os.StartProcess("/usr/bin/sed", []string{"/usr/bin/sed", "-i", "/" + ip + "/d", "/root/.ssh/known_hosts"}, attr)
        if err != nil {
                fmt.Println(err.Error())
        }
        fmt.Println(p)
        pw, _ = p.Wait()
        fmt.Println("sed ", ip, pw.Success())

        outf4, _ := os.Create(ansiblelog[timest]["resultdir"]+"keyscan.log")
	attr4 := &os.ProcAttr{
		Files: []*os.File{outf4, outf4, outf4},
	}
        p, err = os.StartProcess("/usr/bin/ssh-keyscan", []string{"/usr/bin/ssh-keyscan", ip}, attr4)
        if err != nil {
                fmt.Println(err.Error())
        }
        fmt.Println(p)
        pw, _ = p.Wait()
        fmt.Println("ssh-keyscan ", ip, pw.Success())

	buf3, _ := ioutil.ReadFile(ansiblelog[timest]["resultdir"]+"keyscan.log")
	line3 := string(buf3)
	buf4, _ := ioutil.ReadFile("/root/.ssh/known_hosts")
	line4 := string(buf4)
	line4 = line4+"\n"+line3
	ioutil.WriteFile("/root/.ssh/known_hosts", []byte(line4), 0)


	n = 0
        times, _ := strconv.Atoi(ansiblelog[timest]["times"])
     for _, v := range rolelist {
        n=n+1
        if n> times{
           break
        }
        fmt.Println("n=",n)
	p, err := os.StartProcess("/usr/bin/ansible-galaxy", []string{"/usr/bin/ansible-galaxy", "install", v}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ := p.Wait()
	fmt.Println("ansible-galaxy install ", v,pw.Success())

	fout, err := os.Create(ansiblelog[timest]["resultdir"]+"main.yml")
	if err != nil {
		fmt.Println("main.yml", err)
	}
	fout.WriteString("- hosts: all\n  roles:\n  - " + v)
	fout.Close()
	fout, err = os.Create(ansiblelog[timest]["resultdir"]+"host")
	if err != nil {
		fmt.Println("host", err)
	}
	fout.WriteString(ip+" ansible_ssh_user=root ansible_ssh_pass=engine")
	fout.Close()

	p, err = os.StartProcess("/usr/bin/virsh", []string{"/usr/bin/virsh", "shutdown", ansiblelog[timest]["vmname"]}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ = p.Wait()
	fmt.Println("virsh shutdown ", ansiblelog[timest]["vmname"],pw.Success())
        time.Sleep(60 * time.Second)

	p, err = os.StartProcess("/usr/bin/qemu-img", []string{"/usr/bin/qemu-img", "create", "-f", "qcow2", ansiblelog[timest]["sourcefile"], "-b", ansiblelog[timest]["backingfile"]}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ = p.Wait()
	fmt.Println("qemu-img end ",pw.Success())
        time.Sleep(30 * time.Second)

	p, err = os.StartProcess("/usr/bin/virsh", []string{"/usr/bin/virsh", "start", ansiblelog[timest]["vmname"]}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ = p.Wait()
	fmt.Println("virsh start ",pw.Success())
        time.Sleep(60 * time.Second)

	p, err = os.StartProcess("/usr/bin/ansible-playbook", []string{"/usr/bin/ansible-playbook", "-i", ansiblelog[timest]["resultdir"]+"host",ansiblelog[timest]["resultdir"]+"main.yml"}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ = p.Wait()
	fmt.Println("ansible-playbook end ",v,pw.Success())
        if pw.Success(){
           ansiblelog[timest]["successroles"] = ansiblelog[timest]["successroles"]+" "+v
	   liner, _ := json.Marshal(ansiblelog)
	   ioutil.WriteFile("static/data/ansiblelog.json", liner, 0)
        }
   }


	p, err = os.StartProcess("/usr/bin/virsh", []string{"/usr/bin/virsh", "shutdown", ansiblelog[timest]["vmname"]}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ = p.Wait()
	fmt.Println("virsh shutdown ", ansiblelog[timest]["vmname"],pw.Success())
        time.Sleep(60 * time.Second)

	p, err = os.StartProcess("/usr/bin/virsh", []string{"/usr/bin/virsh", "undefine", ansiblelog[timest]["vmname"]}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(p)
	pw, _ = p.Wait()
	fmt.Println("virsh undefine ", ansiblelog[timest]["vmname"],pw.Success())

        ansiblelog[timest]["filterstop"] = time.Now().Format("20060102150405")
	t1, _ := time.Parse("20060102150405", ansiblelog[timest]["filterstop"])
	t2, _ := time.Parse("20060102150405", ansiblelog[timest]["filterstart"])
	ansiblelog[timest]["filtertime"] = strconv.Itoa(int(t1.Sub(t2)) / 1e9)
        ansiblelog[timest]["status"] = "filter finish"
	liner, _ = json.Marshal(ansiblelog)
	ioutil.WriteFile("static/data/ansiblelog.json", liner, 0)

}

func deleteroles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("search.html")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		for k, v := range r.Form {
			fmt.Println(k, ":", strings.Join(v, " "))
		}
		for _, v := range r.Form["clickt"] {
			delete(ansiblelog, v)
			if err := exec.Command("rm", "-rf", ansiblelog[v]["resultdir"]).Run(); err != nil {
				fmt.Printf("Error removing build directory: %s %s\n", ansiblelog[v]["resultdir"], err)
			}

		}
		liner, _ := json.Marshal(ansiblelog)
		ioutil.WriteFile("static/data/ansiblelog.json", liner, 0)

		http.Redirect(w, r, "/deleteroles", 302)
	}
}

