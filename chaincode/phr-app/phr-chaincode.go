// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
 */

package main

/* Imports  
* 4 utility libraries for handling bytes, reading and writing JSON, 
formatting, and string manipulation  
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

/* Define Patient structure, with 6 properties.  
Structure tags are used by encoding/json library
*/
type Patient struct {
	ID string `json:"ID"`
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Age string `json:"Age"`
	HbA1c string `json:"HbA1c"`
	HbA1cDate string `json:"HbA1cDate"`
}


/*
 * The Init method *
 called when the Smart Contract "phr-chaincode" is instantiated by the network
 * Best practice is to have any Ledger initialization in separate function 
 -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method *
 called when an application requests to run the Smart Contract "phr-chaincode"
 The app also specifies the specific smart contract function to call with args
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger
	if function == "queryPatient" {
		return s.queryPatient(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "recordPatient" {
		return s.recordPatient(APIstub, args)
	} else if function == "queryAllPatients" {
		return s.queryAllPatients(APIstub)
	} else if function == "changePatientDiabetesIndication" {
		return s.changePatientDiabetesIndication(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

/*
 * The queryPatient method *
Used to view the records of one particular patient
It takes one argument -- the key for the patient in question
 */
func (s *SmartContract) queryPatient(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	if patientAsBytes == nil {
		return shim.Error("Could not locate patient")
	}
	return shim.Success(patientAsBytes)
}

/*
 * The initLedger method *
Will add test data (1 patient)to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	patient := []Patient{
		Patient{ID: "1", FirstName: "Sarah", LastName: "Bacha", Age: "24", HbA1c:"9", HbA1cDate:"10-10-2018"},
	}

	i := 0
	for i < len(patient) {
		fmt.Println("i is ", i)
		patientAsBytes, _ := json.Marshal(patient[i])
		APIstub.PutState(strconv.Itoa(i+1), patientAsBytes)
		fmt.Println("Added", patient[i])
		i = i + 1
	}

	return shim.Success(nil)
}

/*
 * The recordPatient method *
This method takes in five arguments (attributes to be saved in the ledger). 
 */
func (s *SmartContract) recordPatient(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var patient = Patient{ ID: args[0], FirstName: args[1], LastName: args[2], Age: args[3], HbA1c:"", HbA1cDate:""}

	patientAsBytes, _ := json.Marshal(patient)
	err := APIstub.PutState(args[0], patientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record patient: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryAllPatients method *
allows for assessing all the records added to the ledger(all patients)
This method does not take any arguments. Returns JSON string containing results. 
 */
func (s *SmartContract) queryAllPatients(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "1"
	endKey := "999"

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
		// Add comma before array members,suppress it for the first array member
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

	fmt.Printf("- queryAllPatients:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

/*
 * The changePatientDiabetesIndication method *
The data in the world state can be updated with who has possession. 
This function takes in 3 arguments, patient id, new HbA1c Percentage and Date. 
 */
func (s *SmartContract) changePatientDiabetesIndication(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	if patientAsBytes == nil {
		return shim.Error("Could not locate patient")
	}
	patient := Patient{}

	json.Unmarshal(patientAsBytes, &patient)
	// Do a check on value passed
	// we are skipping this check for this example
	patient.HbA1c = args[1]
	patient.HbA1cDate = args[2]

	patientAsBytes, _ = json.Marshal(patient)
	err := APIstub.PutState(args[0], patientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to change patient HbA1c Percentage: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * main function *
calls the Start function 
The main function starts the chaincode in the container during instantiation.
 */
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}