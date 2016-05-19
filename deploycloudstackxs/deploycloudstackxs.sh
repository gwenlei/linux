virsh shutdown testcentos71
virsh destroy testcentos71
virsh shutdown testcentos71cloudstack
virsh destroy testcentos71cloudstack
rm -f /home/html/downloads/CentOS7-1.qcow2
#wget -P /home/html/downloads/ --no-check-certificate https://192.168.0.82:9090/static/result/20160421103801/output/CentOS7-1.qcow2 
cp /home/code/mycode/go/src/main/static/result/20160421103801/output/CentOS7-1.qcow2 /home/html/downloads/CentOS7-1.qcow2
virsh undefine testcentos71
virsh define xml/testcentos71.xml
virsh start testcentos71
sleep 1m
testmac=`virsh domiflist testcentos71|grep network|awk '{print $5}'`
testip=`virsh net-dhcp-leases default|grep $testmac|awk '{print $5}'|cut -d '/' -f 1`
sed -i /$testip/d /root/.ssh/known_hosts
ssh-keyscan $testip >> /root/.ssh/known_hosts
./sshcopyid.exp root@$testip engine
scp xml/ifcfg* root@$testip:/etc/sysconfig/network-scripts/
ssh root@$testip "sync"
virsh shutdown testcentos71
virsh destroy testcentos71
sed -i /$testip/d /root/.ssh/known_hosts

virsh shutdown testcentos71cloudstack
virsh destroy testcentos71cloudstack
virsh undefine testcentos71cloudstack
virsh define xml/centos71cloudstack.xml
virsh start testcentos71cloudstack
sleep 1m
sed -i /192.168.56.11/d /root/.ssh/known_hosts
ssh-keyscan 192.168.56.11 >> /root/.ssh/known_hosts
./sshcopyid.exp root@192.168.56.11 engine
scp /etc/yum.repos.d/mrepo7.repo root@192.168.56.11:/etc/yum.repos.d/mrepo7.repo
scp /etc/yum.repos.d/cloudstack451_7.repo root@192.168.56.11:/etc/yum.repos.d/cloudstack451_7.repo
scp /home/code/mycode/deploycloudstackxs/sshcopyid.exp root@192.168.56.11:/root
scp /home/code/mycode/deploycloudstackxs/mysqlsecure.sh root@192.168.56.11:/root
scp /home/code/mycode/deploycloudstackxs/mysql.sh root@192.168.56.11:/root
scp /home/code/mycode/deploycloudstackxs/iptables.sh root@192.168.56.11:/root
scp /home/html/downloads/systemvm64template-4.5-xen.vhd.bz2 root@192.168.56.11:/root
ssh root@192.168.56.11 << EOF
chmod +x mysql.sh mysqlsecure.sh iptables.sh sshcopyid.exp
cd /etc/yum.repos.d/
rename .repo .repoback CentOS*
setenforce permissive
sed -i "/SELINUX=enforcing/ c\SELINUX=permissive" /etc/selinux/config
echo "127.0.0.1 localhost localhost.cstack.local" >> /etc/hosts
echo "192.168.56.11 csman.cstack.local csman" >> /etc/hosts
echo "192.168.56.101 xenserver.cstack.local xenserver" >> /etc/hosts
sed -i "/#UseDNS yes/ c\UseDNS no" /etc/ssh/sshd_config
service sshd restart
EOF
virsh reboot testcentos71cloudstack
sleep 1m

ssh root@192.168.56.11 << EOF
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
cd /usr/share/cloudstack-common/scripts/vm/hypervisor/xenserver/
wget http://192.168.0.82/vhd-util
chmod 755 /usr/share/cloudstack-common/scripts/vm/hypervisor/xenserver/vhd-util
#/usr/share/cloudstack-common/scripts/storage/secondary/cloud-install-sys-tmplt -m /exports/secondary -u http://192.168.0.82/downloads/systemvm64template-4.5-xen.vhd.bz2 -h xenserver -F
/usr/share/cloudstack-common/scripts/storage/secondary/cloud-install-sys-tmplt -m /exports/secondary -u file:///root/systemvm64template-4.5-xen.vhd.bz2 -h xenserver -F

/root/mysql.sh
/root/iptables.sh
service cloudstack-management restart
EOF
