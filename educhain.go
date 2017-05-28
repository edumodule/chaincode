/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {

	_, args := stub.GetFunctionAndParameters()
	var Hash, Certificate, Name string
	var empty []string
	var err error

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	Hash = args[0]

	Certificate = args[1]

	Name = args[2]

	err = stub.PutState(Hash, []byte(Certificate))
	if err != nil {
		return shim.Error(err.Error())
	}

	hashes, err := json.Marshal(empty)
	if err != nil {
		return shim.Error(err.Error())
	}

	err = stub.PutState(Name, hashes)
	if err != nil {
		return shim.Error(err.Error())
	}

	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)

}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

	function, args := stub.GetFunctionAndParameters()

	if function != "invoke" {
		return shim.Error("Unknown function call")
	}

	if len(args) < 2 {
		return shim.Error("Incorrect number of arguments. Expecting at least 2")
	}

	if args[0] == "setCertificate" {
		return t.setCertificate(stub, args)
	}

	if args[0] == "getCertificates" {
		return t.getCertificates(stub, args)
	}

	return shim.Error("Unknown action, check the first argument, must be one of 'setCertificate' or 'getCertificates'")
}

func (t *SimpleChaincode) setCertificate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var Name, Hash, Certificate string
	var hashesArray []string
	var err error

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	Name = args[1]
	Hash = args[2]
	Certificate = args[3]

	hashes, err := stub.GetState(Name)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal(hashes, hashesArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	hashesArray = append(hashesArray, Hash)

	err = stub.PutState(Hash, []byte(Certificate))
	if err != nil {
		return shim.Error(err.Error())
	}

	byteArray, _ := json.Marshal(hashesArray)

	err = stub.PutState(Name, byteArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	if transientMap, err := stub.GetTransient(); err == nil {
		if transientData, ok := transientMap["result"]; ok {
			return shim.Success(transientData)
		}
	}
	return shim.Success(nil)
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) getCertificates(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	var Name string
	var hashesArray, response []string
	var err error

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	Name = args[1]

	hashes, err := stub.GetState(Name)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state\"}"
		return shim.Error(jsonResp)
	}

	err = json.Unmarshal(hashes, hashesArray)
	if err != nil {
		return shim.Error(err.Error())
	}

	for _, v := range hashesArray {
		hash, lerr := stub.GetState(v)
		if lerr != nil {
			fmt.Println(lerr)
		}
		response = append(response, string(hash))
	}

	res, _ := json.Marshal(response)

	return shim.Success(res)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
