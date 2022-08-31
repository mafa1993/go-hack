package main

import "shodan/shodan"

func main() {
	client := shodan.New("V3RD7hhexN80bnD7He7pkz7UKjeoRewh")
	client.APIInfo()
}
