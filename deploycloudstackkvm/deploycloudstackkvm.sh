virsh shutdown testcentos71cskvm0
virsh destroy testcentos71cskvm0
virsh shutdown testcentos71cskvm
virsh destroy testcentos71cskvm
rm -f /home/html/downloads/CentOS7-1cskvm.qcow2
#wget -P /home/html/downloads/ --no-check-certificate https://192.168.0.82:9090/static/result/20160421103801/output/CentOS7-1.qcow2 
cp /home/code/mycode/go/src/main/static/result/20160421103801/output/CentOS7-1.qcow2 /home/html/downloads/CentOS7-1cskvm.qcow2
virsh undefine testcentos71cskvm0
virsh define xml/testcentos71.xml
virsh start testcentos71cskvm0
sleep 1m
testmac=`virsh domiflist testcentos71cskvm0|grep network|awk '{print $5}'`
testip=`virsh net-dhcp-leases default|grep $testmac|awk '{print $5}'|cut -d '/' -f 1`
sed -i /$testip/d /root/.ssh/known_hosts
ssh-keyscan $testip >> /root/.ssh/known_hosts
./sshcopyid.exp root@$testip engine
scp xml/ifcfg-eth0 root@$testip:/etc/sysconfig/network-scripts/ifcfg-eth0
ssh root@$testip "sync"
virsh shutdown testcentos71cskvm0
virsh destroy testcentos71cskvm0
virsh undefine testcentos71cskvm0
sed -i /$testip/d /root/.ssh/known_hosts

virsh shutdown testcentos71cskvm
virsh destroy testcentos71cskvm
virsh undefine testcentos71cskvm
virsh define xml/centos71cskvm.xml
virsh start testcentos71cskvm
sleep 1m
sed -i /10.1.6.30/d /root/.ssh/known_hosts
ssh-keyscan 10.1.6.30 >> /root/.ssh/known_hosts
./sshcopyid.exp root@10.1.6.30 engine
scp /etc/yum.repos.d/mrepo7.repo root@10.1.6.30:/etc/yum.repos.d/mrepo7.repo
scp /etc/yum.repos.d/cloudstack451_7.repo root@10.1.6.30:/etc/yum.repos.d/cloudstack451_7.repo
scp /home/code/mycode/deploycloudstackkvm/sshcopyid.exp root@10.1.6.30:/root
scp /home/code/mycode/deploycloudstackkvm/mysqlsecure.sh root@10.1.6.30:/root
scp /home/code/mycode/deploycloudstackkvm/mysql.sh root@10.1.6.30:/root
scp /home/code/mycode/deploycloudstackkvm/iptables.sh root@10.1.6.30:/root
scp /home/img/systemvmtest/systemvm64template-unknown2-kvm.qcow2.bz2 root@10.1.6.30:/root
ssh root@10.1.6.30 << EOF
chmod +x mysql.sh mysqlsecure.sh iptables.sh sshcopyid.exp
cd /etc/yum.repos.d/
rename .repo .repoback CentOS*
setenforce permissive
sed -i "/SELINUX=enforcing/ c\SELINUX=permissive" /etc/selinux/config
echo "127.0.0.1 localhost localhost.cstack.local" >> /etc/hosts
echo "10.1.6.30 csmankvm.cstack.local csmankvm" >> /etc/hosts
echo "10.1.6.20 csagentkvm.cstack.local csagentkvm" >> /etc/hosts
sed -i "/#UseDNS yes/ c\UseDNS no" /etc/ssh/sshd_config
service sshd restart
EOF
virsh reboot testcentos71cskvm
sleep 1m

ssh root@10.1.6.30 << EOF
yum install -y ntp
chkconfig ntpd on
service ntpd start
yum install -y cloudstack-management
mkdir /exports
mkdir -p /exports/primary
mkdir -p /exports/secondary
chmod 777 -R /exports
echo "/exports/primary *(rw,async,no_root_squash)" > /etc/exports
echo "/exports/secondary *(rw,async,no_root_squash)" >> /etc/exports
exportfs -a
#sed -i -e '/#MOUNTD_NFS_V3="no"/ c\MOUNTD_NFS_V3="yes"' -e '/#RQUOTAD_PORT=875/ c\RQUOTAD_PORT=875' -e '/#LOCKD_TCPPORT=32803/ c\LOCKD_TCPPORT=32803' -e '/#LOCKD_UDPPORT=32769/ c\LOCKD_UDPPORT=32769' -e '/#MOUNTD_PORT=892/ c\MOUNTD_PORT=892' -e '/#STATD_PORT=662/ c\STATD_PORT=662' -e '/#STATD_OUTGOING_PORT=2020/ c\STATD_OUTGOING_PORT=2020' /etc/sysconfig/nfs
echo MOUNTD_NFS_V3="yes" >> /etc/sysconfig/nfs
echo RQUOTAD_PORT=875 >> /etc/sysconfig/nfs
echo LOCKD_TCPPORT=32803 >> /etc/sysconfig/nfs
echo LOCKD_UDPPORT=32769 >> /etc/sysconfig/nfs
echo MOUNTD_PORT=892 >> /etc/sysconfig/nfs
echo STATD_PORT=662 >> /etc/sysconfig/nfs
echo STATD_OUTGOING_PORT=2020 >> /etc/sysconfig/nfs
systemctl start nfs-server.service
systemctl enable nfs-server.service

yum -y install mariadb*  
yum -y install vim
sed -i -e '/datadir/ a\innodb_rollback_on_timeout=1' -e '/datadir/ a\innodb_lock_wait_timeout=600' -e '/datadir/ a\max_connections=350' -e '/datadir/ a\log-bin=mysql-bin' -e "/datadir/ a\binlog-format = 'ROW'" -e "/datadir/ a\bind-address = 0.0.0.0" /etc/my.cnf
systemctl start mariadb.service
systemctl enable mariadb.service
/root/mysqlsecure.sh
cloudstack-setup-databases cloud:engine@127.0.0.1 --deploy-as=root:engine
cloudstack-setup-management
yum install nginx -y
service nginx start
#cd /usr/share/cloudstack-common/scripts/vm/hypervisor/xenserver/
#wget http://192.168.0.82/vhd-util
#chmod 755 /usr/share/cloudstack-common/scripts/vm/hypervisor/xenserver/vhd-util
#/usr/share/cloudstack-common/scripts/storage/secondary/cloud-install-sys-tmplt -m /exports/secondary -u http://192.168.0.82/systemvmtest/systemvm64template-unknown2-kvm.qcow2.bz2 -h kvm -F
/usr/share/cloudstack-common/scripts/storage/secondary/cloud-install-sys-tmplt -m /exports/secondary -f /root/systemvm64template-unknown2-kvm.qcow2.bz2 -h kvm -F

/root/mysql.sh
/root/iptables.sh
service cloudstack-management restart
EOF
