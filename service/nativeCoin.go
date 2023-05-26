package service

import (
	"context"
	"fmt"
	"github.com/calmw/ethereum/core/types"
	"github.com/calmw/ethereum/ethclient"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math/big"
	"time"
	"tx-record/blockstore"
	"tx-record/log"
	"tx-record/model"
)

var goroutinePool = make(chan struct{}, 5)

type NativeCoin struct {
	Client        *ethclient.Client
	SleepDuration time.Duration
	ChainId       string
}

var done chan interface{}

func NewNativeCoin(client *ethclient.Client) *NativeCoin {
	return &NativeCoin{
		Client:        client,
		SleepDuration: time.Second * 5,
		ChainId:       "137",
	}
}

func (s NativeCoin) NativeCoinTx() {
	s.PollBlock(done, s.Client, s.BlockNumberStreaming(done, s.Client))
}

func (s NativeCoin) BlockNumberStreaming(done <-chan interface{}, client *ethclient.Client) <-chan uint64 {
	blockNumberOutStream := make(chan uint64)
	go func() {
		for {
			select {
			case <-done:
				log.Logger.Sugar().Error("PollBlock done")
			default:
				number, err := client.BlockNumber(context.Background())
				if err != nil {
					log.Logger.Sugar().Errorf("client.BlockNumber error:%v", err)
					time.Sleep(s.SleepDuration)
					continue
				}
				blockStoreHeight := blockstore.GetBlockStore()
				if blockStoreHeight < number {
					for i := blockStoreHeight + 1; i <= number; i++ {
						start := time.Now()
						blockNumberOutStream <- i
						fmt.Println("处理区块", i, "耗时：", time.Since(start))
					}
				} else {
					time.Sleep(s.SleepDuration)
				}
			}
		}
	}()
	return blockNumberOutStream
}

func (s NativeCoin) PollBlock(done <-chan interface{}, client *ethclient.Client, blockNumberStreaming <-chan uint64) {
	for {
		select {
		case num, ok := <-blockNumberStreaming:
			if !ok {
				log.Logger.Sugar().Error("PollBlock.blockNumberStreaming end")
			}
			goroutinePool <- struct{}{}
			go s.GetAndSaveTx(client, num)
			blockstore.SetBlockStore(num)

		case <-done:
			log.Logger.Sugar().Error("PollBlock done")
		}
	}
}

func (s NativeCoin) GetAndSaveTx(client *ethclient.Client, height uint64) {
	var err error
	var block *types.Block
	retry := 2

	start := time.Now()
	for i := 0; i <= retry; i++ {
		if i == retry {
			return
		}
		block, err = client.BlockByNumber(context.Background(), big.NewInt(int64(height)))
		if err != nil {
			log.Logger.Sugar().Errorf("BlockByNumber error:%v", err)
			time.Sleep(time.Millisecond * 500)
			continue
		}
		break
	}
	fmt.Println("BlockByNumber耗时：", time.Since(start))

	tm := time.Unix(int64(block.Time()), 0)
	number := block.Number().String()
	blockCtime := tm.Format("2006-01-02 15:04:05")

	for _, tx := range block.Transactions() {
		if tx.Value().Cmp(big.NewInt(0)) <= 0 || tx.To() == nil { // to is nil, 0xbcdff8fab41e214dd0461eafa61de134d95cb0d3211702482b2a5fe88275ce83
			continue
		}
		blockTx, _ := tx.TxInfo()
		txInfo := model.TxInfo{
			TxType:      1,
			ChainID:     s.ChainId,
			BlockHeight: number,
			Ctime:       blockCtime,
			//Nonce:       tx.Nonce(),
			From:   blockTx.From.String(),
			To:     tx.To().String(),
			Amount: tx.Value().String(),
			Hash:   tx.Hash().String(),
		}
		SaveTxInfoToDB(txInfo)

		<-goroutinePool
	}
}

func SaveTxInfoToDB(txInfo model.TxInfo) {
	start := time.Now()
	upsert := true
	blockTxUpdateOptions := options.UpdateOptions{Upsert: &upsert}
	_, err := mgm.Coll(&model.TxInfo{}).UpdateOne(context.Background(), bson.M{
		"hash":   txInfo.Hash,
		"from":   txInfo.From,
		"to":     txInfo.To,
		"amount": txInfo.Amount,
	}, bson.M{"$set": bson.M{
		"chainId":     txInfo.ChainID,
		"ctime":       txInfo.Ctime,
		"tx_type":     txInfo.TxType,
		"blockHeight": txInfo.BlockHeight,
	}}, &blockTxUpdateOptions)
	fmt.Println("插入数据耗时：", time.Since(start))
	if err != nil {
		log.Logger.Sugar().Errorf("SaveTxInfoToDB error:%v", err)
	}
}
