sed -i -e "/:OUTPUT/ a\-A INPUT -p tcp -m tcp --dport 111 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p udp -m udp --dport 111 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p tcp -m tcp --dport 2049 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p udp -m udp --dport 2049 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p tcp -m tcp --dport 2020 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p tcp -m tcp --dport 32803 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p udp -m udp --dport 32769 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p tcp -m tcp --dport 892 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p udp -m udp --dport 892 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p tcp -m tcp --dport 875 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p udp -m udp --dport 875 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p tcp -m tcp --dport 662 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p udp -m udp --dport 662 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p tcp -m tcp --dport 80 -j ACCEPT" /etc/sysconfig/iptables
sed -i -e "/:OUTPUT/ a\-A INPUT -p tcp -m tcp --dport 8096 -j ACCEPT" /etc/sysconfig/iptables
service iptables restart
