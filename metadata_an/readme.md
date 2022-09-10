
# 准备工作

1. go get github.com/puerkitobio/goquery   使用goquery包解析html文档
2. xls和doc文档解压后，会生成app.xml 和core.xml 文件，里面记录了一些文档的信息
3. bing搜索的高级应用 site:xxx.com && filetype:docx && instreamset:(url title):docx
    - site:过滤特定域的结果
    - filetype 根据资源类型过滤
    - instreamset 用于过滤结果

# 基础知识

1. xml.NewDecoder(filecontent).Decode(rlt)    xml 解析

用法

go run main.go -site "查询的site" -filetype "文件类型 docx等"