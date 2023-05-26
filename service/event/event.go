package event

import (
	eth "github.com/calmw/ethereum"
	"github.com/calmw/ethereum/common"
	"github.com/calmw/ethereum/crypto"
	"math/big"
	"strings"
)

type Sig string

const (
	TransferERC20 Sig = "Transfer(address,address,uint256)" // event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
)

func BuildQuery(contract common.Address, sig Sig, startBlock *big.Int, endBlock *big.Int) eth.FilterQuery {
	query := eth.FilterQuery{
		FromBlock: startBlock,
		ToBlock:   endBlock,
		Addresses: []common.Address{contract},
		Topics: [][]common.Hash{
			{sig.GetTopic()},
		},
	}
	return query
}

func (es Sig) GetTopic() common.Hash {
	return crypto.Keccak256Hash([]byte(es))
}

func GetEventName(sig Sig) string {
	aSig := string(sig)
	index := strings.Index(aSig, "(")
	return aSig[:index]
}
