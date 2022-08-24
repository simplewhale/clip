package common

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
	utiljson "github.com/json-iterator/go"
	"github.com/shopspring/decimal"
	"math/big"
	"time"
)

var url = "https://bsc-dataseed1.binance.org/"
var pancakeAddress = "0x10ed43c718714eb63d5aa57b78b54704e256024e"
var bnbAddress = "0xbb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c"
var poolFactory = "0xcA143Ce32Fe78f1f7019d7d551a6402fC5350c73"
var teddypair = "0x1a72de436d1386ea002c5180795a495e50a6d7a2"
var teddyAddress = "0xde301d6a2569aefcfe271b9d98f318baee1d30a4"
var testAddress = "0x20fc7b11b3047db1752a8d21febe4469c816488b"
var priKey = "0x4db4afbef1aec182f1b1f67f2d04fd26cda1c02f1ce3ad60f30ff21359a0e505"
var MethodID = "0x7ff36ab5"
var chainId = int64(56)
var feeBnb = big.NewInt(100000000000000000)

type NonceT struct {
	Nonce string `json:"result"`
	Err   string `json:"error"`
}

func hexToBigInt(hex string) *big.Int {
	if hex == "0" {
		return big.NewInt(0)
	}
	n := new(big.Int)
	n, _ = n.SetString(hex[2:], 16)

	return n
}

func GetNonce(from string) (int64, error) {
	address := fmt.Sprint(`"`, from, `"`)
	state := fmt.Sprint(`"`, "pending", `"`)
	arr := []interface{}{address, state}
	param := HttpParam{
		url,
		"eth_getTransactionCount",
		arr}
	result, err := HttpPost(&param)
	if err != nil {
		return 0, err
	}
	nonce := NonceT{}
	err = json.Unmarshal([]byte(result), &nonce)
	if err != nil {
		log.Error(" err", err)
		return 0, err
	}
	nonceB := hexToBigInt(nonce.Nonce)
	return nonceB.Int64(), nil

}

func StringToPrivateKey(privateKeyStr string) (*ecdsa.PrivateKey, error) {
	privateKeyByte, err := hexutil.Decode(privateKeyStr)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateKeyByte)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func CcmTokenTransfer2(amountIn, amountOutMin *big.Int, path []string, from, to string, deadline *big.Int, method, PriKey string, gaslimit uint64, gasprice *big.Int) (string, error) {
	nonce, err := GetNonce(from)
	if err != nil {
		return "", err
	}

	var data []byte

	data, err = MakeERC20TransferData2(method, amountIn, amountOutMin, path, to, deadline)
	if err != nil {
		return "创建data错误", err
	}
	tx := types.NewTransaction(uint64(nonce), ethcommon.HexToAddress(pancakeAddress), new(big.Int).SetInt64(int64(0)), gaslimit, gasprice, data)

	privateKey, err := StringToPrivateKey(PriKey)
	if err != nil {
		return "私钥初始化错误", err
	}
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainId)), privateKey)
	if err != nil {
		return "签名错误", err
	}

	p, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return "格式化错误", err
	}
	return "0x" + hex.EncodeToString(p), nil
}

func CcmTokenTransfer1(value, amountOutMin *big.Int, path []string, from, to string, deadline *big.Int, method, PriKey string, gaslimit uint64, gasprice *big.Int) (string, error) {
	nonce, err := GetNonce(from)
	if err != nil {
		return "", err
	}

	var data []byte

	data, err = MakeERC20TransferData1(method, amountOutMin, path, to, deadline)
	if err != nil {
		return "创建data错误", err
	}
	tx := types.NewTransaction(uint64(nonce), ethcommon.HexToAddress(pancakeAddress), value, gaslimit, gasprice, data)

	privateKey, err := StringToPrivateKey(PriKey)
	if err != nil {
		return "私钥初始化错误", err
	}
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainId)), privateKey)
	if err != nil {
		return "签名错误", err
	}

	p, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return "格式化错误", err
	}
	return "0x" + hex.EncodeToString(p), nil
}

func IntToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

func MakeERC20TransferData2(method string, amountIn, amountOutMin *big.Int, path []string, to string, deadline *big.Int) ([]byte, error) {
	var data []byte
	methodId, err := hexutil.Decode(method)
	if err != nil {
		return methodId, err
	}
	data = append(data, methodId...)

	padIn := ethcommon.LeftPadBytes(amountIn.Bytes(), 32)
	data = append(data, padIn...)

	padOut := ethcommon.LeftPadBytes(amountOutMin.Bytes(), 32)
	data = append(data, padOut...)

	pad := 5 * 32
	paddPd := ethcommon.LeftPadBytes(IntToBytes(pad), 32)
	data = append(data, paddPd...)

	padTo := ethcommon.LeftPadBytes(ethcommon.HexToAddress(to).Bytes(), 32)
	data = append(data, padTo...)

	padDead := ethcommon.LeftPadBytes(deadline.Bytes(), 32)
	data = append(data, padDead...)

	size := ethcommon.LeftPadBytes(IntToBytes(len(path)), 32)
	data = append(data, size...)
	var adds []byte
	var byadd []byte
	for i := 0; i < len(path); i++ {
		byadd = ethcommon.LeftPadBytes(ethcommon.HexToAddress(path[i]).Bytes(), 32)
		adds = append(adds, byadd...)
	}
	data = append(data, adds...)

	return data, nil
}

func MakeERC20TransferData1(method string, amountOutMin *big.Int, path []string, to string, deadline *big.Int) ([]byte, error) {
	var data []byte
	methodId, err := hexutil.Decode(method)
	if err != nil {
		return methodId, err
	}
	data = append(data, methodId...)

	padOut := ethcommon.LeftPadBytes(amountOutMin.Bytes(), 32)
	data = append(data, padOut...)

	pad := 4 * 32
	paddPd := ethcommon.LeftPadBytes(IntToBytes(pad), 32)
	data = append(data, paddPd...)

	padTo := ethcommon.LeftPadBytes(ethcommon.HexToAddress(to).Bytes(), 32)
	data = append(data, padTo...)

	padDead := ethcommon.LeftPadBytes(deadline.Bytes(), 32)
	data = append(data, padDead...)

	size := ethcommon.LeftPadBytes(IntToBytes(len(path)), 32)
	data = append(data, size...)
	var adds []byte
	var byadd []byte
	for i := 0; i < len(path); i++ {
		byadd = ethcommon.LeftPadBytes(ethcommon.HexToAddress(path[i]).Bytes(), 32)
		adds = append(adds, byadd...)
	}
	data = append(data, adds...)

	return data, nil
}

type TxId struct {
	TxId  string `json:"result"`
	Error `json:"error"`
}

type Error struct {
	Error string `json:"message"`
}

func SendTransaction(txEncode string) (string, error) {
	txEncode = fmt.Sprint(`"`, txEncode, `"`)
	arr := []interface{}{txEncode}
	param := HttpParam{
		url,
		"eth_sendRawTransaction",
		arr}
	result, err := HttpPost(&param)
	if err != nil {
		return "", err
	}
	txId := TxId{}
	err = json.Unmarshal([]byte(result), &txId)
	if err != nil || txId.Error.Error != "" {
		fmt.Println("borcast err", txId.Error.Error)
		return "", err
	}
	return txId.TxId, fmt.Errorf(txId.Error.Error)
}

func swapExactTokensForETH(amountIn, amountOutMin *big.Int, path []string, from, to string, deadline *big.Int, PriKey string, gaslimit uint64, gasprice *big.Int) (string, error) {

	fmt.Println("swap token for eth")
	txEncode, err := CcmTokenTransfer2(amountIn, amountOutMin, path, from, to, deadline, "0x18cbafe5", PriKey, gaslimit, gasprice)
	if err == nil {
		return SendTransaction(txEncode)
	}
	return "", err
}

func swapExactETHForTokens(value, amountOutMin *big.Int, path []string, from, to string, deadline *big.Int, PriKey string, gaslimit uint64, gasprice *big.Int) (string, error) {

	fmt.Println("swap eth for token")
	txEncode, err := CcmTokenTransfer1(value, amountOutMin, path, from, to, deadline, "0x7ff36ab5", PriKey, gaslimit, gasprice)
	if err == nil {
		return SendTransaction(txEncode)
	}
	return "", err
}

//tranInfo.to 为授权代币转账的合约地址 / tranInfo.amount 为授权数量
func approve(from, tokenAddress, prikey string) (string, error) {
	fmt.Println("approve..")
	txEncode, err := CcmTokenTransfer(from, tokenAddress, prikey, "0x095ea7b3")
	if err == nil {
		return SendTransaction(txEncode)
	}
	return "", err
}

func CcmTokenTransfer(from, tokenAddress, priKey, method string) (string, error) {
	nonce, err := GetNonce(from)
	if err != nil {
		return "", err
	}

	var data []byte
	data, err = MakeERC20TransferData(method, pancakeAddress)

	if err != nil {
		return "", err
	}
	tx := types.NewTransaction(uint64(nonce), ethcommon.HexToAddress(tokenAddress), new(big.Int).SetInt64(int64(0)), 60000, big.NewInt(6000000000), data)

	privateKey, err := StringToPrivateKey(priKey)
	if err != nil {
		return "", err
	}
	signTx, err := types.SignTx(tx, types.NewEIP155Signer(big.NewInt(chainId)), privateKey)
	if err != nil {
		return "", err
	}

	p, err := rlp.EncodeToBytes(signTx)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(p), nil
}

func MakeERC20TransferData(method string, toAddress string) ([]byte, error) {
	var data []byte
	methodId, err := hexutil.Decode(method)
	if err != nil {
		return methodId, err
	}
	data = append(data, methodId...)
	paddedAddress := ethcommon.LeftPadBytes(ethcommon.HexToAddress(toAddress).Bytes(), 32)
	data = append(data, paddedAddress...)
	paddedAmount := ethcommon.LeftPadBytes(ethcommon.FromHex("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"), 32)
	data = append(data, paddedAmount...)
	return data, nil
}

func Getpair(tokenAddress string) (string, error) {

	obj := fmt.Sprintf(`{
	  "to": "%v",  
	  "data": "%v" 
	}`, poolFactory, "0xe6a43905000000000000000000000000bb4cdb9cbd36b01bd1cbaebf2de08d9173bc095c000000000000000000000000"+tokenAddress[2:])
	state := fmt.Sprint(`"`, "latest", `"`)
	arr := []interface{}{obj, state}

	param := HttpParam{url, "eth_call", arr}
	result, err := HttpPost(&param)
	if err != nil {
		return "", err
	}
	balance := TxId{}
	err = json.Unmarshal([]byte(result), &balance)
	if err != nil {
		log.Error("total balance result err:", err)
		return "", err
	}
	pair := "0x" + balance.TxId[26:]
	return pair, nil
}

func GetAllowance(address, tokenAddress string) (string, error) {
	data := "0xdd62ed3e000000000000000000000000" + address[2:] + "000000000000000000000000" + pancakeAddress[2:]
	obj := fmt.Sprintf(`{
	  "to": "%v",  
	  "data": "%v" 
	}`, tokenAddress, data)
	state := fmt.Sprint(`"`, "latest", `"`)
	arr := []interface{}{obj, state}

	param := HttpParam{url, "eth_call", arr}
	result, err := HttpPost(&param)
	if err != nil {
		return "", err
	}
	amount := utiljson.Get([]byte(result), "result").ToString()
	if err != nil {
		log.Error("total balance result err:", err)
		return "", err
	}
	return amount, nil
}

type BalanceC struct {
	Balance string `json:"result"`
	Err     string `json:"error"`
}

func GetBalance(tokenAddress, address string) (*big.Int, error) {

	obj := fmt.Sprintf(`{
  "from": "%v",
  "to": "%v",
  "data": "%v"
}`, address, tokenAddress, "0x70a08231000000000000000000000000"+address[2:])
	state := fmt.Sprint(`"`, "latest", `"`)
	arr := []interface{}{obj, state}

	param := HttpParam{url, "eth_call", arr}
	result, err := HttpPost(&param)
	if err != nil {
		return new(big.Int), err
	}
	account := BalanceC{}
	err = json.Unmarshal([]byte(result), &account)
	if err != nil || account.Err != "" {
		return new(big.Int), fmt.Errorf(account.Err)
	}
	balanceT := account.Balance
	if balanceT == "" {
		return new(big.Int), fmt.Errorf("无效地址")
	}
	balanceB := hexToBigInt(balanceT)
	return balanceB, nil
}

func newPendingTxFilter() string {
	arr := []interface{}{}
	param := HttpParam{
		url,
		"eth_newPendingTransactionFilter",
		arr}
	result, err := HttpPost(&param)
	if err != nil {
		return ""
	}

	return utiljson.Get([]byte(result), "result").ToString()
}

func getPendingTxs(filter string) string {
	state := fmt.Sprint(`"`, filter, `"`)
	arr := []interface{}{state}
	param := HttpParam{
		url,
		"eth_getFilterChanges",
		arr}
	result, err := HttpPost(&param)
	if err != nil {
		return ""
	}

	return utiljson.Get([]byte(result), "result").ToString()
}

func getTx(txId string) string {
	txId = fmt.Sprint(`"`, txId, `"`)
	arr := []interface{}{txId}
	param := HttpParam{
		url,
		"eth_getTransactionByHash",
		arr}
	result, err := HttpPost(&param)
	if err != nil {
		return ""
	}

	return result
}

func Start() {

	filter := newPendingTxFilter()
	time.Sleep(time.Second * 5)
	result := getPendingTxs(filter)
	list := utiljson.Get([]byte(result))
	for i := 0; i < list.Size(); i++ {
		txId := list.Get(i).ToString()
		txinfo := getTx(txId)
		tx := utiljson.Get([]byte(txinfo), "result")
		blockNumber := tx.Get("blockNumber").ToString()
		//过滤未打包交易
		if blockNumber != "" {
			continue
		}
		//过滤pancake交易
		if tx.Get("to").ToString() != pancakeAddress {
			continue
		}

		input := tx.Get("input").ToString()
		if input[0:10] != MethodID {
			continue
		}
		value := tx.Get("value").ToString()
		if value == "0x0" {
			continue
		}
		fmt.Println(value)
		tokenAddress := "0x" + input[len(input)-40:]
		fmt.Println(tokenAddress)
		if tokenAddress != teddyAddress {
			continue
		}
		amountMinOut := hexToBigInt(input[10:74])
		targertV := hexToBigInt(value)
		pair, _ := Getpair(tokenAddress)
		tb, _ := GetBalance(tokenAddress, pair)
		bb, _ := GetBalance(bnbAddress, pair)
		td := decimal.NewFromBigInt(tb, 0)
		bd := decimal.NewFromBigInt(bb, 0)
		price := bd.Div(td)
		fmt.Println(price)
		am1 := decimal.NewFromBigInt(amountMinOut, 0)
		am2 := decimal.NewFromBigInt(targertV, 0)
		//计算滑点
		sqr, _ := (am2.Sub(am1.Mul(price))).Div(am2).Float64()
		fmt.Println(sqr)
		//计算买入金额占池比
		sq, _ := am2.Div(bd).Float64()
		if sqr < 0.01 || sq < 0.01 {
			continue
		}
		gas := tx.Get("gas").ToString()
		targetGaslimit := hexToBigInt(gas).Int64()
		gasp := tx.Get("gasPrice").ToString()
		targetGasPrice := hexToBigInt(gasp)
		balance, _ := GetBalance(bnbAddress, testAddress)
		balance = balance.Sub(balance, feeBnb)
		var amountIn *big.Int
		if balance.Cmp(targertV) == 1 || balance.Cmp(targertV) == 0 {
			amountIn = targertV
		} else {
			amountIn = balance
		}
		amm := decimal.NewFromBigInt(amountIn, 0)
		//比目标交易少千一滑点
		sqr1 := decimal.NewFromFloat(1.001 - sqr)
		amountMinOut = amm.Div(price).Mul(sqr1).BigInt()
		deadline := big.NewInt(1648023237)
		//增加1wei
		gasprice := targetGasPrice.Add(targetGasPrice, big.NewInt(1))
		path := []string{bnbAddress, tokenAddress}
		result1, err := swapExactETHForTokens(amountIn, amountMinOut, path, testAddress, testAddress, deadline, priKey, uint64(targetGaslimit), gasprice)
		if err != nil || result1 == "" {
			fmt.Println(err)
			continue
		}
		allowance, _ := GetAllowance(testAddress, tokenAddress)
		if allowance == "0x0000000000000000000000000000000000000000000000000000000000000000" {
			approve(testAddress, tokenAddress, priKey)
		}
		//增加滑点
		amountIn1, _ := GetBalance(tokenAddress, testAddress)
		sqr2, _ := decimal.NewFromString("0.99")
		amm1 := decimal.NewFromBigInt(amountIn, 0)
		amountMinOut1 := amm1.Mul(sqr2).BigInt()
		path1 := []string{tokenAddress, bnbAddress}
		//gasprice1 := targetGasPrice.Sub(targetGasPrice,big.NewInt(1))
		gasprice1 := targetGasPrice
		result2, err := swapExactTokensForETH(amountIn1, amountMinOut1, path1, testAddress, testAddress, deadline, priKey, uint64(targetGaslimit), gasprice1)
		if err != nil || result2 == "" {
			fmt.Println(err)
			continue
		}
	}

}

/*func GetPendingTransactions() string {
	arr := []interface{}{}
	param := HttpParam{
		url,
		"eth_pendingTransactions",
		arr}
	result, err := HttpPost(&param)
	if err != nil{
		return ""
	}

	return utiljson.Get([]byte(result),"result").ToString()
}*/
