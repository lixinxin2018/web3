package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"log"
	"time"
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
	Transaction []*Transaction
	//时间戳
	Timestamp uint64
	//难度目标值
	Difficulty uint64
	//随机值
	Nonce uint64
}

// 创建区块
func NewBlock(txs []*Transaction, prevBlockHash []byte) *Block {
	block := Block{
		Version:     00,
		PrevHash:    prevBlockHash,
		MerkeRoot:   []byte{},
		Hash:        []byte{}, //TODO
		Transaction: txs,
		Timestamp:   uint64(time.Now().Unix()),
		Difficulty:  0,
		Nonce:       0,
	}
	//block.setHash()

	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run()

	block.Hash = hash
	block.Nonce = nonce
	block.Difficulty = pow.Target.Uint64()
	block.MerkeRoot = block.MakeMerkeRoot()
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

func (Block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(&Block)
	if err != nil {
		log.Panic("区块序列化失败")
	}
	return buffer.Bytes()
}

func Deserialize(data []byte) Block {
	des := gob.NewDecoder(bytes.NewReader(data))
	var block Block
	err := des.Decode(&block)
	if err != nil {
		log.Panic("区块反序列化失败")
	}
	return block
}

// 生成梅克儿根
func (block *Block) MakeMerkeRoot() []byte {
	var info []byte
	for _, tx := range block.Transaction {
		info = append(info, tx.TXID...)
	}
	hash := sha256.Sum256(info)
	return hash[:]
}
