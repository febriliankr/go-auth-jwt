package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("supersecretstringphrase")

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "super secret key")
}

func isAuthorized(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {

			fmt.Println(r.Header["Token"][0]) // ✅ outputnya beneran sesuai Token yang gua masukin lewat header

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if jwt.GetSigningMethod("HS256") != token.Method {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				fmt.Println(mySigningKey)

				return mySigningKey, nil
			})

			if err != nil {
				fmt.Fprint(w, err.Error(), " :(")
			}

			if token.Valid {
				endpoint(w, r)
			}
		} else {
			fmt.Fprintf(w, "Not authorized")
		}
	})
}

func handleRequests() {
	http.Handle("/", isAuthorized(homePage))
	log.Fatal(http.ListenAndServe(":9001", nil))
}

func main() {
	handleRequests()
}
