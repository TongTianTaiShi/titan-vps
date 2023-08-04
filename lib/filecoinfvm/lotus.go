package filecoinfvm

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/ipfs/go-cid"
	"golang.org/x/xerrors"
)

type request struct {
	Jsonrpc string     `json:"jsonrpc"`
	Method  string     `json:"method"`
	Params  rawMessage `json:"params"`
	ID      int        `json:"id"`
}

type rawMessage []byte

// MarshalJSON returns m as the JSON encoding of m.
func (m rawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	return m, nil
}

// Response defines a JSON RPC response from the spec
// http://www.jsonrpc.org/specification#response_object
type response struct {
	Jsonrpc string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	ID      interface{} `json:"id"`
	Error   *respError  `json:"error,omitempty"`
}

type respError struct {
	Code    errorCode       `json:"code"`
	Message string          `json:"message"`
	Meta    json.RawMessage `json:"meta,omitempty"`
}

type errorCode int

// UnmarshalJSON sets *m to a copy of data.
func (m *rawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*m = append((*m)[0:0], data...)
	return nil
}

type params []interface{}

func requestLotus(out interface{}, req request, addr string) error {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return err
	}

	resp, err := http.Post(addr, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var rsp response
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return err
	}

	if rsp.Error != nil {
		return xerrors.New(rsp.Error.Message)
	}

	b, err := json.Marshal(rsp.Result)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, out)
}

// ChainGetMessage lotus ChainGetMessage api
func ChainGetMessage(out interface{}, c cid.Cid, addr string) error {
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

	return requestLotus(out, req, addr)
}

// StateSearchMsg lotus stateSearchMsg api
func StateSearchMsg(out interface{}, c cid.Cid, addr string) error {
	serializedParams, err := json.Marshal(params{
		c,
	})
	if err != nil {
		return err
	}

	req := request{
		Jsonrpc: "2.0",
		Method:  "Filecoin.StateSearchMsg",
		Params:  serializedParams,
		ID:      1,
	}

	return requestLotus(out, req, addr)
}

// EthGetMessageCidByTransactionHash
func EthGetMessageCidByTransactionHash(out interface{}, tx string, addr string) error {
	serializedParams, err := json.Marshal(params{
		tx,
	})
	if err != nil {
		return err
	}

	req := request{
		Jsonrpc: "2.0",
		Method:  "Filecoin.EthGetMessageCidByTransactionHash",
		Params:  serializedParams,
		ID:      1,
	}

	return requestLotus(out, req, addr)
}
