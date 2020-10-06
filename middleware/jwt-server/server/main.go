package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("mysupersafephrase")

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Super safe information")
}

func isAuthed(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("error in parsing JWT token")
				}
				return mySigningKey, nil
			})
			if err != nil {
				fmt.Fprintf(w, err.Error())
				return
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not Authorized")
		}
	})
}

func handleRequest() {
	http.Handle("/", isAuthed(homePage))

	log.Fatal(http.ListenAndServe(":8090", nil))
}

func main() {
	fmt.Printf("Simple server...")
	handleRequest()
}
