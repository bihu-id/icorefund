// Copyright 2009 The Go Authors. All rights reserved.
// Copyright 2012 The Gorilla Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonrpc

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
)

// ----------------------------------------------------------------------------
// Request and Response
// ----------------------------------------------------------------------------

// ReqObj represents a JSON-RPC request sent by a client.
type ReqObj struct {
	// JSON-RPC protocol.
	Version string `json:"jsonrpc"`

	// A String containing the name of the method to be invoked.
	Method string `json:"method"`

	// Object to pass as request parameter to the method.
	Params interface{} `json:"params"`

	// The request id. This can be of any type. It is used to match the
	// response with the request that it is replying to.
	Id uint64 `json:"id"`
}

// ResObj represents a JSON-RPC response returned to a client.
type ResObj struct {
	Version string           `json:"jsonrpc"`
	Result  *json.RawMessage `json:"result"`
	Error   *json.RawMessage `json:"error"`
	Id      uint64           `json:"id"`
}

// JSON-RPC error object
type ErrorCode int

const (
	E_PARSE       ErrorCode = -32700
	E_INVALID_REQ ErrorCode = -32600
	E_NO_METHOD   ErrorCode = -32601
	E_BAD_PARAMS  ErrorCode = -32602
	E_INTERNAL    ErrorCode = -32603
	E_SERVER      ErrorCode = -32000
)

type Error struct {
	// A Number that indicates the error type that occurred.
	Code ErrorCode `json:"code"` /* required */

	// A String providing a short description of the error.
	// The message SHOULD be limited to a concise single sentence.
	Message string `json:"message"` /* required */

	// A Primitive or Structured value that contains additional information about the error.
	Data interface{} `json:"data"` /* optional */
}

func (e *Error) Error() string {
	return e.Message
}

// ResponseRecorder is an implementation of http.ResponseWriter that
// records its mutations for later inspection in tests.
type ResponseRecorder struct {
	Code      int           // the HTTP response code from WriteHeader
	HeaderMap http.Header   // the HTTP response headers
	Body      *bytes.Buffer // if non-nil, the bytes.Buffer to append written data to
	Flushed   bool
}

// NewRecorder returns an initialized ResponseRecorder.
func NewRecorder() *ResponseRecorder {
	return &ResponseRecorder{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
	}
}

// EncodeReqObj encodes parameters for a JSON-RPC client request.
func EncodeReqObj(method string, args interface{}) ([]byte, error) {
	c := &ReqObj{
		Version: "2.0",
		Method:  method,
		Params:  args,
		Id:      uint64(rand.Int63()),
	}
	return json.Marshal(c)
}

// DecodeResObj decodes the response body of a client request into
// the interface reply.
func DecodeResObj(r io.Reader, reply interface{}) error {
	var c ResObj
	if err := json.NewDecoder(r).Decode(&c); err != nil {
		return err
	}
	//log.Printf("ResObj: %s", c)
	if c.Error != nil {
		jsonErr := &Error{}
		if err := json.Unmarshal(*c.Error, jsonErr); err != nil {
			return &Error{
				Code:    E_SERVER,
				Message: string(*c.Error),
			}
		}
		return jsonErr
	}
	if c.Result == nil {
		return &Error{Message: "response body is null"}
	}
	return json.Unmarshal(*c.Result, reply)
}

func Call(url string, method string, params interface{}, reply interface{}) error {

	j, err := EncodeReqObj(method, params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	//log.Printf("result.Body: %s", resp.Body)
	return DecodeResObj(resp.Body, &reply)
}
