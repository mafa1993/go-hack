package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/yuin/gopher-lua/lua"
)

const PluginsDir = "./plugins/"
const LuaHttpTypeName = "http"

func main() {
	var (
		l     *lua.LState
		files []os.FileInfo
		err   error
		f     string
	)

	l = lua.NewState() //创建lua.Lstate实例
	defer l.Close()

	register(l)

	if files, err = ioutil.ReadDir(PluginsDir); err != nil {
		panic(err)
	}

	// 遍历所有lua脚本
	for idx := range files {
		fmt.Println("plugin ", files[idx])
		f = fmt.Sprintf("%s%s", PluginsDir, files[idx].Name())

		// 执行lua脚本
		if err := l.DoFile(f); err != nil {
			panic(err)
		}
	}
}

func register(l *lua.LState) {
	mt := l.NewTypeMetatable(LuaHttpTypeName) //相当于lua NewTypeMetatable
	l.SetGlobal("http", mt)                   // 注册mt到lua http全局变量

	// 向mt上注册函数 head 和get
	l.SetField(mt, "head", l.NewFunction(head))
	l.SetField(mt, "get", l.NewFunction(get))
}

// lua函数调用的参数和返回值都存在lua.LState对象中
// @return int 返回lua函数调用返回值的个数
func head(l *lua.LState) int {
	var (
		host string
		port uint64
		path string
		resp *http.Response
		err  error
		url  string
	)

	host = l.CheckString(1)        // 以string接受第一个参数
	port = uint64(l.CheckInt64(2)) // int接受第二个参数

	path = l.CheckString(3)

	url = fmt.Sprintf("http://%s:%d/%s", host, port, path)
	if resp, err = http.Head(url); err != nil {
		l.Push(lua.LNumber(0))                            // 以数字形式存入第一个返回值
		l.Push(lua.LBool(false))                          // 存入第二个参数为false
		l.Push(lua.LString(fmt.Sprintf("error:%s", err))) // 错误信息  第三个返回值
		return 3
	}

	l.Push(lua.LNumber(resp.StatusCode))
	l.Push(lua.LBool(resp.Header.Get("www-Authenticate") != ""))
	l.Push(lua.LString(""))
	return 3
}

func get(l *lua.LState) int {
	var (
		host     string
		port     uint64
		path     string
		resp     *http.Response
		err      error
		url      string
		username string
		password string
		client   *http.Client
		req      *http.Request
	)

	host = l.CheckString(1)        // 以string接受第一个参数
	port = uint64(l.CheckInt64(2)) // int接受第二个参数
	username = l.CheckString(3)
	password = l.CheckString(4)
	path = l.CheckString(5)

	url = fmt.Sprintf("http://%s:%d/%s", host, port, path)

	client = new(http.Client)
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		l.Push(lua.LNumber(0))                            // 以数字形式存入第一个返回值
		l.Push(lua.LBool(false))                          // 存入第二个参数为false
		l.Push(lua.LString(fmt.Sprintf("error:%s", err))) // 错误信息  第三个返回值
		return 3
	}

	if username != "" && password != "" {
		req.SetBasicAuth(username, password)
	}

	if resp, err = client.Do(req); err != nil {
		l.Push(lua.LNumber(0))                            // 以数字形式存入第一个返回值
		l.Push(lua.LBool(false))                          // 存入第二个参数为false
		l.Push(lua.LString(fmt.Sprintf("error:%s", err))) // 错误信息  第三个返回值
		return 3
	}

	l.Push(lua.LNumber(resp.StatusCode))
	l.Push(lua.LBool(false))
	l.Push(lua.LString(""))
	return 3
}
