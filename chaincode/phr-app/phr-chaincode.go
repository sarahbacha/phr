// SPDX-License-Identifier: Apache-2.0

/*
  Sample Chaincode based on Demonstrated Scenario

 This code is based on code written by the Hyperledger Fabric community.
  Original code can be found here: https://github.com/hyperledger/fabric-samples/blob/release/chaincode/fabcar/fabcar.go
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

/* Imports
 * 4 utility libraries for handling bytes, reading and writing JSON,
 * formatting, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */

// Define the Smart Contract structure
type SmartContract struct {
}

/* Define Patient structure, with 6 properties.
 * Structure tags are used by encoding/json library
 */

type Account struct {
	ObjectType string `json:"docType"`
	Id         string `json:"id"`
	Role       string `json:"role"`
	FirstName  string `json:"first"`
	LastName   string `json:"last"`
	UserName   string `json:"name"`
	Password   string `json:"password"`
}
type PatientGeneralInformation struct {
	ObjectType string `json:"docType"`
	Id         string `json:"id"`
	AccountId  string `json:"account"`
	DOB        string `json:"DOB"`
	Sex        string `json:"Sex"`
	Race       string `json:"Race"`
}
type PatientDiabetesIndications struct {
	ObjectType string `json:"docType"`
	Id         string `json:"id"`
	AccountId  string `json:"account"`
	HbA1c      string `json:"HbA1c"`
	HbA1cDate  string `json:"HbA1cDate"`
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
	// Account Actions
	case "addAccount":
		return s.account_add(APIstub, args)
	case "readAccount":
		return s.account_read(APIstub, args)
	// Patient Information:
	case "queryPatient":
		return s.queryPatientByAccount(APIstub, args)
	// Patient General Information Actions
	case "recordPatientGeneralInformation":
		return s.recordPatientGeneralInformation(APIstub, args)
	case "editPatientGeneralInformation":
		return s.editPatientGeneralInformation(APIstub, args)
	case "queryPatientGeneralInformation":
		return s.queryPatientGeneralInformation(APIstub, args)
	case "queryPatientGeneralInformationHistory":
		return s.queryPatientGeneralInformationHistory(APIstub, args)
	// Patient Diabetes Indications Actions
	case "recordPatientDiabetesIndications":
		return s.recordPatientDiabetesIndications(APIstub, args)
	case "editPatientDiabetesIndications":
		return s.editPatientDiabetesIndications(APIstub, args)
	case "queryPatientDiabetesIndications":
		return s.queryPatientDiabetesIndications(APIstub, args)
	case "queryPatientDiabetesIndicationsHistory":
		return s.queryPatientDiabetesIndicationsHistory(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}

/*
 * The initLedger method *
 * Will add test data (1 patient)to our network
 */
func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	// Accounts
	fmt.Println("PHR Is Starting Up")

	var accounts []Account
	bytes, err := json.Marshal(accounts)

	if err != nil {
		return shim.Error("Error initializing accounts.")
	}
	err = APIstub.PutState("phr_accounts", bytes)

	return shim.Success(nil)
}

// Account Transactions
func (s *SmartContract) account_add(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	bytes, err := APIstub.GetState("phr_accounts")

	if err != nil {
		return shim.Error("Unable to get accounts.")
	}

	var account Account
	// Build JSON values
	objectType := "\"docType\": \"" + "account" + "\", "
	id := "\"id\": \"" + args[0] + "\", "
	role := "\"role\": \"" + args[1] + "\", "
	first := "\"first\": \"" + args[2] + "\", "
	last := "\"last\": \"" + args[3] + "\", "
	name := "\"name\": \"" + args[4] + "\", "
	password := "\"password\": \"" + args[5] + "\""

	// Make into a complete JSON string
	// Decode into a single account value
	content := "{" + objectType + id + role + first + last + name + password + "}"
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

// General Patient Transactions
/*
 * The recordPatient method *
 * This method takes in five arguments (attributes to be saved in the ledger).
 */
func (s *SmartContract) recordPatientGeneralInformation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	var patient = PatientGeneralInformation{
		ObjectType: "PatientGeneralInformation",
		Id:         args[0],
		AccountId:  args[1],
		DOB:        args[2],
		Sex:        args[3],
		Race:       args[4],
	}

	patientAsBytes, _ := json.Marshal(patient)
	err := APIstub.PutState(args[0], patientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record patient general information: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The editPatientGeneralInformation method *
 * The data in the world state can be updated with who has possession.
 * This function takes in 3 arguments, patient id, new HbA1c Percentage and Date.
 */
func (s *SmartContract) editPatientGeneralInformation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	if patientAsBytes == nil {
		return shim.Error("Could not locate PatientGeneralInformation")
	}
	patient := PatientGeneralInformation{}

	json.Unmarshal(patientAsBytes, &patient)

	patient.ObjectType = "PatientGeneralInformation"
	patient.Id = args[0]
	patient.AccountId = args[1]
	patient.DOB = args[2]
	patient.Sex = args[3]
	patient.Race = args[4]

	patientAsBytes, _ = json.Marshal(patient)
	err := APIstub.PutState(args[0], patientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to edit PatientGeneralInformation ", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryPatientGeneralInformation method *
 * Used to view General Information of one particular patient
 * It takes one argument -- the key for the patient in question
 */
func (s *SmartContract) queryPatientGeneralInformation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	if patientAsBytes == nil {
		return shim.Error("Could not locate Patient General Information")
	}
	return shim.Success(patientAsBytes)
}

/*
 * The queryPatientGeneralInformationHistory method *
 * Used to view the General Information History of one particular patient
 * It takes one argument -- the key for the patient in question
 */
func (s *SmartContract) queryPatientGeneralInformationHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	patientGeneralInformationID := args[0]
	fmt.Printf("- start queryPatientGeneralInformationHistory: %s\n", patientGeneralInformationID)
	resultsIterator, err := APIstub.GetHistoryForKey(patientGeneralInformationID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryPatientGeneralInformationHistory returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// Diabetes Indications Transactions
/*
 * The recordPatientDiabetesIndications method *
 * This method takes in four arguments (attributes to be saved in the ledger).
 */
func (s *SmartContract) recordPatientDiabetesIndications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	var patient = PatientDiabetesIndications{
		ObjectType: "PatientDiabetesIndications",
		Id:         args[0],
		AccountId:  args[1],
		HbA1c:      args[2],
		HbA1cDate:  args[3],
	}

	patientAsBytes, _ := json.Marshal(patient)
	err := APIstub.PutState(args[0], patientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record patient diabetes indications: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The editPatientGeneralInformation method *
 * The data in the world state can be updated with who has possession.
 */
func (s *SmartContract) editPatientDiabetesIndications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	if patientAsBytes == nil {
		return shim.Error("Could not locate PatientDiabetesIndications")
	}
	patient := PatientDiabetesIndications{}

	json.Unmarshal(patientAsBytes, &patient)

	patient.ObjectType = "PatientDiabetesIndications"
	patient.Id = args[0]
	patient.AccountId = args[1]
	patient.HbA1c = args[2]
	patient.HbA1cDate = args[3]

	patientAsBytes, _ = json.Marshal(patient)
	err := APIstub.PutState(args[0], patientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to edit PatientGeneralInformation ", args[0]))
	}

	return shim.Success(nil)
}

/*
 * The queryPatientDiabetesIndications method *
 * Used to view Diabetes Indications of one particular patient
 * It takes one argument -- the key for the patient in question
 */
func (s *SmartContract) queryPatientDiabetesIndications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	if patientAsBytes == nil {
		return shim.Error("Could not locate Patient Diabetes Indications")
	}
	return shim.Success(patientAsBytes)
}

/*
 * The queryPatientDiabetesIndicationsHistory method *
 * Used to view the Diabetes Indications  History of one particular patient
 * It takes one argument -- the key for the patient in question
 */
func (s *SmartContract) queryPatientDiabetesIndicationsHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	patientDiabetesIndicationsID := args[0]
	fmt.Printf("- start queryPatientDiabetesIndicationsHistory: %s\n", patientDiabetesIndicationsID)
	resultsIterator, err := APIstub.GetHistoryForKey(patientDiabetesIndicationsID)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryPatientDiabetesIndicationsHistory returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryPatientByAccount(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	accountID := args[0]

	queryString := fmt.Sprintf("{\"selector\":{\"account\":{\"$eq\":\"%s\"}}}", accountID)

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

// =========================================================================================
// getQueryResultForQueryString executes the passed in query string.
// Result set is built and returned as a byte array containing the JSON results.
// =========================================================================================
func getQueryResultForQueryString(APIstub shim.ChaincodeStubInterface, queryString string) ([]byte, error) {

	fmt.Printf("- getQueryResultForQueryString queryString:\n%s\n", queryString)

	resultsIterator, err := APIstub.GetQueryResult(queryString)
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	buffer, err := constructQueryResponseFromIterator(resultsIterator)
	if err != nil {
		return nil, err
	}

	fmt.Printf("- getQueryResultForQueryString queryResult:\n%s\n", buffer.String())

	return buffer.Bytes(), nil
}

func constructQueryResponseFromIterator(resultsIterator shim.StateQueryIteratorInterface) (*bytes.Buffer, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
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

	return &buffer, nil
}

/*
 * main function *
	*calls the Start function
	*The main function starts the chaincode in the container during instantiation.
*/
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
