package common

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

func Get(url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))
	//client.SetAuth(req)
	//req.Header.Set("Origin", "https://viewblock.io")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	return body, err

}

func getDecimal() {
	tokenAddress := "CW2sMRF3JJ7q8rqamJz3iZcdPRNiv3RYKDQ4LfKTkUm7"
	var url string
	url = `https://api.solscan.io/token/meta?token=%s`
	url = fmt.Sprintf(url, tokenAddress)
	result, err := Get(url)
	if err != nil {
		fmt.Println(err)
	}
	decimal := jsoniter.Get(result, "data", "decimals").ToInt()
	fmt.Println("精度：", decimal)
}
