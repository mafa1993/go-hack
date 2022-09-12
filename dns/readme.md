1. 解析ip地址
2. 枚举各种dns记录
3. 使用dns隧道建立c2通道

# 预备知识

1. go标准库中有LookupAddr(addr string)可以查询dns记录，但是无法指定dns服务器
2. 使用第三方饱
    - go get github.com/miekg/dns
3. 使用tcpdump 查看dns解析的交互过程
    - tcpdump -i eth1 -n udp port 53