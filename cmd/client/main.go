package main

import (
	"fmt"
	"github.com/van-pelt/netcalc/pkg/client"
)

func main() {
	cli := client.NewClient("tcp", "127.0.0.1:2525")
	err := cli.StartCalcTerminal()
	if err != nil {
		fmt.Println(err)
		cli.PrintHelp()
	}
}
