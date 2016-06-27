#!/usr/bin/python
import libvirt
import os
import sys
import re
import subprocess
from collections import namedtuple

from ansible.parsing.dataloader import DataLoader
from ansible.vars import VariableManager
from ansible.inventory import Inventory
from ansible.executor.playbook_executor import PlaybookExecutor
from ansible.playbook.play import Play
from ansible.executor.task_queue_manager import TaskQueueManager

try:
  conn=libvirt.open("qemu:///system")
  p = conn.lookupByName('testcentos72python')
  p.create()
except (ValueError,libvirt.libvirtError):
  pass  
obj = p.interfaceAddresses(0)
print obj['vnet0']['addrs'][0]['addr']

fo = open("host", "wb")
fo.write(obj['vnet0']['addrs'][0]['addr']+" ansible_ssh_user=root ansible_ssh_pass=engine")
fo.close()

def simpletask():
  variable_manager = VariableManager()
  loader = DataLoader()
  inventory = Inventory(loader=loader, variable_manager=variable_manager, host_list='./host')
  playbook_path = './main.yml'
  if not os.path.exists(playbook_path):
      print '[INFO] The playbook does not exist'
      sys.exit()
  Options = namedtuple('Options', ['listtags', 'listtasks', 'listhosts', 'syntax', 'connection','module_path', 'forks', 'remote_user', 'private_key_file', 'ssh_common_args', 'ssh_extra_args', 'sftp_extra_args', 'scp_extra_args', 'become', 'become_method', 'become_user', 'verbosity', 'check'])
  options = Options(listtags=False, listtasks=False, listhosts=False, syntax=False, connection='ssh', module_path=None, forks=100, remote_user='root', private_key_file=None, ssh_common_args=None, ssh_extra_args=None, sftp_extra_args=None, scp_extra_args=None, become=True, become_method=None, become_user='root', verbosity=None, check=False)
  variable_manager.extra_vars = {'hosts': 'vmware'} # This can accomodate various other command line arguments.`
  passwords = {}
  pbex = PlaybookExecutor(playbooks=[playbook_path], inventory=inventory, variable_manager=variable_manager, loader=loader, options=options, passwords=passwords)
  try:
    results = pbex.run()
  except BaseException:
    pass  
  return

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
   print("playbook "+i+" "+r22)
   simpletask()
   break

