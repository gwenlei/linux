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
)

var dat map[string](map[string]string)
var resultdir string

func main() {
	buf, _ := ioutil.ReadFile("static/data/data.json")
	if err := json.Unmarshal(buf, &dat); err != nil {
		panic(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", index)
	http.HandleFunc("/build", build) //设置访问的路由
	http.HandleFunc("/setdat", setdat)
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
	resultdir = dat["resultmap"]["resultdir"] + time.Now().Format("20060102150405") + "/"
	os.MkdirAll(resultdir, 0777)
	os.MkdirAll(resultdir+"json/", 0777)
	os.MkdirAll(resultdir+"script/", 0777)
	os.MkdirAll(resultdir+"cfg/", 0777)
        tmplog:=resultdir+"form.log"
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
	newjson := resultdir+"json/" + json[strings.LastIndex(json, "/")+1:]
	cfg := dat["cfgmap"][r.Form.Get("ostype")]
	newcfg := resultdir+"cfg/" + cfg[strings.LastIndex(cfg, "/")+1:]
	newcfgs := "http://"+dat["servermap"]["server"]+"/" + newcfg
	iso := dat["isomap"][r.Form.Get("ostype")]
	disksizen, _ := strconv.Atoi(r.Form.Get("disksize"))
	disksizens := strconv.Itoa(disksizen * 1024)
	//new json file
	os.Create(newjson)
	buf, _ := ioutil.ReadFile(json)
	line := string(buf)
	line = strings.Replace(line, "DISK_SIZE", disksizens, -1)
	line = strings.Replace(line, "SSH_USERNAME", r.Form.Get("user"), -1)
	line = strings.Replace(line, "SSH_PASSWORD", r.Form.Get("password"), -1)
	line = strings.Replace(line, "VM_NAME", r.Form.Get("vmname"), -1)
	line = strings.Replace(line, "OUTPUT_DIRECTORY", resultdir+"output/", -1)
	line = strings.Replace(line, "ISO_CHECKSUM", dat["md5map"][dat["isomap"][r.Form.Get("ostype")]], -1)
	line = strings.Replace(line, "ISO_URL", iso, -1)
	line = strings.Replace(line, "KS_CFG", newcfgs, -1)
        line = strings.Replace(line, "WIN_CFG", newcfg, -1)
        line = strings.Replace(line, "RESULTDIR", resultdir, -1)
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
			newscript[k] = resultdir+"script/" + v[strings.LastIndex(v, "/")+1:]
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
	ioutil.WriteFile(newjson, []byte(line), 0)

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
	os.Create(newcfg)
	buf, _ = ioutil.ReadFile(cfg)
	line = string(buf)
	line = strings.Replace(line, "SSH_USERNAME", r.Form.Get("user"), -1)
	line = strings.Replace(line, "SSH_PASSWORD", r.Form.Get("password"), -1)
	line = strings.Replace(line, "PARTITIONS", partitions, -1)
	ioutil.WriteFile(newcfg, []byte(line), 0)
        if index:=strings.LastIndex(r.Form.Get("ostype"),"Windows");index>=0 {
          copydir("template/floppy")
        }

	fmt.Println(json)
	fmt.Println(cfg)
	fmt.Println(script)
	fmt.Println(iso)
	fmt.Println(newjson)
	return newjson
}
func callpacker(json string) *os.Process {
	fmt.Println("callpacker", json)
	inf, _ := os.Create(resultdir + "inf.log")
	outf, _ := os.Create(resultdir + "packer.log")
	errf, _ := os.Create(resultdir + "errf.log")
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
	inf, _ := os.Create(resultdir + "inf2.log")
	outf, _ := os.Create(resultdir + "convert.log")
	errf, _ := os.Create(resultdir + "errf2.log")
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

func copydir(path string) {
        err := filepath.Walk(path, func(path string, f os.FileInfo, err error) error {
                if ( f == nil ) {return err}
                if f.IsDir() {
                  os.MkdirAll(strings.Replace(path, "template/", resultdir, -1), 0777)
                  return nil
                }
		newfil, _ := os.Create(strings.Replace(path, "template/", resultdir, -1))
		oldfil, _ := os.Open(path)
		io.Copy(newfil, oldfil)
                oldfil.Close()
                println(path)
                println(strings.Replace(path, "template/", resultdir, -1))
                return nil
        })
        if err != nil {
                fmt.Printf("filepath.Walk() returned %v\n", err)
        }
}
