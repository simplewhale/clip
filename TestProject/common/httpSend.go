package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpParam struct {
	Url    string
	Method string
	Data   []interface{}
}

type HttpReq struct {
	Url     string
	Data    []byte
	IsPost  bool
	Timeout int64
}

var RpcData = `{"jsonrpc":"2.0","id":1,"method":"%s","params":{%s}}`
var RpcStructData = `{"jsonrpc":"2.0","id":1,"method":"%s","params":%s}`

func HttpPost(param *HttpParam) (string, error) {
	curl1 := `{"jsonrpc":"2.0","id":"1","method":"`
	curl2 := `","params":[`
	curl3 := `]}`

	quest := fmt.Sprint(curl1 + param.Method + curl2)
	if len(param.Data) > 0 {
		params := param.Data[0]
		if params == nil {
			quest = fmt.Sprint(quest + curl3)
		} else {
			quest = quest + fmt.Sprint(params)
			if len(param.Data) >= 2 {
				for i := 1; i < len(param.Data); i++ {
					quest = quest + "," + fmt.Sprint(param.Data[i])
				}
			}
		}
	}
	quest = quest + curl3
	//fmt.Println(quest)
	//quest := fmt.Sprintln(curl1+param.method+curl2+"\""+param.data[0]+"\""+","+"\""+param.data[1]+"\""+curl3)
	var jsonStr = []byte(quest)
	req, err := http.NewRequest("POST", param.Url, bytes.NewBuffer(jsonStr))
	if err != nil {
		fmt.Println("req err", err)
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	//SetProxy(client, "http://tunnel.qingtingip.com:8080")
	resp, err := client.Do(req)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("http.post Do error %v\n", err)
		}
	}()
	if resp.StatusCode != 200 {
		for i := 0; i < 5; i++ {
			time.Sleep(time.Duration(i) * time.Second)
			resp, err = client.Do(req)
			if resp.StatusCode == 200 {
				break
			}
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	//fmt.Println("response Body:", string(body))

	return string(body), nil

}

func HttpRequest(q *HttpReq) ([]byte, error) {
	req, err := http.NewRequest("POST", q.Url, bytes.NewReader(q.Data))
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Set("Content-Type", "x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("http.post Do error %v\n", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	fmt.Println("response Body:", string(body))

	return body, nil
}
