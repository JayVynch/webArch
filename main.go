package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Person struct {
	First string
}

type UserClaims struct {
	jwt.RegisteredClaims
	SessionID int64
}

func (u *UserClaims) Valid() error {
	// Check expiration
	if u.ExpiresAt != nil && time.Now().After(u.ExpiresAt.Time) {
		return fmt.Errorf("token is expired")
	}
	// Check NotBefore
	if u.NotBefore != nil && time.Now().Before(u.NotBefore.Time) {
		return fmt.Errorf("token not valid yet")
	}
	// Check IssuedAt
	if u.IssuedAt != nil && time.Now().Before(u.IssuedAt.Time) {
		return fmt.Errorf("token used before issued")
	}
	return nil
}

func main() {

	claims := UserClaims{
		SessionID: 123,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // already expired
		},
	}

	// Validate using built-in method
	if err := claims.Valid(); err != nil {
		fmt.Println("Token invalid:", err)
	} else {
		fmt.Println("Token valid")
	}

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

	// wrongPass := "asdfghjkl"
	err = comparePassword(pass, hashedPwd)
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

// creating a key struct to hold on to keys
type Key struct {
	key       []byte
	createdAt time.Time
}

var currentKey = ""
var keys = map[string]Key{}

func generateNewKey() error {
	newKey := make([]byte, 64)

	_, err := io.ReadFull(rand.Reader, newKey)
	if err != nil {
		return fmt.Errorf("errors while generating new key: %w", err)
	}

	uid, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("errors while generating new kid: %w", err)
	}
	keys[uid.String()] = Key{
		key:       newKey,
		createdAt: time.Now(),
	}

	currentKey = uid.String()

	return nil
}

func createToken(c *UserClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	signedToken, err := token.SignedString(keys[currentKey].key)

	if err != nil {
		return "", fmt.Errorf("could not sign you in: %w", err)
	}

	return signedToken, nil
}

func parseToken(signedToken string) (*UserClaims, error) {
	t, err := jwt.ParseWithClaims(signedToken, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodES512.Alg() {
			return nil, fmt.Errorf("Invalid signing Algorithm")
		}

		kid, ok := t.Header["kid"].(string)

		if !ok {
			return nil, fmt.Errorf("Invalid key")
		}

		k, ok := keys[kid]

		if !ok {
			return nil, fmt.Errorf("invalid key")
		}

		return k, nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error while parsing token %w", err)
	}

	if !t.Valid {
		return nil, fmt.Errorf("token is not valid")
	}
	return t.Claims.(*UserClaims), nil
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
