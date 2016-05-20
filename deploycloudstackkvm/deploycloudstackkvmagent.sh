virsh shutdown testcentos71cskvmagent0
virsh destroy testcentos71cskvmagent0
virsh shutdown testcentos71cskvmagent
virsh destroy testcentos71cskvmagent
rm -f /home/html/downloads/CentOS7-1cskvmagent.qcow2
#wget -P /home/html/downloads/ --no-check-certificate https://192.168.0.82:9090/static/result/20160421103801/output/CentOS7-1.qcow2 
cp /home/code/mycode/go/src/main/static/result/20160421103801/output/CentOS7-1.qcow2 /home/html/downloads/CentOS7-1cskvmagent.qcow2
virsh undefine testcentos71cskvmagent0
virsh define xml/testcentos71agent.xml
virsh start testcentos71cskvmagent0
sleep 1m
testmac=`virsh domiflist testcentos71cskvmagent0|grep network|awk '{print $5}'`
testip=`virsh net-dhcp-leases default|grep $testmac|awk '{print $5}'|cut -d '/' -f 1`
sed -i /$testip/d /root/.ssh/known_hosts
ssh-keyscan $testip >> /root/.ssh/known_hosts
./sshcopyid.exp root@$testip engine
scp xml/ifcfg-eth0.agent root@$testip:/etc/sysconfig/network-scripts/ifcfg-eth0
scp xml/ifcfg-cloudbr0.agent root@$testip:/etc/sysconfig/network-scripts/ifcfg-cloudbr0
scp xml/ifcfg-cloudbr1.agent root@$testip:/etc/sysconfig/network-scripts/ifcfg-cloudbr1
ssh root@$testip "sync"
virsh shutdown testcentos71cskvmagent0
virsh destroy testcentos71cskvmagent0
virsh undefine testcentos71cskvmagent0
sed -i /$testip/d /root/.ssh/known_hosts

virsh shutdown testcentos71cskvmagent
virsh destroy testcentos71cskvmagent
virsh undefine testcentos71cskvmagent
virsh define xml/centos71kvmagent.xml
virsh start testcentos71cskvmagent
sleep 1m
sed -i /10.1.6.20/d /root/.ssh/known_hosts
ssh-keyscan 10.1.6.20 >> /root/.ssh/known_hosts
./sshcopyid.exp root@10.1.6.20 engine
scp /etc/yum.repos.d/mrepo7.repo root@10.1.6.20:/etc/yum.repos.d/mrepo7.repo
scp /etc/yum.repos.d/cloudstack451_7.repo root@10.1.6.20:/etc/yum.repos.d/cloudstack451_7.repo
ssh root@10.1.6.20 << EOF
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
virsh reboot testcentos71cskvmagent
sleep 1m

ssh root@10.1.6.20 << EOF
yum install -y ntp
chkconfig ntpd on
service ntpd start
yum install -y cloudstack-agent
yum -y install vim
sed -i 's/#vnc_listen = "0.0.0.0"/vnc_listen = "0.0.0.0"/g' /etc/libvirt/qemu.conf && sed -i 's/cgroup_controllers=["cpu"]/#cgroup_controllers=["cpu"]/g' /etc/libvirt/qemu.conf
sed -i 's/#listen_tls = 0/listen_tls = 0/g' /etc/libvirt/libvirtd.conf && sed -i 's/#listen_tcp = 1/listen_tcp = 1/g' /etc/libvirt/libvirtd.conf && sed -i 's/#tcp_port = "16509"/tcp_port = "16509"/g' /etc/libvirt/libvirtd.conf && sed -i 's/#auth_tcp = "sasl"/auth_tcp = "none"/g' /etc/libvirt/libvirtd.conf && sed -i 's/#mdns_adv = 1/mdns_adv = 0/g' /etc/libvirt/libvirtd.conf
sed -i 's/#LIBVIRTD_ARGS="--listen"/LIBVIRTD_ARGS="--listen"/g' /etc/sysconfig/libvirtd
sed -i '/cgroup_controllers/d' /usr/lib64/python2.7/site-packages/cloudutils/serviceConfig.py
service libvirtd restart
chkconfig libvirtd on
service cloudstack-agent restart
EOF
