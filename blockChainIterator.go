package main

import (
	"log"

	"github.com/boltdb/bolt"
)

type BlockChainIterator struct {
	Db                 *bolt.DB
	CurrentHashPointer []byte
}

func (bc *BlockChain) NewIterator() BlockChainIterator {
	return BlockChainIterator{
		Db:                 bc.Db,
		CurrentHashPointer: bc.Tail,
	}
}

func (bci *BlockChainIterator) Next() Block {
	db := bci.Db
	CurrentHashPointer := bci.CurrentHashPointer
	var block Block
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBucket))
		if bucket == nil {
			log.Panic("迭代器的bucket不应该为空,请检查")
		}
		blockTemp := bucket.Get(CurrentHashPointer)
		block = Deserialize(blockTemp)
		bci.CurrentHashPointer = block.PrevHash
		return nil
	})
	return block
}
