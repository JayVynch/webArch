package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Person struct {
	First string
}

func main() {
	p1 := Person{
		First: "James",
	}

	p2 := Person{
		First: "Daniel",
	}

	xp := []Person{p1, p2}

	bs, err := json.Marshal(xp)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(bs))
}
