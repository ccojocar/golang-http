package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const TOKEN_HEADER = "X-Forwarded-Access-Token"
const DEX_USERINFO = "https://dex.sso.jx.cosmin.rawlings.it/userinfo"

func handler(w http.ResponseWriter, r *http.Request) {
	token, err := parseAccessToken(r)
	if err != nil {
		printError(w, err)
		return
	}

	userinfo, err := getUserinfo(token)
	if err != nil {
		printError(w, err)
		return
	}

	fmt.Fprintf(w, "User Info:\n%s", userinfo)
}

func printError(w http.ResponseWriter, err error) {
	fmt.Fprintf(w, "Error: %v", err)
}

func parseAccessToken(r *http.Request) (string, error) {
	for k, v := range r.Header {
		if k == TOKEN_HEADER {
			return v[0], nil
		}
	}
	return "", errors.New("no token found")
}

func getUserinfo(accessToken string) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest("GET", DEX_USERINFO, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "bearer "+accessToken)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return "", fmt.Errorf("status code: %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
