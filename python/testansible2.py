import subprocess
print 'start'
cmd='ansible-galaxy search mysql'
p=subprocess.Popen(args=cmd,shell=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,close_fds=True)
(stdoutdata,stderrdata)=p.communicate()
if p.returncode !=0:
 print (cmd+'error')

i=0
for r in str(stdoutdata).split("\n"):
 if i<5 : 
   i+=1
   continue
 j=0
 for r2 in r.split(" "):
   if j==0 : 
     j+=1
     continue
   cmd2='ansible-galaxy install '+r2
   print cmd2
   p2=subprocess.Popen(args=cmd2,shell=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,close_fds=True)
   (stdoutdata2,stderrdata2)=p2.communicate()
   if p2.returncode !=0:
      print (cmd2+' error')
   break

print(p.returncode)
print 'end'
