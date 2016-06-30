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

#shutil.copyfile("/home/html/downloads/CentOS7-2python2.qcow2", "/home/html/downloads/CentOS7-2python.qcow2")
try:
  conn=libvirt.open("qemu:///system")
  p = conn.lookupByName('testubuntu1604grub4')
  if p.state()[0]<>1 :
     p.create()
     time.sleep(30)
except (ValueError,libvirt.libvirtError):
  pass  
obj = p.interfaceAddresses(0)
print obj[obj.keys()[0]]['addrs'][0]['addr']

fo = open("hostu", "wb")
fo.write(obj[obj.keys()[0]]['addrs'][0]['addr']+" ansible_ssh_user=root ansible_ssh_pass=engine")
fo.close()

def simpletask():
  variable_manager = VariableManager()
  loader = DataLoader()
  inventory = Inventory(loader=loader, variable_manager=variable_manager, host_list='./hostu')
  playbook_path = './main.yml'
  if not os.path.exists(playbook_path):
      print '[INFO] The playbook does not exist'
      sys.exit()
  Options = namedtuple('Options', ['listtags', 'listtasks', 'listhosts', 'syntax', 'connection','module_path', 'forks', 'remote_user', 'private_key_file', 'ssh_common_args', 'ssh_extra_args', 'sftp_extra_args', 'scp_extra_args', 'become', 'become_method', 'become_user', 'verbosity', 'check'])
  options = Options(listtags=False, listtasks=False, listhosts=False, syntax=False, connection='ssh', module_path=None, forks=100, remote_user='root', private_key_file=None, ssh_common_args=None, ssh_extra_args=None, sftp_extra_args=None, scp_extra_args=None, become=True, become_method=None, become_user='root', verbosity=None, check=False)
  variable_manager.extra_vars = {'hosts': 'vmware'} # This can accomodate various other command line arguments.`
  passwords = {}
  pbex = PlaybookExecutor(playbooks=[playbook_path], inventory=inventory, variable_manager=variable_manager, loader=loader, options=options, passwords=passwords)
  results =99
  try:
    results = pbex.run()
#    cb = ResultAccumulator()
#    pbex._tqm._stdout_callback = cb
  except BaseException:
    pass  
  return results

cmd='ansible-galaxy list'
p=subprocess.Popen(args=cmd,shell=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,close_fds=True)
(stdoutdata,stderrdata)=p.communicate()
if p.returncode !=0:
 print (cmd+'error')

i=0
for r in str(stdoutdata).split("\n"):
 j=0
 for r2 in r.split(" "):
   if j==0 : 
     j+=1
     continue
   r22=re.sub(r',',"",r2)
   foo = open("main.yml", "wb")
   foo.write("- hosts: all\n  roles:\n  - "+r22)
   foo.close()
   print("playbook "+str(i)+" "+r22)
   result=simpletask()
   fa = open("allrecord", "a")
   fa.write("playbook "+str(i)+" "+r22+" result:"+str(result)+"\n")
   fa.close()
   if result==0 :
      fs = open("successrecord", "a")
      fs.write("playbook "+str(i)+" "+r22+"\n")
      fs.close()
   i+=1
   break

