package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

type Order struct {
	Type            string      `json:"type"`
	Tranid          string      `json:"tranid"`
	Trandate        string      `json:"trandate"`
	Entity          string      `json:"entity"`
	Entityid        string      `json:"entityid"`
	Subsidiary      string      `json:"subsidiary"`
	Location        string      `json:"location"`
	Status          string      `json:"status"`
	Shipdate        string      `json:"shipdate"`
	Shippingaddress string      `json:"shippingaddress"`
	Billingaddress  string      `json:"billingaddress"`
	Memo            string      `json:"memo"`
	Currency        string      `json:"currency"`
	Exchangerate    string      `json:"exchangerate"`
	Item            []OrderItem `json:"item"`
}

type OrderItem struct {
	Item     string `json:"item"`
	Rate     string `json:"rate"`
	Quantity string `json:"quantity"`
}

// prefix
const (
	PrefixOrder = "order"
)

// main
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// init
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// invoke
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoking " + function)

	// functions
	if function == "createOrder" {
		return t.createOrder(stub, args)
	} else if function == "getOrderById" {
		return t.getOrderById(stub, args)
	}

	// result
	message := "invoke did not find func: " + function
	fmt.Println(message)
	return shim.Error(message)
}

// ===========================================================
// createOrder
// ===========================================================
func (t *SimpleChaincode) createOrder(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	// check args
	count := 1
	if len(args) != count {
		return shim.Error(fmt.Sprintf("Incorrect number of arguments. Expecting %d", count))
	}

	for i := 0; i < len(args); i++ {
		if len(args[i]) <= 0 {
			return shim.Error(fmt.Sprintf("argument must be a non-empty string, index=%d", i))
		}
	}

	var order Order
	json.Unmarshal([]byte(args[0]), &order)
	// fmt.Printf("arg0=%s\n", args[0])
	fmt.Printf("id=%s\n", order.Tranid)

	// create
	key, err := stub.CreateCompositeKey(PrefixOrder, []string{order.Tranid})
	if err != nil {
		return shim.Error(err.Error())
	}

	bytes, err := json.Marshal(order)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(key, bytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("create success")
	return shim.Success(nil)
}

// ===========================================================
// getOrderById
// ===========================================================
func (t *SimpleChaincode) getOrderById(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	var id, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting id to query")
	}

	id = args[0]

	key, err := stub.CreateCompositeKey(PrefixOrder, []string{id})
	if err != nil {
		return shim.Error(err.Error())
	}

	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + id + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"does not exist: " + id + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}
