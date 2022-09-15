package main

import (
	"fmt"
	"goplugin/scanner"
	"io/ioutil"
	"log"
	"os"
	"plugin"
)

// 插件调用

const PluginDir = "./plugins/" // 定义放.so 文件的目录

func main() {
	var (
		files []os.FileInfo
		err   error
		p     *plugin.Plugin
		n     plugin.Symbol
		check scanner.Checker
		res   *scanner.Result
	)

	// 获取插件目录所有文件
	if files, err = ioutil.ReadDir(PluginDir); err != nil {
		log.Fatalln(err)
	}

	// 遍历所有的插件
	for idx := range files {
		fmt.Println("插件名", files[idx].Name())

		// 读取插件
		if p, err = plugin.Open(PluginDir + files[idx].Name()); err != nil {
			log.Fatal(err)
		}

		// 加载插件中的New函数
		if n, err = p.Lookup("New"); err != nil {
			log.Fatalln(err)
		}

		newFunc, ok := n.(func() scanner.Checker) //类型断言，检查获取到的n的类型，并返回

		if !ok {
			log.Fatal("new 函数检查出错")
		}

		check = newFunc() // 调用newfunc

		res = check.Check("202.108.22.103", 80) // 调用插件中接口的check方法

		if res.Vulnerable {
			log.Println(res.Details)
		} else {
			log.Println("host not avalable")
		}
	}

}
