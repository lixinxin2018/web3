package main

func main() {
	bc := NewBlockChain("班长")
	cli := CLI{
		bc: bc,
	}
	cli.Run()
	bc.Db.Close()
}
