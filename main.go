package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
)

type Response struct {
	Jsonrpc string `"json: jsonrpc"`
	Id      int    `"id"`
	Result  string `"result"`
}

func main() {
	dir := `keystore`

	var s []string
	files, _ := GetAllFile(dir, s)

	for _, f := range files {
		address := "0x" + f[37:]
		requstInfura(address)
		//time.Sleep(1)
	}
}

func requstInfura(address string) {

	jsonString := `{"jsonrpc":"2.0","method":"eth_getBalance","params":["` + address + `", "latest"],"id":1}`
	url := "https://mainnet.infura.io/v3/{API_KEY}"
	rsps, err := http.Post(url, "application/json", bytes.NewReader([]byte(jsonString)))
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	defer rsps.Body.Close()

	body, err := ioutil.ReadAll(rsps.Body)
	if err != nil {
		fmt.Println("Read body failed:", err)
		return
	}
	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Request failed:", err)
		return
	}
	if response.Result == "0x0" || response.Result == "0x2540be400" {
		return
	}
	balance, err := strconv.ParseInt(response.Result, 16, 64)
	if err != nil {
		fmt.Println(err)
	}
	lastBalance := decimal.NewFromInt(balance).Div(decimal.NewFromInt(1000000000000000000)).Truncate(2).InexactFloat64()
	fmt.Printf("Zero Address is : %s, Balance: %f", address, lastBalance)

	if lastBalance > 0.01 {
		fmt.Printf("#### Address is : %s, Balance: %f", address, lastBalance)
	}
}

func GetAllFile(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}

	for _, fi := range rd {
		if !fi.IsDir() {
			//fullName := pathname + "/" + fi.Name()
			fullName := fi.Name()

			s = append(s, fullName)
		}
	}
	return s, nil
}
