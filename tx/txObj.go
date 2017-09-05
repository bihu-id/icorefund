package tx

import (
	"encoding/json"
	"ico/tool"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type TxObj struct {
	Nonce       string `json:"nonce,omitempty"`
	From        string `json:"from,omitempty"`
	To          string `json:"to,omitempty"`
	Value       string `json:"value,omitempty"`
	GasLimit    string `json:"gasLimit,omitempty"`
	GasPrice    string `json:"gasPrice,omitempty"`
	Data        string `json:"data,omitempty"`
	Hash        string `json:"hash,omitempty"`
	Blocknumber string `json:"blockNumber,omitempty"`
	BlockTime   string `json:"blockTime,omitempty"`
}

func NewTxObj(nonce, to, value, data string) *TxObj {
	return &TxObj{
		Nonce: nonce,
		To:    to,
		Value: value,
		Data:  data,
	}
}

func NewTxObj2(nonce, to, value, gasLimit, gasPrice, data string) *TxObj {
	return &TxObj{
		Nonce:    nonce,
		To:       to,
		Value:    value,
		GasLimit: gasLimit,
		GasPrice: gasPrice,
		Data:     data,
	}
}

func (o *TxObj) ToJson() []byte {
	bytes, err := json.Marshal(o)
	if err != nil {
		return []byte("{ res:0, resMsg: toJson err }")
	} else {
		return bytes
	}
}

func (tx *TxObj) SignedData(privateKey string) (string, error) {

	signedTx, err := tx.Sign(privateKey)
	if err != nil {
		return "", err
	}
	//log.Printf("signedTx: %s", signedTx)
	txb, err := rlp.EncodeToBytes(signedTx)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	return common.Bytes2Hex(txb), nil
}

func (tx *TxObj) Txhash(privateKey string) (string, error) {

	signedTx, err := tx.Sign(privateKey)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	return signedTx.Hash().Hex(), nil
}

func (tx *TxObj) Sign(privateKey string) (*types.Transaction, error) {
	key, err := crypto.ToECDSA(common.Hex2Bytes(privateKey))
	//log.Printf("txData: %s, %b", tx.Data, common.Hex2Bytes(tool.Strip0x(tx.Data)))
	if err != nil {
		return nil, err
	}
	var tempTx *types.Transaction
	if tx.To == "" {
		tempTx = types.NewContractCreation(
			tool.HexToUintWithoutError(tx.Nonce),
			tool.HexToBigInt(tx.Value),
			tool.HexToBigInt(tx.GasLimit),
			tool.HexToBigInt(tx.GasPrice),
			common.FromHex(tx.Data),
		)
	} else {
		tempTx = types.NewTransaction(
			tool.HexToUintWithoutError(tx.Nonce),
			common.HexToAddress(tx.To),
			tool.HexToBigInt(tx.Value),
			tool.HexToBigInt(tx.GasLimit),
			tool.HexToBigInt(tx.GasPrice),
			common.FromHex(tx.Data),
		)
	}
	//log.Printf("tempTx: %s", tempTx)
	return types.SignTx(tempTx, types.HomesteadSigner{}, key)
	//return tempTx.SignECDSA(types.HomesteadSigner{}, key)
}
