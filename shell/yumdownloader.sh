str=("nginx" "golang" "rsync" "net-tools" "tree" "ansible" "yum-utils" ) 
for i in ${str[@]} 
do 
 echo $i 
yumdownloader --resolve $i --destdir /usr/share/nginx/html/repo/centos71-repo/packages  
 
done 
 
createrepo /usr/share/nginx/html/repo/centos71-repo 
