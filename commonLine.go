package main

import (
	"fmt"
	"time"
)

/* func (cli *CLI) AddBlock(txs []*Transaction) {
	cli.bc.AddBlock(txs)
	fmt.Printf("添加区块成功\n")
} */

func (cli *CLI) PrintBlockChain() {
	it := cli.bc.NewIterator()
	fmt.Println(it.Db, it.CurrentHashPointer)
	for {
		bl := it.Next()
		fmt.Printf("Version:%x\n", bl.Version)
		fmt.Printf("PrevHash:%x\n", bl.PrevHash)
		fmt.Printf("MerkeRoot:%x\n", bl.MakeMerkeRoot())
		fmt.Printf("Hash:%x\n", bl.Hash)
		fmt.Printf("Data:%s\n", bl.Transaction[0].TxInputs[0].Sig)
		timeFormat := time.Unix(int64(bl.Timestamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("Timestamp:%s\n", timeFormat)
		fmt.Printf("Difficulty:%d\n", bl.Difficulty)
		fmt.Printf("Nonce:%d\n", bl.Nonce)
		if bl.PrevHash == nil {
			fmt.Println("区块链遍历结束")
			break
		}
	}
}

func (cli *CLI) GetBalance(address string) {
	txOutputs := cli.bc.FindUTXOs(address)
	total := 0.0

	for _, v := range txOutputs {
		total += v.Value
	}
	fmt.Printf("%s余额为:%f\n", address, total)
}

func (cli *CLI) Send(from, to string, amount float64, miner, data string) {
	fmt.Printf("from:%s\n", from)
	fmt.Printf("to:%s\n", to)
	fmt.Printf("amount:%f\n", amount)
	fmt.Printf("miner:%s\n", miner)
	fmt.Printf("data:%s\n", data)
	//1.创建挖矿交易
	coinbase := NewCoinbaseTx(miner, data)
	//2.创建一个普通交易
	tx := NewTransaction(from, to, amount, cli.bc)
	if tx == nil {
		return
	}
	//3.添加到区块
	cli.bc.AddBlock([]*Transaction{coinbase, tx})
	fmt.Println("转账成功")
}
