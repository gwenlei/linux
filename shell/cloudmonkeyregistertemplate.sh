name="manager-base9"
displaytext="manager-base9"
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
ostypeid=`cloudmonkey list ostypes description="$ostype"|grep "$ostype"|awk '{print $10}'`
cloudmonkey set display default
cloudmonkey register template name="${name}" displaytext="${displaytext}" isfeatured=${featured} passwordenabled=${passwordenabled} ispublic=${public} ostypeid=${ostypeid} format=${format} hypervisor=${hypervisor} zoneid=${zoneid} url=${url}

