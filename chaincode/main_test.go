package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"

	"testing"
)

var simpleAsset = new(SimpleChaincode)
var stub = shim.NewMockStub("mockstub", simpleAsset)

// ===========================================================
// test
// ===========================================================
func Test(t *testing.T) {
	// JSON.stringify(json)
	call("createOrder", []string{`{"type":"purchaseorder","tranid":"RSHKPO20001394","trandate":"8/7/2020","entity":"Colgate-Palmolive (HK) Limited 高露潔","entityid":"VCP001","subsidiary":"RSHK","location":"RSHK","status":"Pending Receipt","shipdate":"11/7/2020","shippingaddress":"","billingaddress":"18th floor, Causeway Bay Plaza I, 489 Hennessy Road","memo":"","item":[{"item":"A-AX0007","rate":"21.57","quantity":300},{"item":"A-AX0005","rate":"21.57","quantity":300},{"item":"A-AX0023","rate":"21.57","quantity":120}]}`})
	call("getOrderById", []string{"RSHKPO20001394"})
}

// ===========================================================
// private call
// ===========================================================
func call(method string, argsArray []string) {
	fmt.Println("------ start ------ ", method)
	// all args
	var args [][]byte
	args = append(args, []byte(method))
	fmt.Printf("- args=[")
	fmt.Printf("p0=%s", method)
	if argsArray != nil {
		for i := 0; i < len(argsArray); i++ {
			args = append(args, []byte(argsArray[i]))
			fmt.Printf(",p%d=%s", i+1, argsArray[i])
		}
	}
	fmt.Printf("]")
	fmt.Println("")
	// invoke
	response := stub.MockInvoke("uuid", args)
	fmt.Printf("- status=")
	fmt.Println(response.GetStatus())
	fmt.Printf("- error message=")
	fmt.Println(response.GetMessage())
	fmt.Printf("- payload=")
	fmt.Println(string(response.GetPayload()))
	fmt.Println("------ end ------ ")
	fmt.Println("")
}
