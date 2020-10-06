package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("mysupersafephrase")

func homePage(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()
	if err != nil {
		fmt.Fprint(w, err.Error())
		return
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8090", nil)
	if err != nil {
		fmt.Errorf("Error: %s", err.Error())
	}
	req.Header.Set("Token", validToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintf(w, "Error; %s", err.Error())
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	fmt.Fprint(w, string(body))
}

func handlRequest() {
	http.HandleFunc("/", homePage)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["user"] = "Nathan Yang"
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		fmt.Errorf("Something went wrong: %s", err.Error())
		return "", err
	}

	return tokenString, nil
}

func main() {
	fmt.Printf("A simple client...")

	handlRequest()
}
