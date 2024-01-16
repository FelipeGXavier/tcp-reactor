package main

import (
	"eventloop/pkg/reactor"
)

func main() {
	reactor := reactor.MakeReactor(300)
	reactor.Run()
}