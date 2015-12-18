# install keepalived and ipvsadm
```
yum install mysql-server mysql mysql-deve keepalived ipvsadm
sed -i 's#net.ipv4.ip_forward = 0#net.ipv4.ip_forward = 1#' /etc/sysctl.conf
```
# Mysql + KeepAliaved 搭建 Mysql 双 Master 集群的步骤
## 前提假设：    
所需的软件 mysql、keepalived、ipvsadm 已经安装     
	$myip1 -> 第一台mysql所在的ip     
	$myip2 -> 第二台mysql所在的ip     
	$myuser -> mysql用户名      
	$mypwd -> mysql用户密码      
	两台mysql的用户名和密码是一样的     
	$vip -> mysql 集群虚拟IP。注意，虚拟ip应该与 $myip1、$myip2 在同一网段    
	$inteface -> 虚拟机 Real ip 所在的网卡，例如eth0      
实际环境赋值: centos6.6\mysql-5.1.73\keepalived-1.2.13       
$myip1=192.168.122.68     
$myip2=192.168.122.245    
$myuser=slaver    
$mypwd=engine    
$vip=192.168.122.100/24      
$inteface=eth0    
	 
## 操作步骤：
### 第一步:  配置 Mysql Replication slave 权限     		
####1) 在 $myip1 的 mysql 中执行如下sql      

```
grant replication slave on *.* to '$myuser'@'$myip2' identified by '$mypwd';
flush privileges;
```

####2) 在 $myip2 的 mysql 中执行如下sql      

```
grant replication slave on *.* to '$myuser'@'$myip1' identified by '$mypwd';
flush privileges;
```

### 第二步:  修改 /etc/my.cnf      

####1) 在 $myip1 的 /etc/my.cnf 添加如下配置：     

```
[mysqld]
######For dual master#####
server-id=1
log-bin=master-bin
log-slave-updates
binlog-ignore-db=mysql
binlog-ignore-db=information_schema
replicate-ignore-db=mysql
replicate-ignore-db=information_schema
auto_increment_increment=2
auto_increment_offset=1
###########
```

####2) 在 $myip2 的 /etc/my.conf 添加如下配置：       

```
[mysqld]
######For dual master#####
server-id=2
log-bin=master-bin
log-slave-updates
binlog-ignore-db=mysql
binlog-ignore-db=information_schema
replicate-ignore-db=mysql
replicate-ignore-db=information_schema
auto_increment_increment=2
auto_increment_offset=2
###########
```

###第三步: 重启mysql       

```
ssh  root@$myip1 /etc/init.d/mysqld restart
ssh  root@$myip2 /etc/init.d/mysqld restart
```

###第四步: 启动 MYSQL Slave 线程      

####1) 在 $myip1 的 mysql 中执行如下sql      

```		
CHANGE MASTER TO MASTER_HOST='$myip2';
CHANGE MASTER TO MASTER_USER='$myuser';
CHANGE MASTER TO MASTER_PASSWORD='$mypwd';
start slave;
show slave status\G
```

####2) 在 $myip2 的 mysql 中执行如下sql      

```
CHANGE MASTER TO MASTER_HOST='$myip1';
CHANGE MASTER TO MASTER_USER='$myuser';
CHANGE MASTER TO MASTER_PASSWORD='$mypwd';
start slave;
show slave status\G
```

###第五步： 配置 keepalived.conf 文件        

####1) 拷贝 check_MySQL.sh 文件到两个mysql节点:         
```
#!/bin/bash

###检查mysql服务是否存在###

MYSQL_HOST=localhost
MYSQL_USER=root
CHECK_COUNT=5

counter=1
while true
do
        mysql -h $MYSQL_HOST -u $MYSQL_USER -e "show status;"  >/dev/null 2>&1
        i=$?
        ps aux | grep mysqld | grep -v grep > /dev/null 2>&1
        j=$?
        if [ $i = 0 ] && [ $j = 0 ]
        then
                exit 0
        else
                if [ $i = 1 ] && [ $j = 0 ]
                then
                        exit 0
                else
                        if [ $counter -gt $CHECK_COUNT ]
                        then
                                break
                        fi
                let counter++
                continue
                fi
        fi
        done
/etc/init.d/keepalived stop
exit 1
```
```
scp check_MySQL.sh root@$myip1:/etc/keepalived/check_MySQL.sh
scp check_MySQL.sh root@$myip2:/etc/keepalived/check_MySQL.sh
chmod +x /etc/keepalived/check_MySQL.sh
sh /etc/keepalived/check_MySQL.sh
```

####2) 修改 $myip1 的 keepalived.conf 内容为：        

```		
! Configuration File for keepalived
global_defs {
   router_id Keepalived_MySQL
}
 
vrrp_script check_run {
    script "/etc/keepalived/check_MySQL.sh"
    interval 5
}
 
vrrp_sync_group VG1{
    group {
	VI_1
    }
}
 
vrrp_instance VI_1 {
    state BACKUP
    #修改为 real ip 对应的网卡。
    interface $inteface
    #同一个vrrp_stance,MASTER和BACKUP的virtual_router_id是一致的，同时在整个vrrp内是唯一的。
    virtual_router_id 51
    #一般MASTER可以设置大一些。
    priority 100
    advert_int 1
    nopreempt
    track_script {
	check_run
    }
    authentication {
	auth_type PASS
	auth_pass 111
    }
    virtual_ipaddress {
	$vip
    }
}
```

####3) 修改 $myip2 的 keepalived.conf 内容为：       

```
! Configuration File for keepalived
global_defs {
   router_id Keepalived_MySQL
}
 
vrrp_script check_run {
    script "/etc/keepalived/check_MySQL.sh"
    interval 5
}
 
vrrp_sync_group VG1{
    group {
	VI_1
    }
}
 
vrrp_instance VI_1 {
    state BACKUP
    #修改为 real ip 对应的网卡。
    interface $inteface
    #同一个vrrp_stance,MASTER和BACKUP的virtual_router_id是一致的，同时在整个vrrp内是唯一的。
    virtual_router_id 51
    #一般MASTER可以设置大一些。
    priority 50
    advert_int 1
    nopreempt
    track_script {
	check_run
    }
    authentication {
	auth_type PASS
	auth_pass 111
    }
    virtual_ipaddress {
	$vip
    }
}
```

###第六步： 重启 keepalived 文件:

```
ssh  root@$myip1 service keepalived restart
ssh  root@$myip2 service keepalived restart
```
