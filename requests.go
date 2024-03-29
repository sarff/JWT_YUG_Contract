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

//func download(url, catalog string) {
//	//Get the response bytes from the url
//	fileName := catalog + "/" + url[strings.LastIndex(url, "/")+1:]
//	fi, err := os.Stat(fileName)
//	if os.IsNotExist(err) {
//		os.Mkdir(catalog, 0755)
//	}
//	if Exists(fileName) {
//		if err == nil {
//			if fi.Size() < 1000 {
//				os.Remove(fileName)
//			}
//		}
//	} else {
//		//fi, err := os.Stat(fileName)
//		//if err != nil {
//		response, _ := http.Get(url)
//
//		defer response.Body.Close()
//
//		if response.StatusCode == 200 {
//			//Create a empty file
//			file, err := os.Create(fileName)
//			if err != nil {
//				fmt.Println(err)
//			}
//			defer file.Close()
//
//			//Write the bytes to the fiel
//			io.Copy(file, response.Body)
//
//			fi, err := file.Stat()
//			//fmt.Println(fmt.Sprintf("File: %s, size: %s", file, fi.Size()))
//			if err == nil {
//				if fi.Size() < 1000 && err == nil {
//					//file.Close()
//					os.Remove(fileName)
//				}
//			}
//		}
//	}
//}

func download(url, catalog string) {
	fmt.Println(fmt.Sprintf("catalog: %s  url: %s", catalog, url))
	//Get the response bytes from the url
	//fileName := catalog + "/" + url[strings.LastIndex(url, "/")+1:]
	fileNmaeShort := catalog + "/" + replaceString(url[strings.LastIndex(url, "/")+1:])
	fi, err := os.Stat(fileNmaeShort)
	if os.IsNotExist(err) {
		os.Mkdir(catalog, 0755)
	}
	if Exists(fileNmaeShort) {
		if err == nil {
			if fi.Size() < 1000 {
				os.Remove(fileNmaeShort)
			}
		}
	} else {
		//fi, err := os.Stat(fileName)
		//if err != nil {

		response, _ := http.Get(url)

		defer response.Body.Close()

		if response.StatusCode == 200 {
			//Create a empty file
			// разбить файл и реплейснуть

			file, err := os.Create(fileNmaeShort)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()

			//Write the bytes to the fiel
			io.Copy(file, response.Body)

			fi, err := file.Stat()
			//fmt.Println(fmt.Sprintf("File: %s, size: %s", file, fi.Size()))
			if err == nil {
				if fi.Size() < 1000 && err == nil {
					//file.Close()
					os.Remove(fileNmaeShort)
				}
			}
		}
	}
}

func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
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
		catalog := goDotEnvVariable("IMG_PATH")
		replacer := strings.NewReplacer("\\", "", "/", "", ",", "", " ", "", ".", "")
		out := replacer.Replace(strconv.Itoa(id.Id))

		for _, pict := range id.Picture {
			//fmt.Println("link :", pict)
			if pict != "" {
				download(pict, catalog+out)
			}
			//fmt.Println("file does not exist") // это_true
		}

	}

}

func replaceString(str string) string {

	replacer := strings.NewReplacer("\\", "", "/", "", ",", "", " ", "", ".", "")
	//out = replacer.Replace(v.Code)
	lenname := len(strings.Split(str, ".")) - 1
	splitName := strings.SplitN(str, ".", lenname)
	if lenname == 1 {
		return str
	} else {
		out := replacer.Replace(splitName[0])
		return out + splitName[1]
	}

}
