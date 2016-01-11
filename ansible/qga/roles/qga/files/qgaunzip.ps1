Import-Module BitsTransfer
$url='http://192.168.0.82:9090/static/win/qemu-ga.zip'
Start-BitsTransfer $url C:\win\qemu-ga.zip

Function Unzip-File()
{
    param([string]$ZipFile,[string]$TargetFolder)
    if(!(Test-Path $TargetFolder))
    {
        mkdir $TargetFolder
    }
    $shellApp = New-Object -ComObject Shell.Application
    $files = $shellApp.NameSpace($ZipFile).Items()
    $shellApp.NameSpace($TargetFolder).CopyHere($files)
}
Unzip-File -ZipFile C:\win\qemu-ga.zip -TargetFolder "C:\Program Files"
