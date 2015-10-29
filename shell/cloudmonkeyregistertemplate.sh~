name="manager-base10"
displaytext="manager-base10"
featured="False"
passwordenabled="False"
public="True"
ostype="CentOS 6 (64-bit)"
format="QCOW2"
hypervisor="KVM"
url="http://192.168.122.1/template/manager-base7.qcow2"
zonename="Zone1"
cloudmonkey set display table
zoneid=`cloudmonkey list zones filter=id,name name="$zonename"|grep "$zonename"|awk '{print $4}'`
ostypeid=`cloudmonkey list ostypes description="CentOS 6 (64-bit)"|grep "CentOS 6 (64-bit)"|awk -F"[|]" '{print $5}'|sed -e "s/ //g"`
cloudmonkey set display default
cloudmonkey register template name="${name}" displaytext="${displaytext}" isfeatured=${featured} passwordenabled=${passwordenabled} ispublic=${public} ostypeid=${ostypeid} format=${format} hypervisor=${hypervisor} zoneid=${zoneid} url=${url}

