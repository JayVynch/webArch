package main

import (
	"crypto/hmac"
	"crypto/sha512"
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
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

	// basic auth standards
	// fmt.Println(base64.StdEncoding.EncodeToString([]byte("user:pass")))

	pass := "123456789"

	hashedPwd, err := hashPassword(pass)

	if err != nil {
		fmt.Println("Opps!! something went wrong:", err)
	}

	fmt.Println("HashedPassword is:", string(hashedPwd))

	wrongPass := "asdfghjkl"
	err = comparePassword(wrongPass, hashedPwd)
	if err != nil {
		log.Fatalln("In correct Password")
	}

	log.Println("Yay you are logged in")

	// serving a web application with go
	// http.HandleFunc("/encoder", foo)
	// http.HandleFunc("/decoder", bar)
	// http.ListenAndServe(":8000", nil)
}

func hashPassword(password string) ([]byte, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, fmt.Errorf("error while generating password hash: %w", err)
	}

	return bs, nil
}

func comparePassword(password string, hashedPassword []byte) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil {
		return fmt.Errorf("password does not match: %w", err)
	}
	return nil
}

func signMsg(msg []byte) ([]byte, error) {
	var key []byte

	for i := 1; i < 64; i++ {
		key = append(key, byte(i))
	}

	h := hmac.New(sha512.New, key)

	_, err := h.Write(msg)
	if err != nil {
		return nil, fmt.Errorf("error while hashing message: %w", err)
	}

	signature := h.Sum(nil)
	return signature, nil
}

func checkSig(msg, sig []byte) (bool, error) {
	newSig, err := signMsg(msg)

	if err != nil {
		return false, fmt.Errorf("we encountered some issues while hashing message: %w", err)
	}

	same := hmac.Equal(newSig, sig)
	return same, nil
}

// func foo(w http.ResponseWriter, r *http.Request) {
// 	p1 := Person{
// 		First: "James",
// 	}

// 	err := json.NewEncoder(w).Encode(p1)
// 	if err != nil {
// 		log.Println("Bad json", err)
// 	}
// }

// func bar(w http.ResponseWriter, r *http.Request) {
// 	var p1 Person

// 	err := json.NewDecoder(r.Body).Decode(&p1)

// 	if err != nil {
// 		log.Println("Bad json", err)
// 	}

// 	log.Println("decoded Person:", p1)
// }
