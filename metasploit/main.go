package main

import (
	"fmt"
	"log"
	"msf/rpc"
)

func main() {
	host := "192.168.1.128:55552"
	pass := "123"
	user := "msf"

	if host == "" || pass == "" {
		log.Fatalln("参数错误")
	}

	msf, err := rpc.New(host, user, pass)

	if err != nil {
		log.Panicln(err)
	}

	defer msf.Logout()

	sessions, err := msf.SessionList()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("sessions list =========")
	for _, v := range sessions {
		fmt.Printf("%d,%s\n", v.ID, v.Info)
	}
}
