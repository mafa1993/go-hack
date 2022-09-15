
# 准备工作

1. go1.8以上才支持动态库
2. linux 下动态库编译  go build -buildmode=plugin
3. shard 下  go build -buildmode=c-shared

# go加载动态库的过程

1. 调用plugin.Open(filename) 打开共享对象文件，创建一个*plugin.Plugin实例
2. 在*plugin.Plugin实例上调用Lookup(symbolName string)
3. 使用类型断言将泛型symbol转换为所需类型
4. 根据需要使用生成转换对象
