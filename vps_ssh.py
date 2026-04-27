import paramiko
import sys

key = paramiko.Ed25519Key(filename='C:/3/ssh_key')
client = paramiko.SSHClient()
client.set_missing_host_key_policy(paramiko.AutoAddPolicy())

try:
    client.connect('150.241.95.205', username='root', pkey=key, 
                 banner_timeout=30, auth_timeout=30,
                 server_keepalive=paramiko.channel.Channel.getdefault(timeout=10))
    stdin, stdout, stderr = client.exec_command('docker ps', timeout=30)
    print(stdout.read().decode())
    client.close()
except Exception as e:
    print(f"Error: {e}")
    sys.exit(1)