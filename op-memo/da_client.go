package memo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum/common"
)

var (
	getRoute  = "/da/getObject"
	putRoute  = "/da/putObject"
	initRoute = "/da/warmup"
)

// middleware client
type DAClient struct {
	Client *MemoMiddlewareClient
}

func NewDAClient(rpc string) (*DAClient, error) {
	client := NewMemoMiddlewareClient(rpc)
	err := client.Start()
	if err != nil {
		return nil, err
	}
	return &DAClient{
		Client: client,
	}, nil
}

//
// interface
//

type MemoMiddlewareClient struct {
	rpcaddr string
}

func NewMemoMiddlewareClient(rpcaddr string) *MemoMiddlewareClient {
	return &MemoMiddlewareClient{
		rpcaddr: rpcaddr,
	}
}

func (c *MemoMiddlewareClient) Start() error {
	fmt.Println("waiting for bucket checking")
	if err := warmup(c.rpcaddr); err != nil {
		return err
	}
	return nil
}

func (c *MemoMiddlewareClient) Stop() error {
	return nil
}

// Get returns Blob for each given ID, or an error.
func (c *MemoMiddlewareClient) Get(ctx context.Context, ids [][]byte) ([]eth.Data, error) {
	var datalist []eth.Data
	for _, id := range ids {
		data, err := getObject(ctx, c.rpcaddr, string(id))
		if err != nil {
			return datalist, err
		}
		datalist = append(datalist, data)
	}
	return datalist, nil
}

// Submit submits the Blobs to Data Availability layer.
func (c *MemoMiddlewareClient) Submit(ctx context.Context, data [][]byte, gasPrice float64) ([][]byte, [][]byte, error) {
	var ids [][]byte
	for _, item := range data {
		id, err := putObject(ctx, c.rpcaddr, item)
		if err != nil {
			return ids, nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil, nil
}

//
// http request api for DA requesting
//

func warmup(rpcaddr string) error {
	req, err := http.Get(rpcaddr + initRoute)
	if err != nil {
		return err
	}
	if req.StatusCode != http.StatusOK {
		return nil
	}
	return nil
}

func getObject(ctx context.Context, rpcaddr string, id string) ([]byte, error) {
	// init request
	req, err := http.NewRequestWithContext(ctx, "GET", rpcaddr+getRoute, nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("id", id)
	req.URL.RawQuery = query.Encode()
	// send request and get response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DA: failed to get object due to bad status code")
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func putObject(ctx context.Context, rpcaddr string, data []byte) ([]byte, error) {
	// init payload and request
	payload := make(map[string]string)
	hexdata := common.Bytes2Hex(data)
	payload["data"] = hexdata
	b, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, "POST", rpcaddr+putRoute, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	// send request and get response
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("DA: failed to put object due to bad status code")
	}
	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	res := make(map[string]string)
	if err = json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	if mid, ok := res["mid"]; !ok {
		return nil, fmt.Errorf("DA: no mid is returned after putObject")
	} else {
		return []byte(mid), nil
	}
}
