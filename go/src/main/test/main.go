package main

import (
	//"encoding/json"
	"fmt"
	cloudstack "golang-cloudstack-library"
	"io/ioutil"
	"log"
	"net/url"
	"os"
        "reflect"
        "strconv"
        "strings"
	//	"bufio"
	//	"strconv"
	//	"flag"
	//	"time"
	//	"strings"
)

func main() {

	log.SetOutput(ioutil.Discard)

	endpoint, _ := url.Parse("http://192.168.10.2:8080/client/api")
	//apikey := "EPBWm44S0GpmA0oWomJu_xxIYpfsvg_oYtu1Z2WCvDnPUESOmWqPW04v3LVKbr8UQJUOQxc0nttsbFXP-LqWAw"
	//secretkey := "xNPQp_HU4cQtIlEJJwSbEjU4CAGa-TMqK-VWm5VblAa_mXDnea2Tj6DxPGlCJcUN2ajfqFOit_pL_qcmEgfS_g"
	apikey := "ji2f2JVLM2hGd6pydLg6vHI80h-xELg3yFGxfCGQiXQEXHhT56yGji5l_FEa40ITdVJvzekuX4yK9oBZYD8XDw"
	secretkey := "_7wEyHx7ifBleJNV-Y9Tmr4FlsRNNTRZlDvMpxj1aMEGXQfn8ywaJx_E8PtEwSfbVGgupG2veMTFccARv2_Kpg"

	username := "admin"
	password := "password"

	client, _ := cloudstack.NewClient(endpoint, apikey, secretkey, username, password)

	// for displaying zone information
	//params := cloudstack.NewListZonesParameter()
	//params.Name.Set("FirstZone")

	//zones, _ := client.ListZones(params)
	//b, _ := json.MarshalIndent(zones, "", "    ")

	//fmt.Println("Count:", len(zones))
	//fmt.Println(string(b))

	// For displaying all of the installed templates(templateFilter is all).
	params := cloudstack.NewListTemplatesParameter("all")

	templates, _ := client.ListTemplates(params)
	//b, _ := json.MarshalIndent(templates, "", "    ")

	fmt.Println("Count:", len(templates))
	//fmt.Println(string(b))
	fmt.Println(os.Args[0])

        mapa:=map[string]bool    {
        "account": false,
        "accountid": false,
        "bootable": false,
        "checksum": false,
        "created": false,
        "crosszones": false,
        "details": false,
        "displaytext": true,
        "domain": false,
        "domainid": false,
        "format": true,
        "hostid": false,
        "hostname": false,
        "hypervisor": true,
        "id": false,
        "isdynamicallyscalable": false,
        "isextractable": false,
        "isfeatured": false,
        "ispublic": false,
        "isready": false,
        "name": true,
        "ostypeid": false,
        "ostypename": false,
        "passwordenabled": true,
        "project": false,
        "projectid": false,
        "removed": false,
        "size": true,
        "sourcetemplateid": false,
        "sshkeyenabled": false,
        "status": false,
        "tags": false,
        "templatetag": false,
        "templatetype": false,
        "zoneid": false,
        "zonename": false,
    }
        var liner string
        num :=[]int{}
        for k, v := range templates{
          val := reflect.ValueOf(v).Elem()
          if k==0 {
              for i := 0; i < val.NumField(); i++ {
                     if mapa[strings.ToLower(val.Type().Field(i).Name)]{
                        num=append(num,i)
                     }
              }
              for _, i := range num {
                     fmt.Printf("|"+strconv.Itoa(i)+val.Type().Field(i).Name)
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
                 fmt.Printf("|"+strconv.Itoa(i)+val.Type().Field(i).Name)
                 liner=liner+"|"+val.Type().Field(i).Name
              }else if val.Type().Field(i).Name=="ID" {
                 fmt.Printf("|"+strconv.Itoa(i)+val.Field(i).Field(0).Field(0).MethodByName("String").Call([]reflect.Value{})[0].String())
                 liner=liner+"|"+val.Field(i).Field(0).Field(0).MethodByName("String").Call([]reflect.Value{})[0].String()
              }else{
                 fmt.Printf("|"+strconv.Itoa(i)+val.Field(i).Field(0).MethodByName("String").Call([]reflect.Value{})[0].String())
                 liner=liner+"|"+val.Field(i).Field(0).MethodByName("String").Call([]reflect.Value{})[0].String()
              }
           }
           fmt.Printf("|\n")
           liner=liner+"|\n"
       }
       ioutil.WriteFile("tlist.md", []byte(liner), 0)

	// For registering a new template.
	//*//params := cloudstack.NewRegisterTemplateParameter("gotinydisplaytext", "QCOW2", "KVM", "gotiny", "2077c4ce-6073-11e6-8c8b-5254004ab03d", "http://192.168.1.69/tiny.qcow2", "654d87f3-15a8-4cca-8018-2f451d853871")

	//*//templates, err := client.RegisterTemplate(params)
	//*//if (err == nil) {
	//*//b, _ := json.MarshalIndent(templates, "", "    ")
	//*//// fmt.Println("Count:", len(templates))
	//*//fmt.Println(string(b))
	//*//fmt.Println(os.Args[0])
	//*//} else {
	//*//        fmt.Println(err.Error())
	//*//}
}
