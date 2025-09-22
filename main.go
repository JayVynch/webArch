package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	// format to json
	bs, err := json.Marshal(xp)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println(string(bs))

	xp2 := []Person{}
	err = json.Unmarshal(bs, &xp2)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Return to go structure", xp2)

	// serving a web application with go
	http.HandleFunc("/encoder", foo)
	http.HandleFunc("/decoder", bar)
	http.ListenAndServe(":80", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {

}

func bar(w http.ResponseWriter, r *http.Request) {

}
