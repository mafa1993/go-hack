# go pack数据包处理

1. 明文身份验证
2. syn扫描
3. syn flood

# 准备工作

1. go get github.com/google/gopacket 
    - 支持berkeley数据包过滤
    - 依赖libpcap-dev apt-get intsall libpcap-dev
2. bpf语法 https://tcpdump.org/manpages/pacp-filter.7.html
