Import-Module BitsTransfer
$url='http://192.168.0.82:9090/static/win/qemu-ga.zip'
Start-BitsTransfer $url C:\win\qemu-ga.zip

