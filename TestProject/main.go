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
	//c = cache.New(cache.NoExpiration,cache.NoExpiration)
	//c.Set("1","1",cache.NoExpiration)
	//go test()
	//time.Sleep(10*time.Second)
	//tokenAddress := "CW2sMRF3JJ7q8rqamJz3iZcdPRNiv3RYKDQ4LfKTkUm7"
	//var url string
	//url = `https://api.solscan.io/token/meta?token=%s`
	//url = fmt.Sprintf(url,tokenAddress)
	//result,err := Get(url)
	//if err != nil{
	//	fmt.Println(err)
	//}
	//decimal := jsoniter.Get(result,"data","decimals").ToInt()
	//fmt.Println("?????????",decimal)
	/*log.Println("????????????...")
	g := golimit.NewGoLimit(2) //max_num(?????????????????????)?????????2
	for i := 0; i < 10; i++ {
		//???????????????1.??? ??????>=max_num, ?????????,?????? ??????<max_num
		g.Add()
		//??????????????????????????????????????????????????????
		//g.SetMax(3)
		go func(g *golimit.GoLimit, i int) {
			defer g.Done() //???????????????1
			time.Sleep(time.Second * 2)
			log.Println(i, "done")
		}(g, i)
	}
	log.Println("????????????")
	g.WaitZero() //??????, ???????????????????????????
	log.Println("????????????")*/
	//pj := jsoniter.Get([]byte(`{"Image":{"URL":"http://example.com/example.gif"}}`), "Image","URL")
	for {
		common.Start()
	}
}
