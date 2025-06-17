package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"log"
)

const reward = 12.5

// 1.定义交易结构
type Transaction struct {
	TXID      []byte //交易ID
	TxInputs  []TxInput
	TxOutputs []TxOutput
}

// 定义交易输入
type TxInput struct {
	Txid []byte //交易ID
	//引用的output的索引值
	Index int64
	//解锁脚本,用地址来模拟
	Sig string
}

// 定义交易输出
type TxOutput struct {
	Value         float64
	PublicKeyHash string
}

// 设置交易ID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(tx)
	if err != nil {
		log.Panic("设置交易ID出错,请检查")
	}
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}

// 判断当前的交易是否为挖矿交易
func (tx *Transaction) IsCoinBase() bool {
	input := tx.TxInputs[0]
	//1.交易Input只有一个
	//2.交易id为空
	//3.交易的index为-1
	if len(tx.TxInputs) == 1 && bytes.Equal(input.Txid, []byte{}) && input.Index == -1 {
		return true
	}
	return false
}

// 2.提供创建交易方法

// 2.2将这些UTXO逐一转成inputs
// 2.3创建outputs
// 2.4如果有零钱要找零
func NewTransaction(from, to string, amount float64, bc *BlockChain) *Transaction {
	// 2.1找到最合理的UTXO集合 map[string][]int64
	utxos, resValue := bc.FindNeedUTXOs(from, amount)

	if resValue < amount {
		//无法交易
		fmt.Printf("余额不足,交易失败")
		return nil
	}

	var inputs []TxInput
	var outputs []TxOutput

	for id, indexArray := range utxos {
		for _, index := range indexArray {
			txInput := TxInput{
				Txid:  []byte(id),
				Index: int64(index),
				Sig:   from,
			}
			inputs = append(inputs, txInput)
		}
	}

	txOutput := TxOutput{
		Value:         amount,
		PublicKeyHash: to,
	}
	outputs = append(outputs, txOutput)
	if resValue > amount {
		//找零
		outputs = append(outputs, TxOutput{Value: resValue - amount, PublicKeyHash: from})
	}
	tx := Transaction{
		TXID:      []byte{},
		TxInputs:  inputs,
		TxOutputs: outputs,
	}
	tx.SetHash()
	return &tx
}

// 3.创建挖矿交易
func NewCoinbaseTx(address string, data string) *Transaction {
	//挖矿交易特点,只有一个input,无需引用交易ID,无需引用index,无需签名一般为矿池的名字
	input := TxInput{
		Txid:  []byte{},
		Index: -1,
		Sig:   data,
	}
	output := TxOutput{
		Value:         reward,
		PublicKeyHash: address,
	}
	tx := Transaction{
		TXID:      []byte{},
		TxInputs:  []TxInput{input},
		TxOutputs: []TxOutput{output},
	}
	tx.SetHash()
	return &tx
}

//4.根据交易调整程序
