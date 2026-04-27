$pass = ConvertTo-SecureString 'C9HD496nsKu8' -AsPlainText -Force
$cred = New-Object System.Management.Automation.PSCredential('root', $pass)
$session = New-PSSession -ComputerName '150.241.95.205' -Credential $cred -Authentication Basic
Copy-Item -Path 'C:/3/project.zip' -Destination '/root/project.zip' -ToSession $session
Remove-PSSession $session