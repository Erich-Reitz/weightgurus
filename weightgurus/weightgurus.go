package WeightGurus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
)

type loginData struct {
	email    string
	password string
	web      string
}

type WeightGuruOperation struct {
	Bmi             float64
	BodyFat         float64
	entryTimestamp  string
	MuscleMass      float64
	operationType   string
	ServerTimestamp string
	source          string
	Water           float64
	Weight          float64
}

type weightHistoryParams struct {
	bearerToken string
	startDate   string
}

func GetNonDeletedEntries(email, password string) []WeightGuruOperation {
	bearerToken := login(email, password)

	params := weightHistoryParams{
		bearerToken: bearerToken,
		startDate:   "",
	}
	weightGuruEntries := getWeightGurusEntries(params)
	return weightGuruEntries
}

func WriteNonDeletedEntriesToFile(email, password, fileName string) {

	weightGuruEntries := GetNonDeletedEntries(email, password)
	jsonData, err := json.Marshal(weightGuruEntries)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func login(email, password string) string {
	encodedLoginData, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
		"web":      "True",
	})
	postBody := bytes.NewBuffer(encodedLoginData)

	req := CreateNewPostRequest("https://api.weightgurus.com/v3/account/login", "application/json", postBody)
	body := DoRequestReturnBody(req)

	var weightGurusResponse map[string]interface{}
	json.Unmarshal(body, &weightGurusResponse)

	return weightGurusResponse["accessToken"].(string)
}

func convertResponseInterfaceToWeightGuruOperation(responseInterface interface{}) WeightGuruOperation {
	var weightGuruOperation WeightGuruOperation
	operationMap := responseInterface.(map[string]interface{})
	if operationMap["bmi"] == nil {
		weightGuruOperation.Bmi = 0
	} else {
		weightGuruOperation.Bmi = convertWeightGuruNumToFloat(operationMap["bmi"].(float64))
	}
	if operationMap["bodyFat"] == nil {
		weightGuruOperation.BodyFat = 0
	} else {
		weightGuruOperation.BodyFat = convertWeightGuruNumToFloat(operationMap["bodyFat"].(float64))
	}
	if operationMap["entryTimestamp"] == nil {
		weightGuruOperation.entryTimestamp = ""
	} else {
		weightGuruOperation.entryTimestamp = operationMap["entryTimestamp"].(string)
	}
	if operationMap["muscleMass"] == nil {
		weightGuruOperation.MuscleMass = 0
	} else {
		weightGuruOperation.MuscleMass = convertWeightGuruNumToFloat(operationMap["muscleMass"].(float64))
	}
	if operationMap["operationType"] == nil {
		weightGuruOperation.operationType = ""
	} else {
		weightGuruOperation.operationType = operationMap["operationType"].(string)
	}
	if operationMap["serverTimestamp"] == nil {
		weightGuruOperation.ServerTimestamp = ""
	} else {
		weightGuruOperation.ServerTimestamp = operationMap["serverTimestamp"].(string)
	}
	if operationMap["source"] == nil {
		weightGuruOperation.source = ""
	} else {
		weightGuruOperation.source = operationMap["source"].(string)
	}
	if operationMap["water"] == nil {
		weightGuruOperation.Water = 0
	} else {
		weightGuruOperation.Water = convertWeightGuruNumToFloat(operationMap["water"].(float64))
	}
	if operationMap["weight"] == nil {
		weightGuruOperation.Weight = 0
	} else {
		weightGuruOperation.Weight = convertWeightGuruNumToFloat(operationMap["weight"].(float64))
	}

	return weightGuruOperation
}

func removeDeletedOperation(deletedOperation WeightGuruOperation, weightHistory []WeightGuruOperation) []WeightGuruOperation {
	for i, v := range weightHistory {
		if v.entryTimestamp == deletedOperation.entryTimestamp {
			weightHistory = append(weightHistory[:i], weightHistory[i+1:]...)
		}
	}
	return weightHistory
}

func getWeightGurusOperations(params weightHistoryParams) []interface{} {
	req := prepareWeightHistoryRequest(params)
	body := DoRequestReturnBody(req)

	var response map[string]interface{}
	json.Unmarshal(body, &response)

	operations := response["operations"].([]interface{})
	return operations
}

func getWeightGurusEntries(params weightHistoryParams) []WeightGuruOperation {

	operations := getWeightGurusOperations(params)
	var weightGuruEntries []WeightGuruOperation

	// delete operations can appear before the operations they are deleting
	operationsToDelete := make([]WeightGuruOperation, 0)

	for _, operation := range operations {
		entry := convertResponseInterfaceToWeightGuruOperation(operation)
		if entry.operationType == "delete" {
			operationsToDelete = append(operationsToDelete, entry)
		} else {
			weightGuruEntries = append(weightGuruEntries, entry)
		}
	}

	for _, deleteOperation := range operationsToDelete {
		weightGuruEntries = removeDeletedOperation(deleteOperation, weightGuruEntries)
	}

	return weightGuruEntries
}

func convertWeightGuruNumToFloat(weightGurusNum float64) float64 {
	number := fmt.Sprintf("%.0f", weightGurusNum)

	if len(number) <= 1 {
		log.Fatal("WeightGurus number is too small, behavior undefined")
	}
	decimalPoint := string(number[len(number)-1:])

	wholeNumber := number[:len(number)-1]

	decimalPointFloat, err := strconv.ParseFloat(decimalPoint, 64)
	if err != nil {
		log.Fatal(err)
	}
	decimalPointFloat = decimalPointFloat / 10
	wholeNumberFloat, err := strconv.ParseFloat(wholeNumber, 64)
	if err != nil {
		log.Fatal(err)
	}
	return wholeNumberFloat + decimalPointFloat
}
