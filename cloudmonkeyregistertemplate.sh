name="manager-base8"
displaytext="manager-base8"
featured="False"
passwordenabled="False"
public="True"
ostypeid="cdde0ce4-56ca-11e5-b257-525400cd7603"
format="QCOW2"
hypervisor="KVM"
url="http://192.168.122.1/template/manager-base7.qcow2"
zonename="Zone1"
zoneid=`cloudmonkey list zones filter=id,name name="$zonename"|grep "$zonename"|awk '{print $4}'`
cloudmonkey set display default
cloudmonkey register template name="${name}" displaytext="${displaytext}" isfeatured=${featured} passwordenabled=${passwordenabled} ispublic=${public} ostypeid=${ostypeid} format=${format} hypervisor=${hypervisor} zoneid=${zoneid} url=${url} 
