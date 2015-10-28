vmname="testuserdata14"
templatename="manager-base7"
serviceofferingname="Small Instance"
zonename="Zone1"
networkname="iso-vn-001"

userdata=`cat my-user-data|base64`

cloudmonkey set display table
templateid=`cloudmonkey list templates templatefilter=executable filter=name,id name="$templatename"|grep "$templatename"|awk '{print $2}'`
serviceofferingid=`cloudmonkey list serviceofferings filter=name,id name="$serviceofferingname"|grep "$serviceofferingname"|awk '{print $5}'`
zoneid=`cloudmonkey list zones filter=id,name name="$zonename"|grep "$zonename"|awk '{print $4}'`
networkid=`cloudmonkey list networks filter=name,id name="$networkname"|grep "$networkname"|awk '{print $2}'`

echo templatename=$templatename templateid=$templateid
echo serviceofferingname=$serviceofferingname serviceofferingid=$serviceofferingid
echo zonename=$zonename zoneid=$zoneid
echo networkname=$networkname networkid=$networkid
echo userdata=$userdata

cloudmonkey deploy virtualmachine startvm=false serviceofferingid=$serviceofferingid templateid=$templateid zoneid=$zoneid networkids=$networkid name=$vmname userdata="`cat my-user-data|base64`"


