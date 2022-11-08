package cache

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/config"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

// Any data structure that we know how to cache
type cacheable interface {
	types.SimpleBlock |
		types.SimpleTransaction
}

// getCacheAndChainPath returns path to cache for given chain
func getCacheAndChainPath(chainName string) string {
	cacheDir := config.GetRootConfig().Settings.CachePath
	return path.Join(cacheDir, chainName)
}

// save writes contents of `content` Reader to a file
func save(chainName string, filePath string, content io.Reader) (err error) {
	cacheDir := getCacheAndChainPath(chainName)
	fullPath := path.Join(cacheDir, filePath)
	file, err := os.Create(fullPath)
	if err != nil {
		return
	}

	_, err = io.Copy(file, content)
	return
}

// load opens binary cache file for reading
func load(chainName string, filePath string) (file io.Reader, err error) {
	cacheDir := getCacheAndChainPath(chainName)
	fullPath := path.Join(cacheDir, filePath)
	file, err = os.Open(fullPath)
	return
}

func remove(chainName string, filePath string) (err error) {
	cacheDir := getCacheAndChainPath(chainName)
	fullPath := path.Join(cacheDir, filePath)
	return os.Remove(fullPath)
}

// setItem serializes value into binary format and saves it to a file
func setItem[Data cacheable](
	chainName string,
	filePath string,
	value *Data,
	write func(w *bufio.Writer, d *Data) error,
) (err error) {
	buf := bytes.Buffer{}
	writer := bufio.NewWriter(&buf)
	err = write(writer, value)
	if err != nil {
		return
	}
	reader := bytes.NewReader(buf.Bytes())
	err = save(chainName, filePath, reader)
	return
}

// getItem reads data structure from binary format
func getItem[Data cacheable](
	chainName string,
	filePath string,
	read func(w *bufio.Reader) (*Data, error),
) (value *Data, err error) {
	reader, err := load(chainName, filePath)
	if err != nil {
		return
	}

	bufReader := bufio.NewReader(reader)
	value, err = read(bufReader)
	if err != nil && !os.IsNotExist(err) {
		// Ignore the error, we will re-try next time
		remove(chainName, filePath)
	}
	return
}

// SetBlock stores block in the cache
func SetBlock(chainName string, block *types.SimpleBlock) (err error) {
	filePath := getPathByBlock(ItemBlock, block.BlockNumber)

	return setItem(
		chainName,
		filePath,
		block,
		WriteBlock,
	)
}

// GetBlock reads block from the cache
func GetBlock(chainName string) (block *types.SimpleBlock, err error) {
	filePath := getPathByBlock(ItemBlock, block.BlockNumber)

	return getItem(
		chainName,
		filePath,
		ReadBlock,
	)
}

// SetTransaction stores transaction in the cache
func SetTransaction(chainName string, tx *types.SimpleTransaction) (err error) {
	filePath := getPathByBlock(ItemTransaction, tx.BlockNumber)

	return setItem(
		chainName,
		filePath,
		tx,
		WriteTransaction,
	)
}

// GetTransaction reads transactiono from the cache
func GetTransaction(chainName string, blockNumber types.Blknum, txIndex uint64) (tx *types.SimpleTransaction, err error) {
	filePath := getPathByBlockAndTransactionIndex(ItemTransaction, blockNumber, txIndex)

	return getItem(
		chainName,
		filePath,
		ReadTransaction,
	)
}
