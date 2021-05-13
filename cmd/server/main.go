package main

import (
	"fmt"
	"github.com/van-pelt/netcalc/pkg/network/server"
)

func main() {

	s := server.NewCalcServer("tcp", "127.0.0.1:2525")
	err := s.StartServer()
	if err != nil {
		fmt.Println(err)
	}
	return
}
