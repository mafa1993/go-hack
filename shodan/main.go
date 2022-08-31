package main

import (
	"fmt"
	"log"
	"os"
	"shodan/shodan"
)

func main() {

	//client.APIInfo()
	//client.HostSearch("thinkphp")
	if len(os.Args) != 2 {
		log.Fatalln("参数不对")
	}

	//shodan的api放到环境变量里
	apiKey := os.Getenv("SHODAN_API_KEY")

	s := shodan.New(apiKey)

	info, err := s.APIInfo()

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("query credits:%d,\nnScan credits:%d\n", info.QueryCredits, info.ScanCredits)

	HostSearch, err := s.HostSearch(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}

	for _, host := range HostSearch.Matches {
		fmt.Printf("%18s,%d\n", host.IPString, host.Port)
	}
}
