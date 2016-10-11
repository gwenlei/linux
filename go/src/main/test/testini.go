package main
import ("fmt";goini "github.com/zieckey/goini";"os";cloudstack "golang-cloudstack-library";"strings"; "net/url";"reflect";"encoding/json";"io/ioutil";"log")

func main() {
log.SetOutput(ioutil.Discard)
if len(os.Args) != 3 {
    fmt.Println("usage: ./testini file.ini register|delete|list")
    return
}
ini := goini.New()
err := ini.ParseFile(os.Args[1])
if err != nil {
    fmt.Printf("parse INI file %v failed : %v\n", os.Args[1], err.Error())
    return
}

client:=callclient(ini)
if os.Args[2]=="register"{
   callregister(ini,client)
}else if os.Args[2]=="delete"{
   calldelete(ini,client)
}else if os.Args[2]=="list"{
   calllist(ini,client)
}
}

func callclient(ini *goini.INI) (client *cloudstack.Client) {
        endpoints, _ := ini.SectionGet("main", "endpoint")
        apikey, _ := ini.SectionGet("main", "apikey")
        secretkey, _ := ini.SectionGet("main", "secretkey")
        username, _ := ini.SectionGet("main", "username")
        password, _ := ini.SectionGet("main", "password")
	endpoint, _ := url.Parse(endpoints)
	client, _ = cloudstack.NewClient(endpoint, apikey, secretkey, username, password)
        return client
}

func callregister(ini *goini.INI , client *cloudstack.Client) {
        zonename, _ := ini.SectionGet("main", "zonename")
        ostype, _ := ini.SectionGet("main", "ostype")
        format, _ := ini.SectionGet("main", "format")
        hypervisor, _ := ini.SectionGet("main", "hypervisor")
        passwordenabled, _ := ini.SectionGetBool("main", "passwordenabled")
        registermap, _ := ini.GetKvmap("register")
	params := cloudstack.NewListOstypesParameter(strings.Replace(ostype, "%", "%25", -1))
	ostypes, _ := client.ListOstypes(params)
        var ostypeid string
        if len(ostypes)>0 {
           ostypeid=ostypes[0].Id.String()
        }else{
           fmt.Println("ostype is not exist")
           os.Exit(1)
        }

	params1 := cloudstack.NewListZonesParameter(strings.Replace(zonename, "%", "%25", -1))
	zones, _ := client.ListZones(params1)
        var zoneid string
        if len(zones)>0 {
           zoneid=zones[0].Id.String()
        }else{
           fmt.Println("zonename is not exist")
           os.Exit(1)
        }

	// registering a new template.
        for k, v := range registermap {
           url:= strings.Replace(v, ":", "%3A", -1)
           url = strings.Replace(url, "/", "%2F", -1)
           fmt.Println(k,v)
	   params2 := cloudstack.NewRegisterTemplateParameter(url, format, hypervisor, url, ostypeid, url, zoneid)
           params2.IsPublic.Set(true)
           params2.PasswordEnabled.Set(passwordenabled)
           templates, err := client.RegisterTemplate(params2)
	   if err == nil {
                fmt.Println("return template id : ",templates[0].Id.String())
	   } else {
		fmt.Println(err.Error())
	   }
       }
}

func calldelete(ini *goini.INI , client *cloudstack.Client) {
        deletemap, _ := ini.GetKvmap("delete")
        for _, v := range deletemap {
	  params := cloudstack.NewDeleteTemplateParameter(v)
  	  templates, err := client.DeleteTemplate(params)
	  if err == nil {
                fmt.Println(v)
		b, _ := json.MarshalIndent(templates, "", "    ")
		fmt.Println(string(b))
	  } else {
		fmt.Println(err.Error())
	  }
        }
}

func calllist(ini *goini.INI , client *cloudstack.Client){
        keyword, _ := ini.SectionGet("list", "keyword")
	params := cloudstack.NewListTemplatesParameter("all")
        if keyword != "all" {
          params.Keyword.Set(strings.Replace(keyword, "%", "%25", -1))
        }
	templates, _ := client.ListTemplates(params)
	fmt.Println("total:", len(templates))
        mapa, _ := ini.GetKvmap("list")
        var liner string
        num :=[]int{}
        for k, v := range templates{
          val := reflect.ValueOf(v).Elem()
          if k==0 {
              for i := 0; i < val.NumField(); i++ {
                     if mapa[strings.ToLower(val.Type().Field(i).Name)]=="true"{
                        num=append(num,i)
                     }
              }
              for _, i := range num {
                     //fmt.Printf("|"+strconv.Itoa(i)+val.Type().Field(i).Name)
                     fmt.Printf("|"+val.Type().Field(i).Name)
                     liner=liner+"|"+val.Type().Field(i).Name
              }
              fmt.Printf("|\n")
              liner=liner+"|\n"
              for _, _ = range num {
                     fmt.Printf("|----")
                     liner=liner+"|----"
              }
              fmt.Printf("|\n")
              liner=liner+"|\n"
           }
           //fmt.Printf("|"+val.Field(1).Field(0).MethodByName("String").Call([]reflect.Value{})[0].String()+"|\n")
           for _, i := range num {
              if val.Type().Field(i).Name=="ResourceBase" || val.Type().Field(i).Name=="Tags"{
                 //fmt.Printf("|"+strconv.Itoa(i)+val.Type().Field(i).Name)
                 fmt.Printf("|"+val.Type().Field(i).Name)
                 liner=liner+"|"+val.Type().Field(i).Name
              }else if val.Type().Field(i).Name=="ID" {
                 //fmt.Printf("|"+strconv.Itoa(i)+val.Field(i).Field(0).Field(0).MethodByName("String").Call([]reflect.Value{})[0].String())
                 fmt.Printf("|"+val.Field(i).Field(0).Field(0).MethodByName("String").Call([]reflect.Value{})[0].String())
                 liner=liner+"|"+val.Field(i).Field(0).Field(0).MethodByName("String").Call([]reflect.Value{})[0].String()
              }else{
                 //fmt.Printf("|"+strconv.Itoa(i)+val.Field(i).Field(0).MethodByName("String").Call([]reflect.Value{})[0].String())
                 fmt.Printf("|"+val.Field(i).Field(0).MethodByName("String").Call([]reflect.Value{})[0].String())
                 liner=liner+"|"+val.Field(i).Field(0).MethodByName("String").Call([]reflect.Value{})[0].String()
              }
           }
           fmt.Printf("|\n")
           liner=liner+"|\n"
       }
       ioutil.WriteFile("tlist.md", []byte(liner), 0)
}
