package filecoin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

func requestLotus(req request) (*response, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(httpsAddr, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rsp response
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return nil, err
	}

	if rsp.Error != nil {
		return nil, xerrors.New(rsp.Error.Message)
	}

	return &rsp, nil
}

// chainGetMessage lotus ChainGetMessage api
func chainGetMessage(tx string) error {
	c, err := cid.Decode(tx)
	if err != nil {
		return err
	}

	serializedParams, err := json.Marshal(params{
		c,
	})
	if err != nil {
		return err
	}

	req := request{
		Jsonrpc: "2.0",
		Method:  "Filecoin.ChainGetMessage",
		Params:  serializedParams,
		ID:      1,
	}

	rsp, err := requestLotus(req)
	if err != nil {
		return err
	}

	fmt.Printf("rsp:%v \n", rsp)
	b, err := json.Marshal(rsp.Result)
	if err != nil {
		return err
	}

	fmt.Printf("rsp:%v \n", string(b))

	return nil
}
