// Copyright 2021 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.

package rpcClient

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// GetBlockBodyByNumber fetches the block with transactions from the RPC.
func (conn *Connection) GetBlockBodyByNumber(bn uint64) (types.SimpleBlock[types.SimpleTransaction], error) {
	if conn.HasStore() {
		// We only cache blocks with transaction hashes
		cachedBlock := types.SimpleBlock[string]{BlockNumber: bn}
		if err := conn.Store.Read(&cachedBlock, nil); err == nil {
			// Success, we now have to fill in transaction objects
			result := types.SimpleBlock[types.SimpleTransaction]{}
			result.Transactions = make([]types.SimpleTransaction, 0, len(cachedBlock.Transactions))
			success := true
			for index := range cachedBlock.Transactions {
				tx, err := conn.GetTransactionByNumberAndId(cachedBlock.BlockNumber, uint64(index))
				if err != nil {
					success = false
					break
				}
				result.Transactions = append(result.Transactions, *tx)
			}
			if success {
				cachedBlock.Dup(&result)
				return result, nil
			}
		}
	}

	block, rawBlock, err := loadBlock[types.SimpleTransaction](conn, bn, true)
	block.SetRaw(rawBlock) // may have failed, but it's ok
	if err != nil {
		return block, err
	}

	block.Uncles = make([]base.Hash, 0, len(rawBlock.Uncles))
	for _, uncle := range rawBlock.Uncles {
		block.Uncles = append(block.Uncles, base.HexToHash(uncle))
	}

	ts, _ := strconv.ParseInt(rawBlock.Timestamp, 0, 64)
	block.Transactions = make([]types.SimpleTransaction, 0, len(rawBlock.Transactions))
	for _, rawTx := range rawBlock.Transactions {
		// cast transaction to a concrete type
		rawData, ok := rawTx.(map[string]any)
		if !ok {
			err = errors.New("cannot cast raw block transaction into map")
			return block, err
		}
		raw := types.NewRawTransactionFromMap(rawData)

		// Get the receipt
		var receipt types.SimpleReceipt
		receipt, err = conn.GetReceipt(ReceiptQuery{
			Bn:       uint64(bn),
			Txid:     uint64(raw.TxIndex()),
			TxHash:   raw.TxHash(),
			GasPrice: raw.TxGasPrice(),
			NeedsTs:  true,
			Ts:       ts,
		})
		if err != nil {
			return block, err
		}

		tx := types.NewSimpleTransaction(raw, &receipt, ts)
		block.Transactions = append(block.Transactions, *tx)

		if conn.HasStore() && conn.enabledMap["txs"] && conn.IsFinal(tx.Timestamp) {
			_ = conn.Store.Write(tx, nil)
		}
	}

	if conn.HasStore() && conn.IsFinal(block.Timestamp) {
		_ = conn.Store.Write(&block, nil)
	}

	return block, nil
}

// GetBlockHeaderByNumber fetches the block with only transactions' hashes from the RPC
func (conn *Connection) GetBlockHeaderByNumber(bn uint64) (block types.SimpleBlock[string], err error) {
	if conn.HasStore() {
		block.BlockNumber = bn
		if err := conn.Store.Read(&block, nil); err == nil {
			return block, nil
		}
	}

	block, rawBlock, err := loadBlock[string](conn, bn, false)
	block.SetRaw(rawBlock) // may have failed, but it's ok
	if err != nil {
		return block, err
	}

	block.Uncles = make([]base.Hash, 0, len(rawBlock.Uncles))
	for _, uncle := range rawBlock.Uncles {
		block.Uncles = append(block.Uncles, base.HexToHash(uncle))
	}

	block.Transactions = make([]string, 0, len(rawBlock.Transactions))
	for _, txHash := range rawBlock.Transactions {
		block.Transactions = append(block.Transactions, fmt.Sprint(txHash))
	}

	if conn.HasStoreWritable() && conn.IsFinal(block.Timestamp) {
		_ = conn.Store.Write(&block, nil)
	}

	return block, nil
}

// GetBlockTimestamp returns the timestamp associated with a given block
func (conn *Connection) GetBlockTimestamp(bn *uint64) base.Timestamp {
	if ec, err := conn.getClient(); err != nil {
		logger.Error("Could not connect to RPC client", err)
		return 0
	} else {
		defer ec.Close()

		var blockNumber *big.Int
		if bn != nil {
			blockNumber = big.NewInt(int64(*bn))
		}

		r, err := ec.HeaderByNumber(context.Background(), blockNumber)
		if err != nil {
			logger.Error("Could not connect to RPC client", err)
			return 0
		}

		ts := base.Timestamp(r.Time)
		if ts == 0 {
			// The RPC does not return a timestamp for block zero, so we simulate it with ts from block one less 13 seconds
			// TODO: Chain specific
			return conn.GetBlockTimestamp(utils.PointerOf(uint64(1))) - 13
		}

		return ts
	}
}

// GetBlockHashByHash returns a block's hash if it's a valid block
func (conn *Connection) GetBlockHashByHash(hash string) (string, error) {
	if ec, err := conn.getClient(); err != nil {
		return "", err
	} else {
		defer ec.Close()

		block, err := ec.BlockByHash(context.Background(), common.HexToHash(hash))
		if err != nil {
			return "", err
		}

		return block.Hash().Hex(), nil
	}
}

// GetBlockNumberByHash returns a block's hash if it's a valid block
func (conn *Connection) GetBlockNumberByHash(hash string) (base.Blknum, error) {
	if ec, err := conn.getClient(); err != nil {
		return 0, err
	} else {
		defer ec.Close()

		block, err := ec.BlockByHash(context.Background(), common.HexToHash(hash))
		if err != nil {
			return 0, err
		}

		return block.NumberU64(), nil
	}
}

// GetBlockHashByNumber returns a block's hash if it's a valid block
func (conn *Connection) GetBlockHashByNumber(bn uint64) (string, error) {
	if ec, err := conn.getClient(); err != nil {
		return "", err
	} else {
		defer ec.Close()

		block, err := ec.BlockByNumber(context.Background(), new(big.Int).SetUint64(bn))
		if err != nil {
			return "", err
		}

		return block.Hash().Hex(), nil
	}
}

func (conn *Connection) getBlockReward(bn uint64) *big.Int {
	if bn == 0 {
		return big.NewInt(0)
	} else if bn < byzantiumBlock {
		return big.NewInt(5000000000000000000)
	} else if bn < constantinopleBlock {
		return big.NewInt(3000000000000000000)
	} else {
		return big.NewInt(2000000000000000000)
	}
}

func (conn *Connection) getBlockRaw(bn uint64, withTxs bool) (*types.RawBlock, error) {
	method := "eth_getBlockByNumber"
	params := rpc.Params{fmt.Sprintf("0x%x", bn), withTxs}

	if block, err := rpc.Query[types.RawBlock](conn.Chain, method, params); err != nil {
		return &types.RawBlock{}, err
	} else {
		if bn == 0 {
			// The RPC does not return a timestamp for the zero block, so we make one
			block.Timestamp = fmt.Sprintf("0x%x", conn.GetBlockTimestamp(utils.PointerOf(uint64(0))))
		} else if utils.MustParseUint(block.Timestamp) == 0 {
			return &types.RawBlock{}, fmt.Errorf("block at %s returned an error: %w", fmt.Sprintf("%d", bn), ethereum.NotFound)
		}

		return block, nil
	}
}

// loadBlock fetches block from RPC, but it does not try to fill Transactions field. This is delegated to
// more specialized functions and makes loadBlock generic.
func loadBlock[Tx string | types.SimpleTransaction](conn *Connection, bn uint64, withTxs bool) (block types.SimpleBlock[Tx], rawBlock *types.RawBlock, err error) {
	rawBlock, err = conn.getBlockRaw(bn, withTxs)
	if err != nil {
		return
	}

	ts, err := hexutil.DecodeUint64(rawBlock.Timestamp)
	if err != nil {
		return
	}

	blockNumber, err := hexutil.DecodeUint64(rawBlock.BlockNumber)
	if err != nil {
		return
	}

	gasLimit, err := hexutil.DecodeUint64(rawBlock.GasLimit)
	if err != nil {
		return
	}

	gasUsed, err := hexutil.DecodeUint64(rawBlock.GasUsed)
	if err != nil {
		return
	}

	difficulty, err := hexutil.DecodeUint64(rawBlock.Difficulty)
	if err != nil {
		return
	}

	uncles := make([]base.Hash, 0, len(rawBlock.Uncles))
	for _, uncle := range rawBlock.Uncles {
		uncles = append(uncles, base.HexToHash(uncle))
	}

	block = types.SimpleBlock[Tx]{
		BlockNumber: blockNumber,
		Timestamp:   base.Timestamp(ts), // note that we turn Ethereum's timestamps into types.Timestamp upon read.
		Hash:        base.HexToHash(rawBlock.Hash),
		ParentHash:  base.HexToHash(rawBlock.ParentHash),
		GasLimit:    gasLimit,
		GasUsed:     gasUsed,
		Miner:       base.HexToAddress(rawBlock.Miner),
		Difficulty:  difficulty,
		Uncles:      uncles,
	}
	return
}
