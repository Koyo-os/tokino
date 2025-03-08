package main

import "fmt"

func main() {
	blockchain,err := NewChain()
	if err != nil{
		fmt.Print(err)
	} 

	blockchain.Start()
}