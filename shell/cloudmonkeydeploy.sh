vmname="testuserdata14"
templatename="manager-base7"
serviceofferingname="Small Instance"
zonename="Zone1"
networkname="iso-vn-001"

userdata=`cat my-user-data|base64`

cloudmonkey set display table
templateid=`cloudmonkey list templates templatefilter=executable filter=name,id name="$templatename"|grep "$templatename"|awk -F"[|]" '{print $2}'|sed -e "s/ //g"`
serviceofferingid=`cloudmonkey list serviceofferings filter=name,id name="$serviceofferingname"|grep "$serviceofferingname"|awk -F"[|]" '{print $3}'|sed -e "s/ //g"`
zoneid=`cloudmonkey list zones filter=id,name name="$zonename"|grep "$zonename"|awk -F"[|]" '{print $3}'|sed -e "s/ //g"`
networkid=`cloudmonkey list networks filter=name,id name="$networkname"|grep "$networkname"|awk -F"[|]" '{print $2}'|sed -e "s/ //g"`

echo templatename=$templatename templateid=$templateid
echo serviceofferingname=$serviceofferingname serviceofferingid=$serviceofferingid
echo zonename=$zonename zoneid=$zoneid
echo networkname=$networkname networkid=$networkid
echo userdata=$userdata

id=$(cloudmonkey deploy virtualmachine startvm=false serviceofferingid=$serviceofferingid templateid=$templateid zoneid=$zoneid networkids=$networkid name=$vmname userdata="`cat my-user-data|base64`"|awk '/^id =/ {print $3}')

echo $id
ip=`grep IPADDR my-user-data|awk  -F"[=]" '{print $2}'`
mysql -u root <<EOF
use cloud;
update vm_instance set private_ip_address='$ip' where uuid='$id';
update user_vm_view set ip_address='$ip' where  uuid='$id';
EOF


