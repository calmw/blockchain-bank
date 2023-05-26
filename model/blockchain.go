package model

import (
	"github.com/kamva/mgm/v3"
)

type TxInfo struct {
	mgm.IDField `json:",inline" bson:",inline"`
	TxType      int64  `json:"tx_type"`
	ChainID     string `json:"chainId,omitempty" bson:"chainId"`
	BlockHeight string `json:"blockHeight" bson:"blockHeight"`
	Ctime       string `json:"ctime" bson:"ctime"`
	Nonce       uint64 `json:"nonce" bson:"nonce"`
	From        string `json:"from" bson:"from"`
	To          string `json:"to" bson:"to"`
	//Gas         string `json:"gas" bson:"gas"`
	//GasPrice    string `json:"gasPrice" bson:"gasPrice"`
	Amount string `json:"amount" bson:"amount"`
	//Input       *hexutil.Bytes `json:"input" bson:"input"`
	Hash string `json:"hash"`
}
