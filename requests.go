package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func download(url string, catalog string) {
	fileName := catalog + "/" + url[strings.LastIndex(url, "/")+1:]
	_, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			os.Mkdir(catalog, 0755)
			output, err := os.Create(fileName)
			defer output.Close()

			response, err := http.Get(url)
			if err != nil {
				ErrorLogger.Println(err)
				return
			}
			defer response.Body.Close()
			io.Copy(output, response.Body)
		}
	}
}

func WriteToFile(TextToWrite []byte, fileName string) {
	//fileName := `C:\PriceYUG\price.json` // for other system without path
	//var d = TextToWrite
	err := ioutil.WriteFile(fileName, TextToWrite, 0666)
	if err != nil {
		ErrorLogger.Println(err)
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
	fileName := `C:\PriceYUG\price.json`
	//fileName := `price.json`
	WriteToFile(body, fileName)
}

func getDescr(TokenBraer string) {

	url := goDotEnvVariable("URL_PRICE_DESCR")
	method := "POST"

	payload := strings.NewReader(`{
  		"cats": [` + goDotEnvVariable("CATS") + `],
  		"lang": "RU"
		}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Bearer "+TokenBraer)
	req.Header.Add("Content-Type", "application/json")

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

	fileName := `C:\PriceYUG\price_descr.json`
	//fileName := `price_descr.json`
	WriteToFile(body, fileName)

	var results Data
	json.Unmarshal(body, &results)

	for _, id := range results.Item.Goods {
		//fmt.Println("id :", id.Id)
		catalog := goDotEnvVariable("IMG_PATH") + strconv.Itoa(id.Id)
		for _, pict := range id.Picture {
			//fmt.Println("link :", pict)
			download(pict, catalog)
			//fmt.Println("file does not exist") // это_true
		}

	}

}
