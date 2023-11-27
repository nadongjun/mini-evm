package main

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/mini-evm/vm"
)

func main() {

	data, err := hex.DecodeString(os.Args[1])
	if err != nil {
		fmt.Printf("Error decoding hex data: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(vm.Run(data))
}
