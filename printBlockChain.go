package main

import "fmt"

func (cli *CLI) AddBlock(data string) {
	cli.bc.AddBlock(data)
	fmt.Printf("添加区块成功\n")
}

func (cli *CLI) PrintBlockChain() {
	it := cli.bc.NewIterator()
	fmt.Println(it.Db, it.CurrentHashPointer)
	for {
		bl := it.Next()
		fmt.Printf("Version:%x\n", bl.Version)
		fmt.Printf("PrevHash:%x\n", bl.PrevHash)
		fmt.Printf("MerkeRoot:%x\n", bl.MerkeRoot)
		fmt.Printf("Hash:%x\n", bl.Hash)
		fmt.Printf("Data:%s\n", bl.Data)
		fmt.Printf("Timestamp:%d\n", bl.Timestamp)
		fmt.Printf("Timestamp:%d\n", bl.Difficulty)
		fmt.Printf("Nonce:%d\n", bl.Nonce)
		if bl.PrevHash == nil {
			fmt.Println("区块链遍历结束")
			break
		}
	}
}
