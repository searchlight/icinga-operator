package main

import "github.com/jpillora/go-ogle-analytics"

func main() {
	client, err := ga.NewClient("UA-30305960-4")
	if err != nil {
		panic(err)
	}

	err = client.Send(ga.NewEvent("Foo", "Bar").Label("Bazz"))
	if err != nil {
		panic(err)
	}

	println("Event fired!")
}
