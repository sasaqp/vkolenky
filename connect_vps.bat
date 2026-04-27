@echo off
echo SSH to VPS - keep this window open
ssh -i C:/3/ssh_key -o ServerAliveInterval=10 root@150.241.95.205
pause