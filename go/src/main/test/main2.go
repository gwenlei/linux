package main

import (
	//"encoding/json"
	"fmt"
	cloudstack "golang-cloudstack-library"
	"io/ioutil"
	"log"
	"net/url"
	"os"
        "encoding/json"
	//"reflect"
	//"strconv"
	//"strings"
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

	// registering a new template.
	params := cloudstack.NewRegisterTemplateParameter("gotinydisplaytext2", "QCOW2", "KVM", "gotiny2", "cb48aecd-6db9-11e6-9e4c-5254005357ff", "http%3A%2F%2F192.168.1.69%2Ftiny.qcow2", "d44975e0-4c2e-4c93-8626-8e1d4a3cc9b5")

	templates, err := client.RegisterTemplate(params)
	if err == nil {
		b, _ := json.MarshalIndent(templates, "", "    ")
		// fmt.Println("Count:", len(templates))
		fmt.Println(string(b))
		fmt.Println(os.Args[0])
	} else {
		fmt.Println(err.Error())
	}
}
