package weightgurus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

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

func GetNonDeletedEntries(email, password string) ([]WeightGuruOperation, error) {
	bearerToken, err := login(email, password)
	if err != nil {
		return nil, err
	}

	params := weightHistoryParams{
		bearerToken: bearerToken,
		startDate:   "",
	}
	weightGuruEntries, err := getWeightGurusEntries(params)
	if err != nil {
		return nil, err
	}
	return weightGuruEntries, nil
}

func WriteNonDeletedEntriesToFile(email, password, fileName string) error {

	weightGuruEntries, err := GetNonDeletedEntries(email, password)
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(weightGuruEntries)
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, jsonData, 0644)

	return err
}

func login(email, password string) (string, error) {
	encodedLoginData, _ := json.Marshal(map[string]string{
		"email":    email,
		"password": password,
		"web":      "True",
	})
	postBody := bytes.NewBuffer(encodedLoginData)

	req, err := http.NewRequest(http.MethodPost, "https://api.weightgurus.com/v3/account/login", postBody)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	if err != nil {
		return "", err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var weightGurusResponse map[string]interface{}
	err = json.Unmarshal(body, &weightGurusResponse)

	return weightGurusResponse["accessToken"].(string), err
}

func convertResponseInterfaceToWeightGuruOperation(responseInterface interface{}) WeightGuruOperation {
	var weightGuruOperation WeightGuruOperation
	operationMap := responseInterface.(map[string]interface{})
	if operationMap["bmi"] != nil {
		weightGuruOperation.Bmi = convertWeightGuruNumToFloat(operationMap["bmi"].(float64))
	}
	if operationMap["bodyFat"] != nil {
		weightGuruOperation.BodyFat = convertWeightGuruNumToFloat(operationMap["bodyFat"].(float64))
	}
	if operationMap["entryTimestamp"] != nil {
		weightGuruOperation.entryTimestamp = operationMap["entryTimestamp"].(string)
	}
	if operationMap["muscleMass"] != nil {
		weightGuruOperation.MuscleMass = convertWeightGuruNumToFloat(operationMap["muscleMass"].(float64))
	}
	if operationMap["operationType"] != nil {
		weightGuruOperation.operationType = operationMap["operationType"].(string)
	}
	if operationMap["serverTimestamp"] != nil {
		weightGuruOperation.ServerTimestamp = operationMap["serverTimestamp"].(string)
	}
	if operationMap["source"] != nil {
		weightGuruOperation.source = operationMap["source"].(string)
	}
	if operationMap["water"] != nil {
		weightGuruOperation.Water = convertWeightGuruNumToFloat(operationMap["water"].(float64))
	}
	if operationMap["weight"] == nil {
		weightGuruOperation.Weight = convertWeightGuruNumToFloat(operationMap["weight"].(float64))
	}

	return weightGuruOperation
}

func removeDeletedOperation(deletedOperation WeightGuruOperation, weightHistory *[]WeightGuruOperation) []WeightGuruOperation {
	for i, v := range *weightHistory {
		if v.entryTimestamp == deletedOperation.entryTimestamp {
			// result is everything but the deleted operation
			if i == len(*weightHistory)-1 {
				return (*weightHistory)[:i]
			}
			*weightHistory = append((*weightHistory)[:i], (*weightHistory)[i+1:]...)
			break
		}
	}
	return *weightHistory
}

func getWeightGurusOperations(params weightHistoryParams) ([]interface{}, error) {
	var endpointUrl string
	if params.startDate == "" {
		endpointUrl = "https://api.weightgurus.com/v3/operation/?"
	} else {
		endpointUrl = "https://api.weightgurus.com/v3/operation/?" + params.startDate
	}

	req, err := http.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+params.bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	// handle outside to defer than log fatal
	if err != nil {
		return nil, err
	}

	var response map[string]interface{}
	json.Unmarshal(body, &response)

	operations := response["operations"].([]interface{})
	return operations, nil
}

func getWeightGurusEntries(params weightHistoryParams) ([]WeightGuruOperation, error) {

	operations, err := getWeightGurusOperations(params)
	if err != nil {
		return nil, err
	}
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
		weightGuruEntries = removeDeletedOperation(deleteOperation, &weightGuruEntries)
	}

	return weightGuruEntries, nil
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
