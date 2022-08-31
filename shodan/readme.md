本案例使用http请求，联动shodan api, 获取查询结果

预备知识
1. go返送Get POST HEAD 请求
    - Get(url)(resp *Response,err error)
    - Head(url)(resp *Response,err error)
    - Post(url string,bodytype string,body io.Reader)(resp *Response,err error) bodytype 为请求头的contenttype  
    - PostForm(url,data url.Values)(resp *Response,err error)  Post发送formdata快捷函数
    - PATCH、PUT、DELETE不支持便捷函数
    - NewRequest(method,url,body)(req *Request,err error)
```
// 案例一 
r1, _ := http.Get("http://www.baidu.com")
r2,_ := http.Head("http://www.baidu.com")
defer r1.Close()
defer r2.Close()

form := url.Values()
form.Add("foo","bar")
r3 := http.Post("http://www.baidu.com","application/x-www-form-urlencode",strings.NewReader(form.Encode()))
defer r3.Close()
// 等效于 http.PostForm("http://www.baidu.com",from)

// 案例二  PATCH等请求构建
req,_ := http.NewRequest("DELETE","http://www.baidu.com",nil)
client := &http.Client{}
// 等效 var client http.Client
client.Do(req)


```
2. ioutil.ReadAll(resp.body) // 获取响应的全部数据
3. json.Unmarshal(data,*store) 或者 json.NewDecoder().Decode()