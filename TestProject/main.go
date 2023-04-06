package main

import (
	"TestProject/common"
	"bytes"
	"fmt"
	//jsoniter "github.com/json-iterator/go"
	"github.com/patrickmn/go-cache"
	"io/ioutil"
	"math/big"
	"net/http"
	"regexp"
	"strings"
)

var c *cache.Cache

func test() {
	fmt.Println(c.Get("1"))
}

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

var JsonRpcData = `{"jsonrpc":"2.0","id":"1","method":"%s","params":[%s]}`
var JsonRpcStructData = `{"jsonrpc":"2.0","id":"1","method":"%s","params":%s}`

func Post(url, method string, RPCParams ...string) ([]byte, error) {
	for i, str := range RPCParams {
		if str == "true" || str == "null" {
			continue
		}
		if quote, err := regexp.MatchString("^[0-9]+$", str); err == nil && quote {
			continue
		}
		if quote, err := regexp.MatchString("^[a-zA-Z0-9-*]+$", str); err == nil && quote {
			RPCParams[i] = fmt.Sprintf(`"%s"`, RPCParams[i])
		}
	}
	params := strings.Join(RPCParams, `,`)
	postString := fmt.Sprintf(JsonRpcData, method, params)
	if len(RPCParams) != 0 {
		if argsStart, err := regexp.MatchString(`^{"args":`, RPCParams[0]); err == nil && argsStart {
			postString = fmt.Sprintf(JsonRpcStructData, method, params)
		}
	}
	//fmt.Println(postString)
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(postString))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	return body, err
}

func NormalPost(url, postString string) ([]byte, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(postString))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return []byte(""), err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	return body, err
}

func hexToBigInt(hex string) *big.Int {
	if hex == "0" {
		return big.NewInt(0)
	}
	n := new(big.Int)
	n, _ = n.SetString(hex[2:], 16)

	return n
}

func main() {
	
	for {
		common.Start()
	}
}
