# ansible+mysql+tomcat
## mysql镜像制作

Make centos6.6 image       
```shell
Packer build centos6-6.json      
```
Install mysql     
```shell
#sed -i 's#SELINUX=enforcing#SELINUX=disabled#' /etc/selinux/config
#yum install -y -q mysql-server mysql mysql-deve
#service mysqld start
#chkconfig mysqld on
#mysql -uroot -e "grant all privileges on *.* to 'root'@'%' identified by 'engine' with grant option;"
#mysql -uroot -e "create database testdb;"
#mysql -uroot -e "use testdb;create table account(id int(4),name char(20));insert into account values(1,'jack');"
```
ip `<192.168.122.68>`      

## Tomcat镜像制作    
Make centos6.6镜像     
```shell
packer build centos6-6.json
```
Install tomcat    
```shell
#sed -i 's#SELINUX=enforcing#SELINUX=disabled#' /etc/selinux/config
#yum install -y -q java-1.8.0-openjdk
#yum install -y -q tomcat6  tomcat6-webapps tomcat6-admin-webapps
#sed -i 's#</tomcat-users>#<role rolename="manager" /><user username="clouder" password="engine" roles="manager" /></tomcat-users>#' /etc/tomcat6/tomcat-users.xml
#service tomcat6 start
#chkconfig tomcat6 on
```
ip `<192.168.122.245>`     

## Webapp deploy
Github clone HelloTomcat      
```shell
#git clone https://github.com/pjq/HelloTomcat.git
#cd HelloTomcat
```
Modify WEB-INF/src/DataManager.java      
```java
import java.util.Properties;
...
        public Connection getConnection() {
-                //String url = "jdbc:mysql://192.168.122.68:3306/testdb";
-                //String username = "root";
-                //String password = "engine";
+InputStream inputStream = this.getClass().getClassLoader().getResourceAsStream("mysql.properties");
+Properties p = new Properties();
+try {
+p.load(inputStream);
+} catch (IOException e1) {
+e1.printStackTrace();
+}
```
Compile DataManager.java
```shell
#javac DataManager.java
#cp DataManager.class ../classes
```
Add WEB-INF/classes/mysql.properties
```text
ip=192.168.122.245 
user=root
passwd=engine
```
Modify WEB-INF/web.xml
```text
        <servlet-mapping>
                <servlet-name>query</servlet-name>
                <url-pattern>/</url-pattern>
        </servlet-mapping>
```
Make war
```shell
#cd HelloTomcat
#jar cvf test6.war *
```
Deploy war in tomcat   
login http://192.168.122.245:8080   
 ![tomcat6 homepage](/images/tomcat6.png) 
Tomcat Manager with user:clouder password:engine  
![tomcat6 managerpage](/images/tomcat6manager.png)
WAR file to deploy select HelloTomcat/test6.war   
now http://192.168.122.245:8080/test6 will show table account list
![tomcat6 webapp page](/images/webapp.png)
## Ansible playbook
This is the whole playbook directory
```list
├── hosts
├── roles
│   ├── mysql
│   │   ├── files
│   │   ├── tasks
│   │   │   └── main.yml
│   │   └── templates
│   │       └── my.cnf
│   └── tomcat
│       ├── files
│       ├── tasks
│       │   └── main.yml
│       └── templates
│           └── mysql.properties
└── roles.yml
```
File hosts include ip and variables
```text
# hosts
[mysql]
192.168.122.68 ansible_ssh_user=root ansible_ssh_pass=engine

[tomcat]
192.168.122.245 ansible_ssh_user=root ansible_ssh_pass=engine

[tomcat:vars]
mysql_ip=192.168.122.245
mysql_user=root
mysql_password=engine
```
File roles.yml
```
- hosts: tomcat
  roles:
    - { role: tomcat }
- hosts: mysql
  roles:
    - { role: mysql }
```
File roles/tomcat/tasks/main.yml
```text
  - name: copy mysql.properties
    template: src=mysql.properties dest=/usr/share/tomcat6/webapps/test6/WEB-INF/classes/mysql.properties
  - name: restart tomcat6
    service: name=tomcat6 state=restarted
```
File roles/tomcat/templates/mysql.properties
```
ip={{mysql_ip}} 
user={{mysql_user}}
passwd={{mysql_password}}
```
File roles/mysql/tasks/main.yml
```
  - name: copy my.cnf
    template: src=my.cnf dest=/etc/my.cnf
  - name: restart mysqld
    service: name=mysqld state=restarted
```
File roles/mysql/templates/my.cnf
```
[mysqld]
datadir=/var/lib/mysql
socket=/var/lib/mysql/mysql.sock
user=mysql
# Disabling symbolic-links is recommended to prevent assorted security risks
symbolic-links=0

[mysqld_safe]
log-error=/var/log/mysqld.log
pid-file=/var/run/mysqld/mysqld.pid
```
When mysql ip changes, change the hosts file and run command
```shell
ansible-playbook -i hosts roles.yml
```
