package main

import (
	"context"
	"flag"
	"fmt"

	memo "github.com/ethereum-optimism/optimism/op-memo"
)

var (
	rpcaddr = flag.String("http", "localhost:15678", "middleware http rpc")
)

func main() {
	flag.Parse()
	dac, err := memo.NewDAClient(*rpcaddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := [][]byte{}
	ids, _, err := dac.Client.Submit(context.Background(), data, -1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ids[0])

	gotData, err := dac.Client.Get(context.Background(), ids)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(gotData[0])
}
