package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

// 定义结构
type Block struct {
	//版本号
	Version uint64
	//前一区块
	PrevHash []byte
	//MerkeRoot
	MerkeRoot []byte
	//当前区块哈希值
	Hash []byte
	//数据
	Data []byte
	//时间戳
	Timestamp uint64
	//难度目标值
	Difficulty uint64
	//随机值
	Nonce uint64
}

// 创建区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := Block{
		Version:    00,
		PrevHash:   prevBlockHash,
		MerkeRoot:  []byte{},
		Hash:       []byte{}, //TODO
		Data:       []byte(data),
		Timestamp:  uint64(time.Now().Unix()),
		Difficulty: 0,
		Nonce:      0,
	}
	//block.setHash()

	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()

	block.Hash = hash
	block.Nonce = nonce
	block.Difficulty = pow.Target.Uint64()
	return &block
}

// 实现 一个辅助函数,功能是将uint转换为byte
func Uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

// 生成哈希
/* func (Block *Block) setHash() {
	tmp := [][]byte{
		Uint64ToByte(Block.Version),
		Block.PrevHash,
		Block.MerkeRoot,
		Block.Data,
		Uint64ToByte(Block.Timestamp),
		Uint64ToByte(Block.Difficulty),
		Uint64ToByte(Block.Nonce)}
	blockInfo := bytes.Join(tmp, []byte(""))
	blockInfo = append(blockInfo, Uint64ToByte(Block.Version)...)
	blockInfo = append(blockInfo, Block.PrevHash...)
	blockInfo = append(blockInfo, Block.MerkeRoot...)
	blockInfo = append(blockInfo, Block.Data...)
	blockInfo = append(blockInfo, Uint64ToByte(Block.Timestamp)...)
	blockInfo = append(blockInfo, Uint64ToByte(Block.Difficulty)...)
	blockInfo = append(blockInfo, Uint64ToByte(Block.Nonce)...)

	hash := sha256.Sum256(blockInfo)
	Block.Hash = hash[:]
}
*/
// 引入区块链
type BlockChain struct {
	Db   *bolt.DB
	Tail []byte //存储最后一个区块的哈希
}

const blockChainDb = "blockChain.db"
const blockBucket = "blockBucket"

func (Block *Block) toByte() []byte {
	return []byte{}
}

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
			bk.Put(genesisBlock.Hash, genesisBlock.toByte())
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
	/* prevHash := BlockChain.Blocks[len(BlockChain.Blocks)-1].Hash
	block := NewBlock(data, prevHash)
	BlockChain.Blocks = append(BlockChain.Blocks, block) */

}

// 创世块
func InitBlock() *Block {
	block := NewBlock("奖励老李50BIC", nil)
	return block
}
