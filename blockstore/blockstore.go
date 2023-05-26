package blockstore

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

type BlockStore struct {
	Height uint64
}
type TtStore struct {
	Height int64
}

func GetBlockStore() uint64 {
	var blockStore BlockStore
	currentAbPath := getCurrentAbPathByCaller()
	blockStoreFile, _ := filepath.Abs(currentAbPath + "/blockStore.json")
	blockStoreFile1, _ := os.Open(blockStoreFile)
	all, _ := io.ReadAll(blockStoreFile1)
	json.Unmarshal(all, &blockStore)
	return blockStore.Height
}

func SetBlockStore(height uint64) {
	blockStore := BlockStore{
		height,
	}
	currentAbPath := getCurrentAbPathByCaller()
	blockStoreFile, _ := filepath.Abs(currentAbPath + "/blockStore.json")
	blockStoreFile1, _ := os.OpenFile(blockStoreFile, os.O_WRONLY, os.ModeAppend)
	writer := io.Writer(blockStoreFile1)
	blockStoreB, _ := json.Marshal(&blockStore)
	writer.Write(blockStoreB)
}

func GetTtHeight() int64 {
	var ttStore TtStore
	currentAbPath := getCurrentAbPathByCaller()
	blockStoreFile, _ := filepath.Abs(currentAbPath + "/tt.json")
	blockStoreFile1, _ := os.Open(blockStoreFile)
	all, _ := io.ReadAll(blockStoreFile1)
	json.Unmarshal(all, &ttStore)
	return ttStore.Height
}

func SetTtHeight(height int64) {
	ttStore := TtStore{
		height,
	}
	currentAbPath := getCurrentAbPathByCaller()
	blockStoreFile, _ := filepath.Abs(currentAbPath + "/tt.json")
	blockStoreFile1, _ := os.OpenFile(blockStoreFile, os.O_WRONLY, os.ModeAppend)
	writer := io.Writer(blockStoreFile1)
	blockStoreB, _ := json.Marshal(&ttStore)
	writer.Write(blockStoreB)
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
