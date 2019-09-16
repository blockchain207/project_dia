/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the dia structure, with 4 properties.  Structure tags are used by encoding/json library
type Dia struct {
	Carat   string `json:"carat"`
	Color   string `json:"color"`
	Clarity string `json:"clarity"`
	Owner   string `json:"owner"`
	Lost	bool   `json:"lost"`
}

/*
 * The Init method is called when the Smart Contract "fabdia" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "fabdia"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryDia" {
		return s.queryDia(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createDia" {
		return s.createDia(APIstub, args)
	} else if function == "queryAllDias" {
		return s.queryAllDias(APIstub)
	} else if function == "changeDiaOwner" {
		return s.changeDiaOwner(APIstub, args)
	} else if function == "lostDia" {
		return s.lostDia(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

func (s *SmartContract) queryDia(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	diaAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(diaAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	dias := []Dia{
		Dia{Carat: "2.03", Color: "E", Clarity: "SI1", Owner: "Tomoko", Lost: false},
		Dia{Carat: "1.59", Color: "D", Clarity: "SI3", Owner: "Brad", Lost: false},
		Dia{Carat: "2.12", Color: "F", Clarity: "SI2", Owner: "Jin Soo", Lost: false},
		Dia{Carat: "1.78", Color: "D", Clarity: "SI2", Owner: "Max", Lost: false},
		Dia{Carat: "2.42", Color: "F", Clarity: "SI1", Owner: "Adriana", Lost: false},
		Dia{Carat: "2.11", Color: "E", Clarity: "SI3", Owner: "Michel", Lost: false},
		Dia{Carat: "2.29", Color: "D", Clarity: "SI2", Owner: "Aarav", Lost: false},
		Dia{Carat: "1.98", Color: "F", Clarity: "SI3", Owner: "Pari", Lost: false},
		Dia{Carat: "1.66", Color: "E", Clarity: "SI1", Owner: "Valeria", Lost: false},
		Dia{Carat: "1.84", Color: "E", Clarity: "SI3", Owner: "Shotaro", Lost: false},
	}

	i := 0
	for i < len(dias) {
		fmt.Println("i is ", i)
		diaAsBytes, _ := json.Marshal(dias[i])
		APIstub.PutState("Dia"+strconv.Itoa(i), diaAsBytes)
		fmt.Println("Added", dias[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createDia(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var dia = Dia{Carat: args[1], Color: args[2], Clarity: args[3], Owner: args[4]}

	diaAsBytes, _ := json.Marshal(dia)
	APIstub.PutState(args[0], diaAsBytes)

	json.Unmarshal(diaAsBytes, &dia)
        dia.Lost = false

	return shim.Success(nil)
}

func (s *SmartContract) queryAllDias(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "Dia0"
	endKey := "Dia999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllDias:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeDiaOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	diaAsBytes, _ := APIstub.GetState(args[0])
	dia := Dia{}

	json.Unmarshal(diaAsBytes, &dia)
	dia.Owner = args[1]

	diaAsBytes, _ = json.Marshal(dia)
	APIstub.PutState(args[0], diaAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) lostDia(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

        if len(args) != 1 {
                return shim.Error("Incorrect number of arguments. Expecting 1")
        }

        diaAsBytes, _ := APIstub.GetState(args[0])
        dia := Dia{}

	json.Unmarshal(diaAsBytes, &dia)
	if dia.Lost == false {
		dia.Lost = true
		dia.Owner = "Insure Co"
	}

        diaAsBytes, _ = json.Marshal(dia)
        APIstub.PutState(args[0], diaAsBytes)

        return shim.Success(nil)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
