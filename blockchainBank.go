package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

func main() {
	client, err := ethclient.Dial("https://rpc-mumbai.maticvigil.com")
	if err != nil {
		log.Panic("dail failed")
	}

	block, err := client.BlockByNumber(context.Background(), big.NewInt(35382948))
	if err != nil {
		log.Printf("BlockByNumber error:%v", err)
	}
	for _, tx := range block.Transactions() {
		fmt.Println(tx.Hash(), block.Header().Time, tx.Value(), tx.To(), tx.ChainId(), tx.Gas(), tx.Nonce(), tx.Type())
	}
}
