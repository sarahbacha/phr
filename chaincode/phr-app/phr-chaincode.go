// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
*/

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

/* Imports
* 4 utility libraries for handling bytes, reading and writing JSON,
formatting, and string manipulation
* 2 specific Hyperledger Fabric specific libraries for Smart Contracts
*/

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Patient structure, with 6 properties.
Structure tags are used by encoding/json library
*/

type Account struct {
	Id        string `json:"id"`
	FirstName string `json:"first"`
	LastName  string `json:"last"`
	UserName  string `json:"name"`
	Password  string `json:"password"`
}
type Patient struct {
	AccountId string `json:"account"`
	Gender    string `json:"Gender"`
	Age       string `json:"Age"`
	HbA1c     string `json:"HbA1c"`
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
	switch function {
	case "initLedger":
		return s.initLedger(APIstub)
	case "addAccount":
		return s.account_add(APIstub, args)
	case "readAccount":
		return s.account_read(APIstub, args)
	case "queryPatient":
		return s.queryPatient(APIstub, args)
	case "recordPatient":
		return s.recordPatient(APIstub, args)
	case "editPatient":
		return s.editPatient(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

func (s *SmartContract) account_add(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	bytes, err := APIstub.GetState("phr_accounts")

	if err != nil {
		return shim.Error("Unable to get accounts.")
	}

	var account Account

	// Build JSON values
	id := "\"id\": \"" + args[0] + "\", "
	first := "\"first\": \"" + args[1] + "\", "
	last := "\"last\": \"" + args[2] + "\", "
	name := "\"name\": \"" + args[3] + "\", "
	password := "\"password\": \"" + args[4] + "\""

	// Make into a complete JSON string
	// Decode into a single account value
	content := "{" + id + first + last + name + password + "}"
	err = json.Unmarshal([]byte(content), &account)
	fmt.Printf("Query Response  content :\n", content)
	var accounts []Account

	// Decode JSON into account array
	// Add latest account
	err = json.Unmarshal(bytes, &accounts)
	accounts = append(accounts, account)

	// Encode as JSON
	// Put back on the block
	bytes, err = json.Marshal(accounts)
	err = APIstub.PutState("phr_accounts", bytes)
	return shim.Success(nil)
}
func (s *SmartContract) account_read(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	bytes, err := APIstub.GetState("phr_accounts")
	if err != nil {
		return shim.Error("Unable to get accounts." + err.Error())
	}

	var accounts []Account

	// From JSON to data structure
	err = json.Unmarshal(bytes, &accounts)
	found := false

	// Look for match
	for _, account := range accounts {
		// Match
		if account.UserName == args[0] && account.Password == args[1] {
			// Sanitize
			account.Password = ""

			// JSON encode
			bytes, err = json.Marshal(account)
			found = true
			break
		}
	}

	// Nope
	if found != true {
		bytes, err = json.Marshal(nil)
	}
	//fmt.Printf("Query Response:%s\n", bytes)
	return shim.Success(bytes)
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
	// // Accounts
	// fmt.Println("PHR Is Starting Up")

	// var accounts []Account
	// bytes, err := json.Marshal(accounts)

	// if err != nil {
	// 	return shim.Error("Error initializing accounts.")
	// }
	// err = APIstub.PutState("phr_accounts", bytes)

	// // Patient
	// var patients []Patient

	// bytes, err = json.Marshal(patients)

	// if err != nil {
	// 	return shim.Error("Error initializing patients.")
	// }

	// err = APIstub.PutState("phr_patients", bytes)

	// return shim.Success(nil)
	patient := []Patient{}

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

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var patient = Patient{AccountId: args[0], Gender: args[1], Age: args[2], HbA1c: args[3], HbA1cDate: args[4]}

	patientAsBytes, _ := json.Marshal(patient)
	err := APIstub.PutState(args[0], patientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record patient: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The changePatientDiabetesIndication method *
The data in the world state can be updated with who has possession.
This function takes in 3 arguments, patient id, new HbA1c Percentage and Date.
*/
func (s *SmartContract) editPatient(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	if patientAsBytes == nil {
		return shim.Error("Could not locate patient")
	}
	patient := Patient{}

	json.Unmarshal(patientAsBytes, &patient)
	// Do a check on value passed
	// we are skipping this check for this example
	patient.Age = args[1]
	patient.Gender = args[2]
	patient.HbA1c = args[3]
	patient.HbA1cDate = args[4]

	patientAsBytes, _ = json.Marshal(patient)
	err := APIstub.PutState(args[0], patientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to edit patient ", args[0]))
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
