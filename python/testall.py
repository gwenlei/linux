#!/usr/bin/python
import libvirt
import os
import sys
import re
import subprocess
import time
import shutil
from collections import namedtuple

from ansible.parsing.dataloader import DataLoader
from ansible.vars import VariableManager
from ansible.inventory import Inventory
from ansible.executor.playbook_executor import PlaybookExecutor
from ansible.playbook.play import Play
from ansible.executor.task_queue_manager import TaskQueueManager

try:
  conn=libvirt.open("qemu:///system")
  pv = conn.lookupByName('testubuntu1604grub3')
  if pv.state()[0]==1 :
     print("vm shutdowning")
     pv.shutdown()
     time.sleep(30)
  cmd='qemu-img create -f qcow2 /home/html/downloads/Ubuntu16-04\(2\).qcow2 -b /home/html/downloads/Ubuntu16-04\(4\).qcow2 '
  p=subprocess.Popen(args=cmd,shell=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,close_fds=True)
  (stdoutdata,stderrdata)=p.communicate()
  if p.returncode !=0:
     print (cmd+'error')
  if pv.state()[0]<>1 :
     print("vm starting")
     pv.create()
     time.sleep(30)
except (ValueError,libvirt.libvirtError):
  print("vm error")
  pass
print("vm running")  
obj = pv.interfaceAddresses(0)
print obj[obj.keys()[0]]['addrs'][0]['addr']

fo = open("hostu2", "wb")
fo.write(obj[obj.keys()[0]]['addrs'][0]['addr']+" ansible_ssh_user=root ansible_ssh_pass=engine")
#fo.write(obj[obj.keys()[0]]['addrs'][0]['addr'])
fo.close()

fa = open("allrecord2", "w")
fa.close()
fs = open("successrecord2", "w")
fs.close()

cmd='ansible-galaxy search mysql'
p=subprocess.Popen(args=cmd,shell=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,close_fds=True)
(stdoutdata,stderrdata)=p.communicate()
if p.returncode !=0:
 print (cmd+'error')

booklist=[]
i=0
for r in str(stdoutdata).split("\n"):
 i+=1
 if i<6 :
   if i==2 : print (r) 
   continue
 if i>200 : break
 j=0
 for r2 in r.split(" "):
   if j==0 : 
     j+=1
     continue
   booklist.append(r2)
   break

i=0
for r2 in booklist :
   i+=1
   cmd='ansible-galaxy install '+r2
   print cmd + '  '+str(i)
   p=subprocess.Popen(args=cmd,shell=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,close_fds=True)
   (stdoutdata,stderrdata)=p.communicate()
   print str(stdoutdata)
   print p.returncode
   if p.returncode !=0:
      print (cmd+' error')
i=0
for r2 in booklist :
      i+=1
      time.sleep(3)
      foo = open("main2.yml", "wb")
      foo.write("- hosts: all\n  roles:\n  - "+r2)
      foo.close()
      print("playbook "+str(i)+" "+r2)
      cmd='ansible-playbook -i hostu2 main2.yml'
      print cmd
      p=subprocess.Popen(args=cmd,shell=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,close_fds=True)
      (stdoutdata,stderrdata)=p.communicate()
      print str(stdoutdata)
      result=p.returncode
      fa = open("allrecord2", "a")
      fa.write("playbook "+str(i)+" "+r2+" result:"+str(result)+"\n")
      fa.close()
      if result==0 :
         fs = open("successrecord2", "a")
         fs.write("playbook "+str(i)+" "+r2+"\n")
         fs.close()
#      else :
#         cmd='ansible-galaxy remove '+r2
#         print cmd
#         p=subprocess.Popen(args=cmd,shell=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,close_fds=True)
#         (stdoutdata,stderrdata)=p.communicate()
#         print str(stdoutdata)



