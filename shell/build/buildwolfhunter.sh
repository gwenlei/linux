echo `date '+%Y%m%d-%H%M%S'`
curday=`date '+%Y%m%d'`
#curday=20160216

if [ ! -d "/home/clouder/wolfhunter_repo/$curday" ] ; then
echo noupdate `date '+%Y%m%d-%H%M%S'`
exit -1
fi

echo buildwolfhunter begin at `date '+%Y%m%d-%H%M%S'`
cd /home/clouder/leisj/build
wolfname=leiwolfhunter$curday
wolfsourcefile=/home/clouder/leisj/build/temp/temp$curday.qcow2
wolfoutfile=/home/clouder/leisj/build/temp/wolfhunter-$curday.qcow2
backingfile=/home/clouder/wolfhunter_repo/base/wolfhunter-0830.qcow2
wolfxml=/home/clouder/leisj/build/xml/vm.xml
wolfxmltemp=/home/clouder/leisj/build/temp/temp$curday.xml

rm -f $wolfsourcefile
qemu-img create -f qcow2 $wolfsourcefile -o backing_file=$backingfile 100G
sed "s#VMNAME#$wolfname#g" $wolfxml > $wolfxmltemp
sed -i "s#SOURCEFILE#$wolfsourcefile#g" $wolfxmltemp

virsh shutdown $wolfname
echo `date '+%Y%m%d-%H%M%S'`
sleep 10
virsh destroy $wolfname
virsh undefine $wolfname
virsh define $wolfxmltemp
virsh start $wolfname

echo wait vm boot for 1m 
echo `date '+%Y%m%d-%H%M%S'`
sleep 1m
sed -i "/192.168.173.2/"d /root/.ssh/known_hosts
ssh-keyscan 192.168.173.2 >> /root/.ssh/known_hosts
echo sshcopyid for root@192.168.173.2  
#ssh-copy-id -i ~/.ssh/id_rsa.pub root@192.168.173.2
./sshlogin.exp root@192.168.173.2 engine

####update ceph in wolfhunter.qcow2
if [ ! -d "/home/clouder/wolfhunter_repo/$curday/ceph" ] ; then
echo no ceph update `date '+%Y%m%d-%H%M%S'`
else
echo update ceph `date '+%Y%m%d-%H%M%S'`
scp /home/clouder/wolfhunter_repo/$curday/ceph/* root@192.168.173.2:/var/www/html/repo/centos71-ceph-cloudstack-repo/packages/
ssh root@192.168.173.2 <<EOF 
cd /var/www/html/repo/centos71-ceph-cloudstack-repo
createrepo .
rm -f /etc/udev/rules.d/70-persistent-net.rules
rm -f /etc/udev/rules.d/75-persistent-net-generator.rules
echo "# " > /etc/udev/rules.d/75-persistent-net-generator.rules
EOF
fi

####update cloudstack-agent in wolfhunter.qcow2
if [ ! -d "/home/clouder/wolfhunter_repo/$curday/cloudstack/agent" ] ; then
echo no agent update `date '+%Y%m%d-%H%M%S'`
else
echo update agent `date '+%Y%m%d-%H%M%S'`
ssh root@192.168.173.2 <<EOF 
cd /var/www/html/repo/centos71-ceph-cloudstack-repo/packages/
rm -f cloudstack*
EOF
scp /home/clouder/wolfhunter_repo/$curday/cloudstack/agent/* root@192.168.173.2:/var/www/html/repo/centos71-ceph-cloudstack-repo/packages/
ssh root@192.168.173.2 <<EOF 
cd /var/www/html/repo/centos71-ceph-cloudstack-repo
createrepo .
rm -f /etc/udev/rules.d/70-persistent-net.rules
rm -f /etc/udev/rules.d/75-persistent-net-generator.rules
echo "# " > /etc/udev/rules.d/75-persistent-net-generator.rules
EOF
fi

####update cloudstack-management in wolfhunter.qcow2
if [[ ! -d "/home/clouder/wolfhunter_repo/$curday/cloudstack/management" && ! -d "/home/clouder/wolfhunter_repo/$curday/auto-script" ]] ; then
echo no manager update `date '+%Y%m%d-%H%M%S'`
else

managername=leimanager$curday
managersourcefile=/home/clouder/leisj/build/temp/manager-base$curday.qcow2
managerxml=/home/clouder/leisj/build/xml/manager.xml
managerxmltemp=/home/clouder/leisj/build/temp/manager$curday.xml
sed "s#VMNAME#$managername#g" $managerxml > $managerxmltemp
sed -i "s#SOURCEFILE#$managersourcefile#g" $managerxmltemp
virsh destroy $managername
virsh undefine $managername
scp root@192.168.173.2:/var/www/html/repo/centos71-ceph-cloudstack-repo/template/manager-base.qcow2 $managersourcefile
virsh undefine
virsh define $managerxmltemp
virsh start $managername
echo wait manager boot for 1m
echo `date '+%Y%m%d-%H%M%S'`
sleep 1m
sed -i "/192.168.173.10/"d /root/.ssh/known_hosts
ssh-keyscan 192.168.173.10 >> /root/.ssh/known_hosts
echo sshcopyid to root@192.168.173.10 
#ssh-copy-id -i ~/.ssh/id_rsa.pub root@192.168.173.10
./sshlogin.exp root@192.168.173.10 engine

if [ ! -d "/home/clouder/wolfhunter_repo/$curday/cloudstack/management" ] ; then
echo no mamangement update `date '+%Y%m%d-%H%M%S'`
else
echo update management
scp /home/clouder/wolfhunter_repo/$curday/cloudstack/management/* root@192.168.173.10:/root

ssh root@192.168.173.10 <<EOF
cd /root
yum remove -y cloudstack-common-4.5.1-shapeblue0.el6.x86_64 cloudstack-management-4.5.1-shapeblue0.el6.x86_64 cloudstack-awsapi-4.5.1-shapeblue0.el6.x86_64
rpm -ivh cloudstack-management-4.5.1-1.el6.x86_64.rpm cloudstack-awsapi-4.5.1-1.el6.x86_64.rpm cloudstack-common-4.5.1-1.el6.x86_64.rpm
cloudstack-setup-databases cloud:engine@localhost --deploy-as=root:engine
cloudstack-setup-management
service cloudstack-management restart
echo wait cloudstack-management restart for 5m
echo `date '+%Y%m%d-%H%M%S'`
sleep 5m
rm -f /etc/udev/rules.d/70-persistent-net.rules
rm -f /etc/udev/rules.d/75-persistent-net-generator.rules
echo "# " > /etc/udev/rules.d/75-persistent-net-generator.rules
EOF
fi

####update auto-script in wolfhunter.qcow2
if [ ! -d "/home/clouder/wolfhunter_repo/$curday/auto-script" ] ; then
echo no auto-script update `date '+%Y%m%d-%H%M%S'`
else
echo update auto-script
scp /home/clouder/wolfhunter_repo/$curday/auto-script/* root@192.168.173.10:/home/deploy
fi

echo `date '+%Y%m%d-%H%M%S'`
virsh shutdown $managername
sleep 20
virsh destroy $managername
scp $managersourcefile root@192.168.173.2:/var/www/html/repo/centos71-ceph-cloudstack-repo/template/manager-base.qcow2 

fi

virsh shutdown $wolfname
sleep 20
virsh destroy $wolfname
echo qemu-img convert ing at `date '+%Y%m%d-%H%M%S'`
qemu-img convert -O qcow2 $wolfsourcefile $wolfoutfile
#cp $wolfoutfile /home/clouder/wolfhunter_repo/$curday/
echo finish at `date '+%Y%m%d-%H%M%S'`
exit 0
