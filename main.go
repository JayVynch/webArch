package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Person struct {
	First string
}

func main() {
	// p1 := Person{
	// 	First: "James",
	// }

	// p2 := Person{
	// 	First: "Daniel",
	// }

	// xp := []Person{p1, p2}

	// // format to json
	// bs, err := json.Marshal(xp)

	// if err != nil {
	// 	log.Panic(err)
	// }

	// fmt.Println(string(bs))

	// xp2 := []Person{}
	// err = json.Unmarshal(bs, &xp2)

	// if err != nil {
	// 	log.Panic(err)
	// }

	// fmt.Println("Return to go structure", xp2)

	// serving a web application with go
	http.HandleFunc("/encoder", foo)
	http.HandleFunc("/decoder", bar)
	http.ListenAndServe(":8000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	p1 := Person{
		First: "James",
	}

	err := json.NewEncoder(w).Encode(p1)
	if err != nil {
		log.Println("Bad json", err)
	}
}

func bar(w http.ResponseWriter, r *http.Request) {
	var p1 Person

	err := json.NewDecoder(r.Body).Decode(&p1)

	if err != nil {
		log.Println("Bad json", err)
	}

	log.Println("decoded Person:", p1)
}
