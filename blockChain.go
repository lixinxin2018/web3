package main

import (
	"log"

	"github.com/boltdb/bolt"
)

// 引入区块链
type BlockChain struct {
	Db   *bolt.DB
	Tail []byte //存储最后一个区块的哈希
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"

func NewBlockChain() *BlockChain {
	genesisBlock := InitBlock()
	var lastHash []byte //最后一个块的哈希
	db, err := bolt.Open(blockChainDb, 0600, nil)
	if err != nil {
		log.Panic("打开数据库失败")
	}
	//增 改
	db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(blockBucket))
		if bk == nil {
			//没有抽屉,需要创建
			bk, err = tx.CreateBucket([]byte(blockBucket))
			if err != nil {
				log.Panic("创建bucket失败")
			}
			bk.Put(genesisBlock.Hash, genesisBlock.Serialize())
			bk.Put([]byte("LastHashKey"), genesisBlock.Hash)
			lastHash = genesisBlock.Hash
		} else {
			lastHash = bk.Get([]byte("lastHashKey"))
		}

		return nil
	})
	defer db.Close()

	return &BlockChain{
		Db:   db,
		Tail: lastHash,
	}

}

// 添加区块
func (BlockChain *BlockChain) AddBlock(data string) {
	// 前区块哈希
	db := BlockChain.Db
	lastHash := BlockChain.Tail
	db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(blockBucket))
		if bk == nil {
			log.Panic("bk不应该为空,请检查")
		}
		block := NewBlock(data, lastHash)
		bk.Put(block.Hash, block.Serialize())
		//记录最后一个区块的哈希值
		bk.Put([]byte("LastHashKey"), block.Hash)
		BlockChain.Tail = block.Hash
		return nil
	})

}
