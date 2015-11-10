for i in $(yum search collectd |grep collectd |awk '{print $1}'|sed -e "s/\..*//g" )
do
 echo $i
 yumdownloader --resolve $i
done

