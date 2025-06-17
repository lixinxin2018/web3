package main

import (
	"fmt"
	"os"
	"strconv"
)

// 接收命令操作区块链
type CLI struct {
	bc *BlockChain
}

const Usage = `
	printChain 		"打印区块数据"
	getBalance --address ADDRESS "获取指定地址的余额"
	send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO,由MINER挖矿,同时写入数据"
`

func (cli *CLI) Run() {
	//得到所有命令
	args := os.Args
	//分析命令
	if len(args) < 2 {
		fmt.Printf("%v:", Usage)
		return
	}
	cmd := args[1]
	switch cmd {
	/* case "addBlock":
	fmt.Printf("添加区块\n")
	if len(args) == 4 && args[2] == "--data" {
		//data := args[3]
		cli.AddBlock([]*Transaction{})
	} else {
		fmt.Printf("添加区块参数使用不当,请检查")
		fmt.Printf("%v:", Usage)
	} */
	case "printChain":
		cli.PrintBlockChain()
	case "getBalance":
		if len(args) == 4 && args[2] == "--address" {
			address := args[3]
			cli.GetBalance(address)
		} else {
			fmt.Printf("添加区块参数使用不当,请检查")
			fmt.Printf("%v:", Usage)
		}
	case "send":
		if len(args) != 7 {
			fmt.Printf("参数使用不当,请检查")
			return
		} else {
			fmt.Printf("转账开始...\n")
			from := args[2]
			to := args[3]
			amount, _ := strconv.ParseFloat(args[4], 64)
			miner := args[5]
			data := args[6]
			cli.Send(from, to, amount, miner, data)
		}
	default:
		fmt.Printf("使用了无效的参数命令,请检查")
		fmt.Printf("%v:", Usage)
	}
}
