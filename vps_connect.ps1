$ErrorActionPreference = "Stop"
$session = New-Object System.Management.Automation.Runspaces.Runspace
$session.Open()

$pipe = $session.CreatePowerShellShell("ssh -i C:/3/ssh_key -o ServerAliveInterval=10 root@150.241.95.205 `"(echo test; sleep 1)`"")

Start-Sleep -Seconds 20

$results = $pipe.ReadToEnd()

$results

$session.Close()