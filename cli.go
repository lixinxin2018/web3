package main

import (
	"fmt"
	"os"
)

// 接收命令操作区块链
type CLI struct {
	bc *BlockChain
}

const Usage = `
	addBlock --data DATA "add data to blockchain"
	printChain 		"print all blockchain data"
`

func (cli *CLI) Run() {
	//得到所有命令
	args := os.Args
	//分析命令
	if len(args) < 2 {
		fmt.Printf(Usage)
		return
	}
	cmd := args[1]
	switch cmd {
	case "addBlock":
		fmt.Printf("添加区块\n")
		if len(args) == 4 && args[2] == "--data" {
			data := args[3]
			cli.AddBlock(data)
		} else {
			fmt.Printf("添加区块参数使用不当,请检查")
			fmt.Println(Usage)
		}
	case "printChain":
		cli.PrintBlockChain()
	default:
		fmt.Printf("使用了无效的参数命令,请检查")
		fmt.Printf(Usage)
	}
}
