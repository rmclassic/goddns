package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func main() {
	http.ListenAndServe(":8058", http.HandlerFunc(handler))
}

func handler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	host := q.Get("hostname")
	ip := q.Get("myip")
	auth := r.Header.Get("Authorization")
	authParts := strings.Split(auth, " ")
	if auth == "" || len(authParts) != 2 || authParts[0] != "Basic" {
		fmt.Println("error: empty authorization data")
		w.WriteHeader(401)
		return
	}

	rawCreds := make([]byte, 65535)
	n, err := base64.NewDecoder(base64.StdEncoding, strings.NewReader(authParts[1])).Read(rawCreds)
	if err != nil {
		fmt.Println("error: invalid authorization format")
		w.WriteHeader(401)
		return
	}

	authData := strings.Split(string(rawCreds[:n]), ":")

	if len(authData) != 2 {
		fmt.Println("error: invalid authorization format")
		w.WriteHeader(401)
		return
	}

	fmt.Printf("host %s has ip %s, username %s password %s", host, ip, authData[0], authData[1])
	w.WriteHeader(200)
}
