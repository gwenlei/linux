yum install -y -q mysql-server mysql mysql-deve
service mysqld start
chkconfig mysqld on
mysql -uroot -e "grant all privileges on *.* to 'root'@'%' identified by 'engine' with grant option;"
mysql -uroot -e "create database testdb;"
mysql -uroot -e "create table account(id int(4),name char(20));"
mysql -uroot -e "insert into account values(1,'jack');"
