package main

import (
	"fmt"
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

func NewBlockChain(address string) *BlockChain {
	genesisBlock := InitBlock1(address)
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
			lastHash = bk.Get([]byte("LastHashKey"))
		}

		return nil
	})

	return &BlockChain{
		Db:   db,
		Tail: lastHash,
	}

}

// 创世块
func InitBlock1(address string) *Block {
	coinbase := NewCoinbaseTx(address, "Go一起创世块")
	return NewBlock([]*Transaction{coinbase}, nil)
}

// 添加区块
func (BlockChain *BlockChain) AddBlock(txs []*Transaction) {
	// 前区块哈希
	db := BlockChain.Db
	lastHash := BlockChain.Tail
	db.Update(func(tx *bolt.Tx) error {
		bk := tx.Bucket([]byte(blockBucket))
		if bk == nil {
			log.Panic("bk不应该为空,请检查")
		}
		block := NewBlock(txs, lastHash)
		block.PrevHash = lastHash
		bk.Put(block.Hash, block.Serialize())
		//记录最后一个区块的哈希值
		bk.Put([]byte("LastHashKey"), block.Hash)
		BlockChain.Tail = block.Hash
		return nil
	})
}

// 查找余额
func (bc *BlockChain) FindUTXOs(address string) []TxOutput {
	var UTXOS []TxOutput
	txs := bc.FindUTXOTransaction(address)
	for _, tx := range txs {
		for _, output := range tx.TxOutputs {
			if address == output.PublicKeyHash {
				UTXOS = append(UTXOS, output)
			}
		}
	}
	return UTXOS
}

// 获取需要的UTXO
func (bc *BlockChain) FindNeedUTXOs(from string, amount float64) (map[string][]int64, float64) {
	utxos := make(map[string][]int64)
	var calc float64
	txs := bc.FindUTXOTransaction(from)
	for _, tx := range txs {
		for i, output := range tx.TxOutputs {
			//总额小于转账金额
			if calc < amount {
				//output中找到自己的UTXO
				if output.PublicKeyHash == from {
					//统计当前UTXO的总额
					calc += output.Value
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], int64(i))
				}
				//加完满足条件
				if calc >= amount {
					//fmt.Printf("找到了满足条件的金额:%f\n", calc)
					return utxos, calc
				}
			}
		}
	}
	if calc >= amount {
		fmt.Printf("金额不足无法交易:%f\n", calc)
		return nil, calc
	}
	return utxos, calc
}

func (bc *BlockChain) FindUTXOTransaction(address string) []*Transaction {
	var trs []*Transaction
	//定义一个map保存消费过的output,key是这个output的交易ID,value是这个交易中索引的数组
	//map[txid][]int64
	//spentOutputs := make(map[string][]int64)
	spentOutputs := make(map[string][]int64)

	it := bc.NewIterator()
	//遍历区块
	for {
		block := it.Next()
		txs := block.Transaction
		//遍历交易
		for _, tx := range txs {
			//fmt.Printf("current txid:%x\n", tx.TXID)
			//遍历output,找到自己的UTXO
		OUTPUT:
			for i, output := range tx.TxOutputs {
				//遍历消费过UTXO
				if spentOutputs[string(tx.TXID)] != nil {
					for _, j := range spentOutputs[string(tx.TXID)] {
						//确认消费过的UTXO
						if int64(i) == j {
							continue OUTPUT
						}
					}
				}
				//output中找到自己的UTXO
				if output.PublicKeyHash == address {
					trs = append(trs, tx)
					//fmt.Printf("输出中找到的余额是:%f\n", UTXOS[0].Value)
				}
			}
			//如果当前交易是挖矿交易,那么不做遍历,直接跳过
			if !tx.IsCoinBase() {
				//遍历input,找到自己消费过的UTXO的集合
				for _, input := range tx.TxInputs {
					if input.Sig == address {
						spentOutputs[string(input.Txid)] = append(spentOutputs[string(input.Txid)], input.Index)
					}
				}
			} else {
				//fmt.Printf("该交易是挖矿交易:%v\n", tx.TXID)
			}
		}
		if block.PrevHash == nil {
			break
		}
	}
	return trs
}
