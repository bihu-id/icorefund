package refund

import (
	"fmt"
	"ico/chain"
	"ico/tool"
	"ico/tx"
	"log"
	"math/big"
)

var (
	url                = "http://101.132.85.51:8545" // "http://127.0.0.1:8545" // eth node
	abi                = `[{"constant":true,"inputs":[],"name":"sold","outputs":[{"name":"","type":"uint128"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"stop","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"PUBLIC_SALE_PRICE","outputs":[{"name":"","type":"uint128"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"destFoundation","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"owner_","type":"address"}],"name":"setOwner","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"time","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"FUTURE_DISTRIBUTE_LIMIT","outputs":[{"name":"","type":"uint128"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"soldByChannels","outputs":[{"name":"","type":"uint128"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"endTime","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"key","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"total","type":"uint256"}],"name":"canBuy","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"startTime_","type":"uint256"}],"name":"setStartTime","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"userBuys","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"dst","type":"address"},{"name":"wad","type":"uint256"},{"name":"_token","type":"address"}],"name":"transferTokens","outputs":[],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"finalize","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"stopped","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"startTime","outputs":[{"name":"","type":"uint256"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"SELL_HARD_LIMIT","outputs":[{"name":"","type":"uint128"}],"payable":false,"type":"function"},{"constant":false,"inputs":[{"name":"authority_","type":"address"}],"name":"setAuthority","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"SELL_SOFT_LIMIT","outputs":[{"name":"","type":"uint128"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"owner","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"TOTAL_SUPPLY","outputs":[{"name":"","type":"uint128"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"summary","outputs":[{"name":"_sold","type":"uint128"},{"name":"_startTime","type":"uint256"},{"name":"_endTime","type":"uint256"}],"payable":false,"type":"function"},{"constant":false,"inputs":[],"name":"start","outputs":[],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"authority","outputs":[{"name":"","type":"address"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"moreThanSoftLimit","outputs":[{"name":"","type":"bool"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"USER_BUY_LIMIT","outputs":[{"name":"","type":"uint128"}],"payable":false,"type":"function"},{"constant":true,"inputs":[],"name":"MAX_GAS_PRICE","outputs":[{"name":"","type":"uint128"}],"payable":false,"type":"function"},{"inputs":[{"name":"startTime_","type":"uint256"},{"name":"destFoundation_","type":"address"}],"payable":false,"type":"constructor"},{"payable":true,"type":"fallback"},{"anonymous":true,"inputs":[{"indexed":true,"name":"sig","type":"bytes4"},{"indexed":true,"name":"guy","type":"address"},{"indexed":true,"name":"foo","type":"bytes32"},{"indexed":true,"name":"bar","type":"bytes32"},{"indexed":false,"name":"wad","type":"uint256"},{"indexed":false,"name":"fax","type":"bytes"}],"name":"LogNote","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"authority","type":"address"}],"name":"LogSetAuthority","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"name":"owner","type":"address"}],"name":"LogSetOwner","type":"event"}]`
	icoAddress         = "0xd939cdb6b110c96b1758adab5cab1e836ddbdd2b" //ico contract
	tokenAddress       = "0x2770Ad33F5798542DE5DE4bC8Beae3FAdB8De1E5" //token contract
	start        int64 = 4230322                                      //ico start block
	end          int64 = 4233965                                      //4233965                                      //ico end block
	rate         int64 = 200000                                       //token sale rate: 1eth = 200000
	decimal            = 18                                           //token decimal

	refundAddress    = "" //refund address
	refundPrivateKey = "" //refund address private key

)

func ICORefund() {
	holders := GetHolders()
	SendEther(holders)
}

func ICORefund1() {
	holders := make(map[string]*big.Int)

	/*
		holders["0x9a27b56cf576e45ad10d8d3cf26dbe207463e813"] = tool.HexToBigInt("0x17d2a320dd74555000000") //index=0 refund=9.000000000000000000 totalRefund=9.000000000000000000
		holders["0x4e8a164c25a02e30fad798c7257c70b01eae213e"] = tool.HexToBigInt("0x16d07629f842c77600000") //index=1 refund=8.619000000000000000 totalRefund=17.619000000000000000

		/*
	*/

	//ShowHolder(holders)
	SendEther(holders)
}

func GetHolders() map[string]*big.Int {
	holders := make(map[string]*big.Int)
	for index := start; index <= end; index++ {
		senders := GetSenderAtBlock(index, icoAddress)
		fmt.Printf("Block: %d, %d \n", index, len(senders))
		for key := range senders {

			tokenBalance, err := chain.BalanceOf(url, tokenAddress, key)
			if err != nil {
				fmt.Printf("Get TokenBalance err: %s %s\n", key, err)
				break
			}
			value := tool.HexToBigInt(tokenBalance)
			if value.Cmp(new(big.Int).SetInt64(0)) > 0 {
				holders[key] = value
				//fmt.Printf("Address: %s, %s\n", key, val)
			}
		}
	}
	return holders
}

func ShowHolder(holders map[string]*big.Int) {
	totalToken := new(big.Int)
	index := 0
	totalRefund := new(big.Int)
	str1 := ""
	str2 := ""
	for key, value := range holders {
		refundEther := new(big.Int).Div(value, big.NewInt(rate))
		totalRefund = new(big.Int).Add(totalRefund, refundEther)
		//fmt.Printf("Address %d: %s, %s, refund: %s\n", index, key, value, tool.ToBalance(refundEther.String(), 18))
		str1 += fmt.Sprintf("Address %d: %s, %s, refund: %s\n", index, key, value, tool.ToBalance(refundEther.String(), 18))
		str2 += fmt.Sprintf("holders[\"%s\"] = tool.HexToBigInt(\"%s\") //index=%d refund=%s totalRefund=%s\n", key, fmt.Sprintf("0x%x", value), index, tool.ToBalance(refundEther.String(), 18), tool.ToBalance(totalRefund.String(), 18))
		//fmt.Printf("holders[\"%s\"] = tool.HexToBigInt(\"%s\") //index=%d refund=%s totalRefund=%s\n", key, fmt.Sprintf("0x%x", value), index, tool.ToBalance(refundEther.String(), 18), tool.ToBalance(totalRefund.String(), 18))
		totalToken = new(big.Int).Add(totalToken, value)
		index++
	}
	tokenBalance := tool.ToBalance(totalToken.String(), decimal)
	totalEther := new(big.Int).Div(totalToken, big.NewInt(rate))
	fmt.Printf("total: %d, %s, %s, %s, %s\n", len(holders), totalToken, tokenBalance, totalEther, totalRefund)
	fmt.Println(str1)
	fmt.Println(str2)
}

func SendEther(holders map[string]*big.Int) {
	totalToken := new(big.Int)

	totalRefund := new(big.Int)
	successedRefund := new(big.Int)
	failedHolders := make(map[string]*big.Int)

	var index int64

	nonce, err := chain.GetTransactionCount(url, refundAddress)
	if err != nil {
		return
	}
	nonceInt := tool.HexToIntWithoutError(nonce)

	for key, value := range holders {
		fmt.Printf("Address %d: %s,%s\n", index, key, value)

		// -------------------------------------------------------------------------------------------------
		// refundEther = value / decimal / rate * 1e19
		refundEther := new(big.Int).Div(value, big.NewInt(rate)) // wei
		// -------------------------------------------------------------------------------------------------

		totalRefund = new(big.Int).Add(totalRefund, refundEther)
		totalToken = new(big.Int).Add(totalToken, value)
		_, err := send(nonceInt+index, key, refundEther)
		if err != nil {
			//fmt.Printf("Address %d: %t, %s, refund: %s (%s), %s\n", index, false, key, tool.ToBalance(refundEther.String(), 18), refundEther, err)
			failedHolders[key] = value
			break
		} else {
			index++
			//fmt.Printf("Address %d: %t, %s, refund: %s (%s), %s\n", index, true, key, tool.ToBalance(refundEther.String(), 18), refundEther, hash)
		}

	}
	tokenBalance := tool.ToBalance(totalToken.String(), decimal)
	totalEther := new(big.Int).Div(totalToken, big.NewInt(rate))
	fmt.Printf("total: %d/%d, %s (%s), refund: %s, %s/%s\n", len(holders)-len(failedHolders), len(holders), tokenBalance, totalToken, totalEther, successedRefund, totalRefund)
	fmt.Printf("failedHolder: %d \n%s", len(failedHolders), failedHolders)
}

func send(nonce int64, to string, value *big.Int) (string, error) {

	txObj := tx.NewTxObj(tool.IntToHex(nonce), to, fmt.Sprintf("0x%x", value), "")
	txObj.GasLimit = tool.IntToHex(21000)
	txObj.GasPrice = tool.IntToHex(5000000001)

	fmt.Printf("TxObj: %s\n", txObj)
	signedData, err := txObj.SignedData(refundPrivateKey)
	if err != nil {
		log.Printf("SignTx err: %s, nonce %d, value %s", to, nonce, tool.ToBalance(value.String(), 18))
		return "", err
	}
	fmt.Printf("signedData: %s\n", signedData)

	reply, err := chain.SendRawTransaction(url, "0x"+signedData)
	fmt.Printf("txHash: %s,%s\n", reply, err)

	// go func() {
	// 	result, err := self.SendSignedTx(signedData)
	// 	log.Printf("txHash: %s,%s", result, err)
	// }()

	return reply, err
}

func GetSenderAtBlock(blockNum int64, toAddr string) map[string]int64 {

	senders := make(map[string]int64)

	block, err := chain.GetBlock(url, tool.IntToHex(blockNum), true)
	if err != nil || block.Transactions == nil {
		return senders
	}
	//fmt.Printf("Block : %d, %d \n ", blockNum, len(block.Transactions))
	for index := 0; index < len(block.Transactions); index++ {
		tx, ok := block.Transactions[index].(map[string]interface{})
		if !ok {
			return senders
		}
		from := tx["from"].(string)
		to := tx["to"]
		//fmt.Printf("%s, from: %s, to: %s, value: %s \n", block, from, to, value)
		if to != nil && to.(string) == toAddr {
			senders[from] = 1
		}
	}
	return senders
}
