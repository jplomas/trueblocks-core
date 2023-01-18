package rpc

import (
	"context"
	"log"
	"math/big"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// import (
// 	"context"
// 	"log"
// 	"math/big"
// 	"sync"

// 	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
// 	"github.com/ethereum/go-ethereum/ethclient"
// )

// TODO: this looks weird, but before we were creating and deleting the client with every call which
// TODO: overran the number of TPC connection the OS would create (on a Mac). Since then, we now
// TODO: open the connection once and just use it allowing the operating system to clean it up
// TODO: DUPLICATED DUE TO CYCLICAL IMPORT
var perProviderClientMap = map[string]*ethclient.Client{}
var clientMutex sync.Mutex

func GetClient(provider string) *ethclient.Client {
	clientMutex.Lock()
	if perProviderClientMap[provider] == nil {
		// TODO: I don't like the fact that we Dail In every time we want to us this
		// TODO: If we make this a cached item, it needs to be cached per chain, see timestamps
		ec, err := ethclient.Dial(provider)
		if err != nil || ec == nil {
			log.Println("Missdial(" + provider + "):")
			log.Fatalln(err)
		}
		perProviderClientMap[provider] = ec
	}
	clientMutex.Unlock()
	return perProviderClientMap[provider]
}

// TODO: DUPLICATED DUE TO CYCLICAL IMPORT
// TxHashFromNumberAndId returns a transaction's hash if it's a valid transaction
func TxHashFromNumberAndId(chain string, blkNum, txId uint64) (string, error) {
	provider := config.GetRpcProvider(chain)
	ec := GetClient(provider)
	defer ec.Close()

	block, err := ec.BlockByNumber(context.Background(), new(big.Int).SetUint64(blkNum))
	if err != nil {
		return "", err
	}

	tx, err := ec.TransactionInBlock(context.Background(), block.Hash(), uint(txId))
	if err != nil {
		return "", err
	}

	return tx.Hash().Hex(), nil
}

// TODO: DUPLICATED DUE TO CYCLICAL IMPORT
func TxFromNumberAndId(chain string, blkNum, txId uint64) (ethTypes.Transaction, error) {
	provider := config.GetRpcProvider(chain)
	ec := GetClient(provider)
	defer ec.Close()

	block, err := ec.BlockByNumber(context.Background(), new(big.Int).SetUint64(blkNum))
	if err != nil {
		return ethTypes.Transaction{}, err
	}

	tx, err := ec.TransactionInBlock(context.Background(), block.Hash(), uint(txId))
	if err != nil {
		return ethTypes.Transaction{}, err
	}

	return *tx, nil
}
