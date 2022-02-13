package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func WriteToFile(TextToWrite []byte) {
	fileName := "price.json"
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
  "format": "json",
  "type": "regular",
  "cats": [498,177,552,196,1463,197,1217,1215,85,218,129,494,142,493,98,99,102,103,1247,560,1421,492,118,111,225,1514,106,1422,1037,1036,414,354,1460,367,372,601,345,347,1394,1448,1311,1438,1302,1449,1280,1211,1279,1450,1312,1456,1487,1512,1457,1301,404,392,1439,390,395,1446,254,252,1441,250,256,796,58,1187,56,51,50,233,244,231],
  "ext_cols": ['artikul','descr','url','bigimg'],
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
