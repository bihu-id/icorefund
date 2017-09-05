package chain

import (
	"ico/conf"
	"ico/jsonrpc"
	"ico/tx"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type Block struct {
	Number           string        `json:"number"`
	Hash             string        `json:"hash"`
	ParentHash       string        `json:"parentHash"`
	Nonce            string        `json:"nonce"`
	Sha3Uncles       string        `json:"sha3Uncles"`
	LogsBloom        string        `json:"logsBloom"`
	TransactionsRoot string        `json:"transactionsRoot"`
	StateRoot        string        `json:"stateRoot"`
	Miner            string        `json:"miner"`
	Difficulty       string        `json:"difficulty"`
	TotalDifficulty  string        `json:"totalDifficulty"`
	ExtraData        string        `json:"extraData"`
	Size             string        `json:"size"`
	GasLimit         string        `json:"gasLimit"`
	GasUsed          string        `json:"gasUsed"`
	Timestamp        string        `json:"timestamp"`
	Transactions     []interface{} `json:"transactions"`
	Uncles           []interface{} `json:"uncles"`
}

type Transaction struct {
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	BlockNumber      string `json:"blockNumber"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	GasPrice         string `json:"gasPrice"`
	Gas              string `json:"gas"`
	Input            string `json:"input"`
}

func GetLatestBlockNumber(url string) (string, error) {
	var reply string
	err := jsonrpc.Call(url, "eth_blockNumber", []string{}, &reply)
	return reply, err
}

func GetGasPrice(url string) (result string, err error) {
	var reply string
	er := jsonrpc.Call(url, "eth_gasPrice", []string{}, &reply)
	return reply, er
}

func GetBlock(url, blockNumber string, hasTx bool) (Block, error) {
	//var reply interface{}
	var reply = Block{}
	var params = [2]interface{}{}
	params[0] = blockNumber
	params[1] = hasTx
	err := jsonrpc.Call(url, "eth_getBlockByNumber", params, &reply)
	return reply, err
}

func GetBalance(url, address string) (result string, err error) {
	var reply string
	var params = [2]string{}
	params[0] = address
	params[1] = "latest"
	er := jsonrpc.Call(url, "eth_getBalance", params, &reply)
	return reply, er
}

func GetTransactionCount(url, address string) (result string, err error) {
	var reply string
	var params = [2]string{}
	params[0] = address
	params[1] = "latest"
	er := jsonrpc.Call(url, "eth_getTransactionCount", params, &reply)
	return reply, er
}

func GetTransactionReceipt(url, txHash string) (result interface{}, err error) {
	var reply interface{}
	var params = [1]string{}
	params[0] = txHash
	er := jsonrpc.Call(url, "eth_getTransactionReceipt", params, &reply)
	return reply, er
}

func GetTxReceipt(url, txHash string, reply interface{}) (err error) {
	//reply := types.Receipt{}
	var params = [1]string{}
	params[0] = txHash
	er := jsonrpc.Call(url, "eth_getTransactionReceipt", params, &reply)
	return er
}

func SendRawTransaction(url, txData string) (result string, err error) {
	var reply string
	var params = [1]string{}
	params[0] = txData
	er := jsonrpc.Call(url, "eth_sendRawTransaction", params, &reply)
	return reply, er
}

func Call(url string, param *tx.TxObj, defaultBlock string) (result interface{}, err error) {
	var reply interface{}
	var params = make([]interface{}, 2)
	params[0] = param
	params[1] = defaultBlock
	er := jsonrpc.Call(url, "eth_call", params, &reply)
	return reply, er
}

func BalanceOf(url, tokenAddress, accountAddress string) (balance string, err error) {

	reply, err := CallWithBlock(url, conf.TokenContractABI, tokenAddress, "latest", "balanceOf", common.HexToAddress(accountAddress))
	//log.Printf("SignData: %s, %s", signedData, err)
	if err != nil {
		return "", err
	}
	return reply.(string), nil
}

func CallWithBlock(url, jsonAbi, contractAddress, defaultBlock, functionName string, args ...interface{}) (interface{}, error) {
	abiObj, err := abi.JSON(strings.NewReader(jsonAbi))
	if err != nil {
		return "", err
	}
	dataByte, err := abiObj.Pack(functionName, args...)
	if err != nil {
		return "", err
	}
	//log.Printf("dataByte: %s", common.Bytes2Hex(dataByte))
	txObj := tx.NewTxObj("", contractAddress, "", "0x"+common.Bytes2Hex(dataByte))

	txObj.GasLimit = ""
	txObj.GasPrice = ""
	//log.Printf("TxObj: %s", txObj)
	reply, err := Call(url, txObj, defaultBlock)
	return reply, err
}
