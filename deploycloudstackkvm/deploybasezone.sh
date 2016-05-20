#!/bin/bash

cli=cloudmonkey
dns_ext=8.8.8.8
dns_int=10.1.6.30
gw=10.1.6.1
nmask=255.255.255.0
hpvr=KVM
pod_start=10.1.6.225
pod_end=10.1.6.234
vlan_start=10.1.6.235
vlan_end=10.1.6.254

#Put space separated host ips in following
host_ips=10.1.6.20
host_user=root
host_passwd=engine
sec_storage=nfs://10.1.6.30/exports/secondary
prm_storage=nfs://10.1.6.30/exports/primary

zone_id=`$cli create zone dns1=$dns_ext internaldns1=$dns_int name=MyZone networktype=Basic | grep ^id\ = | awk '{print $3}'`
echo "Created zone" $zone_id

phy_id=`$cli create physicalnetwork name=phy-network zoneid=$zone_id | grep ^id\ = | awk '{print $3}'`
echo "Created physical network" $phy_id
$cli add traffictype traffictype=Guest physicalnetworkid=$phy_id
echo "Added guest traffic"
$cli add traffictype traffictype=Management physicalnetworkid=$phy_id
echo "Added mgmt traffic"
$cli update physicalnetwork state=Enabled id=$phy_id
echo "Enabled physicalnetwork"

nsp_id=`$cli list networkserviceproviders name=VirtualRouter physicalnetworkid=$phy_id | grep ^id\ = | awk '{print $3}'`
vre_id=`$cli list virtualrouterelements nspid=$nsp_id | grep ^id\ = | awk '{print $3}'`
$cli api configureVirtualRouterElement enabled=true id=$vre_id
$cli update networkserviceprovider state=Enabled id=$nsp_id
echo "Enabled virtual router element and network service provider"

nsp_sg_id=`$cli list networkserviceproviders name=SecurityGroupProvider physicalnetworkid=$phy_id | grep ^id\ = | awk '{print $3}'`
$cli update networkserviceprovider state=Enabled id=$nsp_sg_id
echo "Enabled security group provider"

netoff_id=`$cli list networkofferings name=DefaultSharedNetworkOfferingWithSGService | grep ^id\ = | awk '{print $3}'`
net_id=`$cli create network zoneid=$zone_id name=guestNetworkForBasicZone displaytext=guestNetworkForBasicZone networkofferingid=$netoff_id | grep ^id\ = | awk '{print $3}'`
echo "Created network $net_id for zone" $zone_id

pod_id=`$cli create pod name=MyPod zoneid=$zone_id gateway=$gw netmask=$nmask startip=$pod_start endip=$pod_end | grep ^id\ = | awk '{print $3}'`
echo "Created pod"

$cli create vlaniprange podid=$pod_id networkid=$net_id gateway=$gw netmask=$nmask startip=$vlan_start endip=$vlan_end forvirtualnetwork=false
echo "Created IP ranges for instances"

cluster_id=`$cli add cluster zoneid=$zone_id hypervisor=$hpvr clustertype=CloudManaged podid=$pod_id clustername=MyCluster | grep ^id\ = | awk '{print $3}'`
echo "Created cluster" $cluster_id

#Put loop here if more than one
for host_ip in $host_ips;
do
  $cli add host zoneid=$zone_id podid=$pod_id clusterid=$cluster_id hypervisor=$hpvr username=$host_user password=$host_passwd url=http://$host_ip;
  echo "Added host" $host_ip;
done;

$cli create storagepool zoneid=$zone_id podid=$pod_id clusterid=$cluster_id name=MyNFSPrimary url=$prm_storage
echo "Added primary storage"

$cli add secondarystorage zoneid=$zone_id url=$sec_storage
echo "Added secondary storage"

$cli update zone allocationstate=Enabled id=$zone_id
echo "Basic zone deloyment completed!"
