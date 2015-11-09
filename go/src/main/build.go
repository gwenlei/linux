package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
    "strings"
    "os"
    "bufio"
    "io"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()         //解析url传递的参数，对于POST则解析响应包的主体（request body）
    //注意:如果没有调用ParseForm方法，下面无法获取表单的数据
    fmt.Println(r.Form)   //这些信息是输出到服务器端的打印信息
    fmt.Println("path", r.URL.Path)
    fmt.Println("scheme", r.URL.Scheme)
    fmt.Println(r.Form["url_long"])
    for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, ""))
    }
    fmt.Fprintf(w, "Hello astaxie!") //这个写入到w的是输出到客户端的
}

func login(w http.ResponseWriter, r *http.Request) {
    fmt.Println("method:", r.Method) //获取请求的方法
    if r.Method == "GET" {
        t, _ := template.ParseFiles("login.gtpl")
        t.Execute(w, nil)
    } else {
        //请求的是登陆数据，那么执行登陆的逻辑判断
        r.ParseForm()
        fmt.Println("username:", r.Form["username"])
        fmt.Println("password:", r.Form["password"])
    }
}

func main() {
    http.HandleFunc("/", sayhelloName)       //设置访问的路由
    http.HandleFunc("/login", login)         //设置访问的路由
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
        for _,v:=range r.Form["software"]{
          fmt.Println("v:=",v)
        }
        for k, v := range r.Form {
        fmt.Println("key:", k)
        fmt.Println("val:", strings.Join(v, " "))
        }
        fmt.Println("cloudstackip:", r.Form.Get("cloudstackip"))
        json:= buildjson(r)
        callpacker(json)
    }
}

func buildjson(r *http.Request)(json string){
     jsondir:="result/test2"
     os.MkdirAll(jsondir, 0777)
     json="result/test2/test2.json"
     f,_:=os.Create(json)
     f.WriteString("just a test")
     template:="template/json/centos6-6.json"
     ft,_:=os.Open(template)
     buf:=bufio.NewReader(ft)
     for {
      line,err := buf.ReadString('\n')
      if err== io.EOF{
        break
      }
      f.WriteString(line)
     }

     fmt.Println(json)
     defer ft.Close()
     defer f.Close()
     return json
}

func callpacker(json string){
     fmt.Println("callpacker",json)
}
