package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
)

// 工作量证明
type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// 创建工作量证明
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		Block: block,
	}
	//指定难度值
	targetStr := "0010000000000000000000000000000000000000000000000000000000000000"

	temInt := big.Int{}
	temInt.SetString(targetStr, 16)

	pow.Target = &temInt
	return &pow
}

// 挖矿
func (ProofOfWork *ProofOfWork) Run() ([]byte, uint64) {
	block := ProofOfWork.Block
	temInt := big.Int{}
	var destHash []byte
	var nonce uint64
	for {
		tmp := [][]byte{
			Uint64ToByte(block.Version),
			block.PrevHash,
			block.MerkeRoot,
			block.Data,
			Uint64ToByte(block.Timestamp),
			Uint64ToByte(ProofOfWork.Target.Uint64()),
			Uint64ToByte(nonce)}
		blockInfo := bytes.Join(tmp, []byte(""))
		srcHash := sha256.Sum256(blockInfo)
		temInt.SetBytes(srcHash[:])
		if temInt.Cmp(ProofOfWork.Target) == -1 {
			fmt.Printf("挖矿成功得到的Hash和Nonce是:%x%d", srcHash[:], nonce)
			return destHash, nonce
		} else {
			nonce++
		}
	}
}
