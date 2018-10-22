package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type SmartContract struct {
}
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
	ObjectType       string `json:"docType"`
	Id               string `json:"id"`
	AccountId        string `json:"account"`
	DOB              string `json:"DOB"`
	Sex              string `json:"Sex"`
	Race             string `json:"Race"`
	PatientReporting string `json:"reportingid"`
}
type PatientReporting struct {
	ObjectType string `json:"docType"`
	PatientId  string `json:"patient"`
	ReporterId string `json:"reporter"`
	Status     string `json:"status"`
}
type PatientDiabetesIndications struct {
	ObjectType       string `json:"docType"`
	Id               string `json:"id"`
	AccountId        string `json:"account"`
	HbA1c            string `json:"HbA1c"`
	HbA1cDate        string `json:"HbA1cDate"`
	PatientReporting string `json:"reportingid"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}
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
	// Patient Reporting Actions
	case "recordPatientReporting":
		return s.recordPatientReporting(APIstub, args)
	case "getPatientReportingAccess":
		return s.getPatientReportingAccess(APIstub, args)
	default:
		return shim.Error("Invalid Smart Contract function name.")
	}
}
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

	if strings.HasPrefix(args[1], "Patient") {
		var patientReporting = PatientReporting{
			ObjectType: "PatientReporting",
			PatientId:  args[0],
			ReporterId: args[0],
			Status:     "Active",
		}
		reportingbytes, errreporting := json.Marshal(patientReporting)
		patientID := args[0]
		values := []string{}
		values = append(values, patientID)
		values = append(values, patientID)
		key := strings.Join(values, "")
		errreporting = APIstub.PutState(key, reportingbytes)
		if errreporting != nil {
			return shim.Error("Unable to put reporting." + errreporting.Error())
		}
	}

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
	Record Patient General Information
	arguments: [id, patientID, reporterID, dob, sex, race]
*/
func (s *SmartContract) recordPatientGeneralInformation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	status, err := s.getPatientReportingStatus(APIstub, args)
	if err != nil {
		return shim.Error(err.Error())
	}
	if !strings.HasPrefix(status, "Active") {
		return shim.Error(fmt.Sprintf("Transaction is refused"))
	}
	patientID := args[1]
	reporterID := args[2]
	values := []string{}
	values = append(values, patientID)
	values = append(values, reporterID)
	key := strings.Join(values, "")

	var patient = PatientGeneralInformation{
		ObjectType:       "PatientGeneralInformation",
		Id:               args[0],
		AccountId:        args[1],
		DOB:              args[3],
		Sex:              args[4],
		Race:             args[5],
		PatientReporting: key,
	}

	patientAsBytes, _ := json.Marshal(patient)
	erro := APIstub.PutState(args[0], patientAsBytes)
	if erro != nil {
		return shim.Error(fmt.Sprintf("Failed to record patient general information: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
	Edit General Patient Information
	arguments: [id, patientID, reporterID, dob, sex, race]
*/
func (s *SmartContract) editPatientGeneralInformation(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}
	status, err := s.getPatientReportingStatus(APIstub, args)
	if err != nil {
		return shim.Error(err.Error())
	}
	if !strings.HasPrefix(status, "Active") {
		return shim.Error(fmt.Sprintf("Transaction is refused"))
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	if patientAsBytes == nil {
		return shim.Error("Could not locate PatientGeneralInformation")
	}
	patient := PatientGeneralInformation{}

	json.Unmarshal(patientAsBytes, &patient)

	patientID := args[1]
	reporterID := args[2]
	values := []string{}
	values = append(values, patientID)
	values = append(values, reporterID)
	key := strings.Join(values, "")

	patient.ObjectType = "PatientGeneralInformation"
	patient.Id = args[0]
	patient.AccountId = args[1]
	patient.DOB = args[3]
	patient.Sex = args[4]
	patient.Race = args[5]
	patient.PatientReporting = key

	patientAsBytes, _ = json.Marshal(patient)
	err = APIstub.PutState(args[0], patientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to edit PatientGeneralInformation ", args[0]))
	}

	return shim.Success(nil)
}

/*
	Query Patient General Information
	arguments: [id]
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
	Query Patient General Information History
	arguments: [id]
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
	Record Patient Diabetes Indications
	arguments: [id, patientID, reporterID, HbA1c, HbA1cDate]
*/
func (s *SmartContract) recordPatientDiabetesIndications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	status, err := s.getPatientReportingStatus(APIstub, args)
	if err != nil {
		return shim.Error(err.Error())
	}
	if !strings.HasPrefix(status, "Active") {
		return shim.Error(fmt.Sprintf("Transaction is refused"))
	}
	patientID := args[1]
	reporterID := args[2]
	values := []string{}
	values = append(values, patientID)
	values = append(values, reporterID)
	key := strings.Join(values, "")

	var patient = PatientDiabetesIndications{
		ObjectType:       "PatientDiabetesIndications",
		Id:               args[0],
		AccountId:        args[1],
		HbA1c:            args[3],
		HbA1cDate:        args[4],
		PatientReporting: key,
	}

	patientAsBytes, _ := json.Marshal(patient)
	err = APIstub.PutState(args[0], patientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record patient diabetes indications: %s", args[0]))
	}

	return shim.Success(nil)
}

/*
	Edit Patient Diabetes Indications
	arguments: [id, patientID, reporterID, HbA1c, HbA1cDate]
*/
func (s *SmartContract) editPatientDiabetesIndications(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	if len(args) != 5 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	status, err := s.getPatientReportingStatus(APIstub, args)
	if err != nil {
		return shim.Error(err.Error())
	}
	if !strings.HasPrefix(status, "Active") {
		return shim.Error(fmt.Sprintf("Transaction is refused"))
	}

	patientAsBytes, _ := APIstub.GetState(args[0])
	if patientAsBytes == nil {
		return shim.Error("Could not locate PatientDiabetesIndications")
	}
	patient := PatientDiabetesIndications{}

	json.Unmarshal(patientAsBytes, &patient)

	patientID := args[1]
	reporterID := args[2]
	values := []string{}
	values = append(values, patientID)
	values = append(values, reporterID)
	key := strings.Join(values, "")

	patient.ObjectType = "PatientDiabetesIndications"
	patient.Id = args[0]
	patient.AccountId = args[1]
	patient.HbA1c = args[3]
	patient.HbA1cDate = args[4]
	patient.PatientReporting = key

	patientAsBytes, _ = json.Marshal(patient)
	err = APIstub.PutState(args[0], patientAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to edit PatientGeneralInformation ", args[0]))
	}

	return shim.Success(nil)
}

/*
	Query Patient Diabetes Indications
	arguments: [id]
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
	Query Patient Diabetes Indications History
	arguments: [id]
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

// Patient Actions
/*
	Query Patient
	arguments: [patientID]
*/
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

// Patient Reporting Action
func (s *SmartContract) getPatientReportingStatus(APIstub shim.ChaincodeStubInterface, args []string) (string, error) {

	patientID := args[1]
	reporterID := args[2]

	values := []string{}
	values = append(values, patientID)
	values = append(values, reporterID)
	key := strings.Join(values, "")
	PatientReportingAsBytes, err := APIstub.GetState(key)
	if err != nil {
		return "", errors.New("Could not locate PatientReportingAsBytes")
	}
	var PatientReporting PatientReporting
	err = json.Unmarshal(PatientReportingAsBytes, &PatientReporting)
	role := PatientReporting.Status
	return role, nil
}

/*
	Record Patient Reporting
	arguments: [PatientId, ReporterID, Status]
*/
func (s *SmartContract) recordPatientReporting(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var patientReporting = PatientReporting{
		ObjectType: "PatientReporting",
		PatientId:  args[0],
		ReporterId: args[1],
		Status:     args[2],
	}
	values := []string{}
	values = append(values, args[0])
	values = append(values, args[1])
	key := strings.Join(values, "")

	patientReportingAsBytes, _ := json.Marshal(patientReporting)
	err := APIstub.PutState(key, patientReportingAsBytes)
	if err != nil {
		return shim.Error(fmt.Sprintf("Failed to record patient reporting: %s", key))
	}

	return shim.Success(nil)
}

/*
	Query Patient Reporting Statuses
	arguments: [PatientId]
*/
func (s *SmartContract) getPatientReportingAccess(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	patientID := args[0]
	queryString := fmt.Sprintf("{\"selector\": {\"docType\": {\"$eq\": \"PatientReporting\"},\"patient\": { \"$eq\": \"%s\" }}}", patientID)

	queryResults, err := getQueryResultForQueryString(APIstub, queryString)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(queryResults)
}

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

func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
