 mkdir %SystemDrive%\win

 :: Fetch netframework4.0
bitsadmin /create myDownloadJob3
bitsadmin /addfile myDownloadJob3 https://192.168.0.82:9090/static/win/NDP46-KB3045557-x86-x64-AllOS-ENU.exe %SystemDrive%\win\NDP46-KB3045557-x86-x64-AllOS-ENU.exe
bitsadmin /addfile myDownloadJob3 https://192.168.0.82:9090/static/win/uppowershell.ps1 %SystemDrive%\win\uppowershell.ps1
bitsadmin /addfile myDownloadJob3 https://192.168.0.82:9090/static/win/winrm.ps1 %SystemDrive%\win\winrm.ps1
bitsadmin /addfile myDownloadJob3 https://192.168.0.82:9090/static/upload/qemu-ga-x64.msi %SystemDrive%\win\qemu-ga-x64.msi
bitsadmin /SetSecurityFlags myDownloadJob3 30
bitsadmin /resume myDownloadJob3
ping -n 120 127.0.0.1 >nul
bitsadmin /complete myDownloadJob3


%SystemDrive%\win\NDP46-KB3045557-x86-x64-AllOS-ENU.exe /q
::msiexec.exe /i "C:\Users\Administrator\Downloads\CloudInstanceManager(1).msi" /passive

 cd %SystemDrive%\win
 powershell Set-ExecutionPolicy Unrestricted
 powershell .\uppowershell.ps1
 powershell .\winrm.ps1

echo ==^> Disabling UAC
reg add "HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Policies\System" /f /v EnableLUA /t REG_DWORD /d 0

netsh advfirewall firewall set rule group="File and Printer Sharing" new enable=yes

reg add "HKEY_LOCAL_MACHINE\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon" /v AutoAdminLogon /d 0 /f

echo ==^> enable remote desktop
netsh advfirewall firewall add rule name="Open Port 3389" dir=in action=allow protocol=TCP localport=3389
reg add "HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\Terminal Server" /v fDenyTSConnections /t REG_DWORD /d 0 /f

# allow RDP on firewall
Enable-NetFirewallRule -DisplayName 'Remote Desktop - User Mode (TCP-in)'


import-module Servermanager
add-windowsfeature NPAS-Policy-Server,NPAS-RRAS-Services
add-windowsfeature Web-Static-Content,Web-Default-Doc,Web-Dir-Browsing,Web-Http-Errors,Web-Http-Logging,Web-Request-Monitor,Web-Filtering,Web-Stat-Compression,Web-Mgmt-Console,Web-Ftp-Server
add-windowsfeature WDS

