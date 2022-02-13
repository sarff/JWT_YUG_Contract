package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// For HMAC signing method, the key can be any []byte. It is recommended to generate
// a key using crypto/rand or something equivalent. You need the same key for signing
// and validating.

var userKey = goDotEnvVariable("USER_KEY")
var secret = goDotEnvVariable("SECRET")
var hmacSampleSecret = []byte(secret)

func return_sign() (string, error) {
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_key": userKey,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Unix() + 180,
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)

	//fmt.Println(tokenString, err)

	return tokenString, err
}

func return_token(tokenString string) (string, error) {
	url := goDotEnvVariable("URL")
	method := "POST"

	payload_string := fmt.Sprintf("%s%s%s", "{\n\t\t\"requestToken\": \"", tokenString, "\"}")
	payload := strings.NewReader(payload_string)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var jwt JWT
	json.Unmarshal(body, &jwt)
	link := jwt.JWT.Link

	return link, err
}
