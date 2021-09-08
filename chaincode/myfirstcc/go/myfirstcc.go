package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type SimpleAsset struct {
}
type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (t *SimpleAsset) Init(stub shim.ChaincodeStubInterface) peer.Response {

	return shim.Success(nil)
}

func (t *SimpleAsset) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	fn, args := stub.GetFunctionAndParameters()

	var result string
	var err error
	if fn == "set" {
		result, err = set(stub, args)
	} else if fn == "get" {
		result, err = get(stub, args)
	} else if fn == "getAllKeys" {
		result, err = getAllKeys(stub)
	} else if fn == "getHistoryForKey" {
		result, err = getHistoryForKey(stub, args)
	} else {
		return shim.Error("Not supported chaincode function.")
	}

	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(result))
}

func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key")
	}

	value, err := stub.GetState(args[0])
	if err != nil {
		return "", fmt.Errorf("Failed to get asset: %s with error: %s", args[0], err)
	}
	if value == nil {
		return "", fmt.Errorf("Asset not found: %s", args[0])
	}
	return string(value), nil
}

func getAllKeys(stub shim.ChaincodeStubInterface) (string, error) {
	iter, err := stub.GetStateByRange("A", "z")
	if err != nil {
		return "", fmt.Errorf("Failed to get all keys with error: %s", err)
	}
	defer iter.Close()

	var buffer string
	buffer = "["

	comma := false
	for iter.HasNext() {
		res, err := iter.Next()
		if err != nil {
			return "", fmt.Errorf("Failed to get iterator's next value: %s", err)
		}
		if comma == true {
			buffer += ","
		}
		buffer += string(res.Value)

		comma = true
	}
	buffer += "]"

	fmt.Println(buffer)

	return string(buffer), nil
}

func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 2 {
		return "", fmt.Errorf("Incorrect arguments. Expecting a key and a value")
	}

	var data = Data{Key: args[0], Value: args[1]}
	dataAsBytes, _ := json.Marshal(data)

	err := stub.PutState(args[0], dataAsBytes)
	if err != nil {
		return "", fmt.Errorf("Failed to set asset: %s", args[0])
	}
	return string(dataAsBytes), nil
}

func getHistoryForKey(stub shim.ChaincodeStubInterface, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("Incorrect number of arguments. Expecting 1")
	}
	key := args[0]

	fmt.Println("getHistoryForKey: " + key)

	iter, err := stub.GetHistoryForKey(key)
	if err != nil {
		return "", err
	}
	defer iter.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	comma := false
	for iter.HasNext() {
		res, err1 := iter.Next()
		if err1 != nil {
			return "", err1
		}
		if comma == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(res.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		if res.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(res.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(res.Timestamp.Seconds, int64(res.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(res.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		comma = true
	}
	buffer.WriteString("]")

	return (string)(buffer.Bytes()), nil
}

func main() {
	if err := shim.Start(new(SimpleAsset)); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
