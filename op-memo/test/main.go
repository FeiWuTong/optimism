package main

import (
	"context"
	"flag"
	"fmt"

	memo "github.com/ethereum-optimism/optimism/op-memo"
	"github.com/ethereum/go-ethereum/common"
)

var (
	defaultData = "0x095ea7b30000000000000000000000000776b0a89b77163a33bb5f062f3a781f920adcf100000000000000000000000000000000000000000000000000005af3107a4000"
	rpcaddr     = flag.String("http", "localhost:15678", "middleware http rpc")
	testdata    = flag.String("data", defaultData, "test data for submission and retrieval")
)

func main() {
	flag.Parse()
	dac, err := memo.NewDAClient(*rpcaddr)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := [][]byte{common.FromHex(*testdata)}
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
