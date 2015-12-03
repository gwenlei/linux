package main

import (
	"encoding/json"
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
        "path/filepath"
        "os/exec"
)

var dat map[string](map[string]string)
var reportlog map[string](map[string]string)

func main() {
        dat=make(map[string](map[string]string))
        reportlog=make(map[string](map[string]string))
	buf, _ := ioutil.ReadFile("static/data/data.json")
        if len(buf)>0{
	if err := json.Unmarshal(buf, &dat); err != nil {
		panic(err)
	}
        }
	buf, _ = ioutil.ReadFile("static/data/reportlog.json")
        if len(buf)>0{
	if err := json.Unmarshal(buf, &reportlog); err != nil {
		panic(err)
	}
        }
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/build", build) //设置访问的路由
	http.HandleFunc("/setdat", setdat)
        http.HandleFunc("/report", report)
	err := http.ListenAndServe(dat["servermap"]["server"], nil) //设置监听的端口
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
		tmp := [...]string{"jsonmap", "cfgmap", "isomap", "md5map", "scriptmap", "resultmap", "servermap"}
		for _, vt := range tmp {
			for k, v := range r.Form[vt+"+fieldid"] {
				dat[vt][v] = r.Form[vt+"+fieldvalue"][k]
			}
		}
		newdataf, _ := os.Create("static/data/log/data.json" + time.Now().Format("20060102150405"))
		dataf, _ := os.Open("static/data/data.json")
		io.Copy(newdataf, dataf)
		defer newdataf.Close()
		defer dataf.Close()
		line, _ := json.Marshal(dat)
		ioutil.WriteFile("static/data/data.json", line, 0)
		http.Redirect(w, r, "/setdat", 302)
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
		p := callpacker(timest)
                go calltransform(p, timest)
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
                        if err := exec.Command("rm", "-rf",dat["resultmap"]["resultdir"] +v).Run(); err != nil {
                           fmt.Printf("Error removing build directory: %s %s\n", dat["resultmap"]["resultdir"] +v,err)
                        }
		}
                liner, _ := json.Marshal(reportlog)
	        ioutil.WriteFile("static/data/reportlog.json", liner, 0)
                
                http.Redirect(w, r, "/report", 302)
	}
}

func buildjson(r *http.Request) (timest string) {
        timest=string(time.Now().Format("20060102150405"))
        
        //report
        reportlog[timest]=make(map[string]string)
        reportlog[timest]["resultdir"]= dat["resultmap"]["resultdir"] + timest + "/"
        reportlog[timest]["timestamp"]=timest
        reportlog[timest]["ostype"]=r.Form.Get("ostype")
        reportlog[timest]["vmname"]=r.Form.Get("vmname")
        reportlog[timest]["user"]=r.Form.Get("user")
        reportlog[timest]["password"]=r.Form.Get("password")
        reportlog[timest]["disksize"]=r.Form.Get("disksize")
        reportlog[timest]["compat"]=r.Form.Get("compat")
        reportlog[timest]["outputdir"]=reportlog[timest]["resultdir"]+"output/"
        reportlog[timest]["status"]="waiting"
        for k, v := range r.Form["part"] {
             reportlog[timest]["part"]=reportlog[timest]["part"]+v+":"+r.Form["size"][k]+" "
	}
        for _, v := range r.Form["software"] {
             reportlog[timest]["software"]=reportlog[timest]["software"]+v+"\n"
	}   
   	liner, _ := json.Marshal(reportlog)
        ioutil.WriteFile("static/data/reportlog.json", liner, 0)
	
	os.MkdirAll(reportlog[timest]["resultdir"], 0777)
	os.MkdirAll(reportlog[timest]["resultdir"]+"json/", 0777)
	os.MkdirAll(reportlog[timest]["resultdir"]+"script/", 0777)
	os.MkdirAll(reportlog[timest]["resultdir"]+"cfg/", 0777)
        tmplog:=reportlog[timest]["resultdir"]+"form.log"
        os.Create(tmplog)
        var tmplogs string
        for k, v := range r.Form {
	    tmplogs=tmplogs+k+":"
            for _, v1 := range v {
	       tmplogs=tmplogs+v1+" "
	    }
            tmplogs=tmplogs+"\n"
	}
        ioutil.WriteFile(tmplog, []byte(tmplogs), 0)
	json := dat["jsonmap"][r.Form.Get("ostype")]
	reportlog[timest]["newjson"]= reportlog[timest]["resultdir"]+"json/" + json[strings.LastIndex(json, "/")+1:]
	cfg := dat["cfgmap"][r.Form.Get("ostype")]
	reportlog[timest]["newcfg"]= reportlog[timest]["resultdir"]+"cfg/" + cfg[strings.LastIndex(cfg, "/")+1:]
	reportlog[timest]["newcfgs"] = "http://"+dat["servermap"]["server"]+"/" + reportlog[timest]["newcfg"]
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
	line = strings.Replace(line, "KS_CFG", reportlog[timest]["newcfgs"], -1)
        line = strings.Replace(line, "WIN_CFG", reportlog[timest]["newcfg"], -1)
        line = strings.Replace(line, "RESULTDIR", reportlog[timest]["resultdir"], -1)
	line = strings.Replace(line, "HEADLESS", r.Form.Get("headless"), -1)
	var script = make([]string, 10)
	var newscript = make([]string, 10)
	n := copy(script, r.Form["software"])
	copy(newscript, script)
	fmt.Println("n=", n)
	if n > 0 {
		var scriptfiles string
		for k, v := range script {
			fmt.Println(k, v)
			newscript[k] = reportlog[timest]["resultdir"]+"script/" + v[strings.LastIndex(v, "/")+1:]
			scriptfiles = scriptfiles + "\"" + newscript[k] + "\""
			n = n - 1
			// copy script
			newscriptf, _ := os.Create(newscript[k])
			scriptf, _ := os.Open(v)
			io.Copy(newscriptf, scriptf)
			defer scriptf.Close()
			defer newscriptf.Close()

			if n == 0 {
				break
			}
			scriptfiles = scriptfiles + ",\n"
		}
		fmt.Println("scriptfiles=", scriptfiles)
		buf, _ := ioutil.ReadFile("template/json/provisioners.json")
		line = line + ",\n" + string(buf)
		line = strings.Replace(line, "SCRIPTFILES", scriptfiles, -1)
		line = strings.Replace(line, "SSH_PASSWORD", r.Form.Get("password"), -1)
	}
	line = line + "}"
	ioutil.WriteFile(reportlog[timest]["newjson"], []byte(line), 0)

	// new cfg file part
	var partitions string
        if index:=strings.LastIndex(r.Form.Get("ostype"),"CentOS");index>=0 {
	   for k, v := range r.Form["part"] {
		sizen, _ := strconv.Atoi(r.Form["size"][k])
		sizens := strconv.Itoa(sizen * 1024)
                if v=="swap" {
                  partitions = partitions + "part swap --size=" + sizens + "\n"
                }else{
		  partitions = partitions + "part " + v + " --fstype=ext4 --size=" + sizens + "\n"
                }
	   }
        }
        if index:=strings.LastIndex(r.Form.Get("ostype"),"Ubuntu");index>=0 {
	   for k, v := range r.Form["part"] {
		sizen, _ := strconv.Atoi(r.Form["size"][k])
		sizens := strconv.Itoa(sizen * 1024)
                if k==0 {
                partitions = partitions + "d-i partman-auto/method string regular\nd-i partman-auto/expert_recipe string boot-root :: "
                if v=="swap" {
                partitions = partitions + "64 "+sizens+" 300% $primary{ } linux-swap method{ swap } format{ } . "
                }else if v=="/boot" {
                partitions = partitions + "64 "+sizens+" 200 ext4 $primary{ } $bootable{ } method{ format } format{ } use_filesystem{ } filesystem{ ext4 } mountpoint{ /boot } "
                }else{
		partitions = partitions + sizens+" 4000 -1 ext4 $primary{ } method{ format } format{ } use_filesystem{ } filesystem{ ext4 } mountpoint{ "+v+" } . "
                }
                }else {
                if v=="swap" {
                partitions = partitions + "64 "+sizens+" 300% linux-swap method{ swap } format{ } . "
                }else if v=="/boot" {
                partitions = partitions + "64 "+sizens+" 200 ext4 $bootable{ } method{ format } format{ } use_filesystem{ } filesystem{ ext4 } mountpoint{ /boot } "
                }else{
		partitions = partitions + sizens+" 4000 -1 ext4 method{ format } format{ } use_filesystem{ } filesystem{ ext4 } mountpoint{ "+v+" } . "
                }
                }
	   }
        }
	os.Create(reportlog[timest]["newcfg"])
	buf, _ = ioutil.ReadFile(cfg)
	line = string(buf)
	line = strings.Replace(line, "SSH_USERNAME", r.Form.Get("user"), -1)
	line = strings.Replace(line, "SSH_PASSWORD", r.Form.Get("password"), -1)
	line = strings.Replace(line, "PARTITIONS", partitions, -1)
	ioutil.WriteFile(reportlog[timest]["newcfg"], []byte(line), 0)
        if index:=strings.LastIndex(r.Form.Get("ostype"),"Windows");index>=0 {
          copydir("template/floppy",timest)
        }

	fmt.Println(reportlog[timest]["newjson"])
	return timest
}
func callpacker(timest string) *os.Process {
	fmt.Println("callpacker", timest)
	inf, _ := os.Create(reportlog[timest]["resultdir"] + "inf.log")
	outf, _ := os.Create(reportlog[timest]["resultdir"] + "packer.log")
	errf, _ := os.Create(reportlog[timest]["resultdir"] + "errf.log")
	attr := &os.ProcAttr{
		//Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		Files: []*os.File{inf, outf, errf},
	}
	p, err := os.StartProcess("/home/packerdir/packer", []string{"/home/packerdir/packer", "build", reportlog[timest]["newjson"]}, attr)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("p=[", p, "]")
        reportlog[timest]["packerpid"]=strconv.Itoa(p.Pid)
	return p
}
func calltransform(p *os.Process, timest string){
        checkstatus(p,"packer",timest)
	if reportlog[timest]["compat"] == "0.1" {
		fmt.Println("calltransform")
		inf, _ := os.Create(reportlog[timest]["resultdir"] + "inf2.log")
		outf, _ := os.Create(reportlog[timest]["resultdir"] + "convert.log")
		errf, _ := os.Create(reportlog[timest]["resultdir"] + "errf2.log")
		attr := &os.ProcAttr{
			//Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
			Files: []*os.File{inf, outf, errf},
		}
		output:=reportlog[timest]["resultdir"]+"output/"+reportlog[timest]["vmname"]
		newoutput:=reportlog[timest]["resultdir"]+"output/tr"+reportlog[timest]["vmname"]
		fmt.Println("output=[", output, "]")
		fmt.Println("newoutput=[", newoutput, "]")
		p2, err := os.StartProcess("/bin/qemu-img", []string{"/bin/qemu-img", "convert", "-f", "qcow2", output, "-O", "qcow2", "-o", "compat=0.10", newoutput}, attr)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("p2=[", p2, "]")
		reportlog[timest]["transformpid"]=strconv.Itoa(p2.Pid)
		go checkstatus(p2,"transform",timest)		
	}
}
func checkstatus(p *os.Process,pname string,timest string){
        fmt.Println("checkstatus", pname,p)
        reportlog[timest]["status"]=pname+" running"
        liner, _ := json.Marshal(reportlog)
	ioutil.WriteFile("static/data/reportlog.json", liner, 0)
	pw, _ := p.Wait()
        fmt.Println("checkstatus over", p)
        fmt.Println("timest=", timest)
        if status:=pw.Success();status==true {
            reportlog[timest]["status"]=pname+" success"
            fmt.Println("checkstatus over success ",pname, p)
        }else{
	    reportlog[timest]["status"]=pname+" exit"
            fmt.Println("checkstatus over exit ",pname, p)
        }
        fmt.Println("checkstatus over SystemTime ",pname, p,pw.SystemTime())
        fmt.Println("checkstatus over UserTime ",pname, p,pw.UserTime())
        reportlog[timest][pname]=string(pw.SystemTime())
        liner, _ = json.Marshal(reportlog)
	ioutil.WriteFile("static/data/reportlog.json", liner, 0)

}

func copydir(path string,timest string) {
        err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
                if ( f == nil ) {return err}
                if f.IsDir() {
                  os.MkdirAll(strings.Replace(path, "template/", reportlog[timest]["resultdir"], -1), 0777)
                  return nil
                }
		newfil, _ := os.Create(strings.Replace(path, "template/", reportlog[timest]["resultdir"], -1))
		oldfil, _ := os.Open(path)
		io.Copy(newfil, oldfil)
                oldfil.Close()
                println(path)
                println(strings.Replace(path, "template/", reportlog[timest]["resultdir"], -1))
                return nil
        })
        if err != nil {
                fmt.Printf("filepath.Walk() returned %v\n", err)
        }
}
