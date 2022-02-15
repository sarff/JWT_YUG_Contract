package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func WriteToFile(TextToWrite []byte) {
	fileName := `C:\PriceYUG\price.json` // for other system without path
	//var d = TextToWrite
	err := ioutil.WriteFile(fileName, TextToWrite, 0666)
	if err != nil {
		fmt.Println("write fail")
	}
	fmt.Println("write success")
}

func getPrice(TokenBraer string) {

	url := goDotEnvVariable("URL_PRICE")
	method := "POST"

	payload := strings.NewReader(`{
  		"type": ` + goDotEnvVariable("TYPE_PROD") + `,
  		"cats": [` + goDotEnvVariable("CATS") + `],
  		"ext_cols": [` + goDotEnvVariable("EXT_COLS") + `],
  		"type_prod": "new"
		}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+TokenBraer)
	req.Header.Add("Content-Type", "text/plain")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	WriteToFile(body)
}
