package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"log"
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
		log.Printf("error: empty authorization data")
		w.WriteHeader(401)
		return
	}

	rawCreds := make([]byte, 65535)
	n, err := base64.NewDecoder(base64.StdEncoding, strings.NewReader(authParts[1])).Read(rawCreds)
	if err != nil {
		log.Printf("error: invalid authorization format")
		w.WriteHeader(401)
		return
	}

	authData := strings.Split(string(rawCreds[:n]), ":")

	if len(authData) != 2 {
		log.Printf("error: invalid authorization format")
		w.WriteHeader(401)
		return
	}

	err = updateDNS(host, ip, authData[0], authData[1])
	if err != nil {
		log.Printf("error: dns update failed: %s", err.Error())
		w.WriteHeader(500)
		return
	}

	log.Printf("host %s has ip %s, username %s password %s", host, ip, authData[0], authData[1])
	w.WriteHeader(200)
}

func updateDNS(host, ip, username, password string) error {
	var req http.Request
	req.URL, _ = url.Parse(fmt.Sprintf("http://dynupdate.no-ip.com/nic/update?hostname=%s&myip=%s", host, password))
	req.SetBasicAuth(username, password)

	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		return err
	}

	io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return fmt.Errorf("invalid http code received %d", resp.StatusCode)
	}

	return nil
}
