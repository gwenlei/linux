cd C:\
$src = 'http://192.168.0.82:9090/static/win/qemu-ga.zip'
$des = "$env:win/qemu-ga.zip"
Invoke-WebRequest -uri 'http://192.168.0.82:9090/static/win/qemu-ga.zip' -OutFile "C:\win\qemu-ga.zip"
Unblock-File $des

