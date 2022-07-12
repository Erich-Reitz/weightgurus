package WeightGurus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type loginData struct {
	email    string
	password string
	web      string
}

type WeightGuruOperation struct {
	bmi             float64
	bodyFat         float64
	entryTimestamp  string
	muscleMass      float64
	operationType   string
	serverTimestamp string
	source          string
	water           float64
	Weight          float64
}

func AddBearerTokenToRequest(req *http.Request, bearerToken string) *http.Request {
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	return req
}

func DoRequestReturnBody(req *http.Request) []byte {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return body
}

func CreateNewGetRequest(endpointUrl string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, endpointUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	return req
}

func CreateNewPostRequest(endpointUrl, contentType string, body io.Reader) *http.Request {
	req, err := http.NewRequest(http.MethodPost, endpointUrl, body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", contentType)
	return req
}

func Test(email, password string) {
	params := weightHistoryParams{
		bearerToken: login(email, password),
		startDate:   "",
	}
	entries := getWeightGurusEntries(params)
	for _, entry := range entries {
		fmt.Println(entry)
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

type weightHistoryParams struct {
	bearerToken string
	startDate   string
}

func getEndpointUrl(params weightHistoryParams) string {
	if params.startDate == "" {
		return "https://api.weightgurus.com/v3/operation/?"
	}

	return "https://api.weightgurus.com/v3/operation/?" + params.startDate
}

func prepareWeightHistoryRequest(params weightHistoryParams) *http.Request {
	req := CreateNewGetRequest(getEndpointUrl(params))
	req = AddBearerTokenToRequest(req, params.bearerToken)
	return req
}

func convertResponseInterfaceToWeightGuruOperation(responseInterface interface{}) WeightGuruOperation {
	var weightGuruOperation WeightGuruOperation
	operationMap := responseInterface.(map[string]interface{})
	if operationMap["bmi"] == nil {
		weightGuruOperation.bmi = 0
	} else {
		weightGuruOperation.bmi = convertWeightGuruNumToFloat(operationMap["bmi"].(float64))
	}
	if operationMap["bodyFat"] == nil {
		weightGuruOperation.bodyFat = 0
	} else {
		weightGuruOperation.bodyFat = convertWeightGuruNumToFloat(operationMap["bodyFat"].(float64))
	}
	if operationMap["entryTimestamp"] == nil {
		weightGuruOperation.entryTimestamp = ""
	} else {
		weightGuruOperation.entryTimestamp = operationMap["entryTimestamp"].(string)
	}
	if operationMap["muscleMass"] == nil {
		weightGuruOperation.muscleMass = 0
	} else {
		weightGuruOperation.muscleMass = convertWeightGuruNumToFloat(operationMap["muscleMass"].(float64))
	}
	if operationMap["operationType"] == nil {
		weightGuruOperation.operationType = ""
	} else {
		weightGuruOperation.operationType = operationMap["operationType"].(string)
	}
	if operationMap["serverTimestamp"] == nil {
		weightGuruOperation.serverTimestamp = ""
	} else {
		weightGuruOperation.serverTimestamp = operationMap["serverTimestamp"].(string)
	}
	if operationMap["source"] == nil {
		weightGuruOperation.source = ""
	} else {
		weightGuruOperation.source = operationMap["source"].(string)
	}
	if operationMap["water"] == nil {
		weightGuruOperation.water = 0
	} else {
		weightGuruOperation.water = convertWeightGuruNumToFloat(operationMap["water"].(float64))
	}
	if operationMap["weight"] == nil {
		weightGuruOperation.Weight = 0
	} else {
		weightGuruOperation.Weight = convertWeightGuruNumToFloat(operationMap["weight"].(float64))
	}

	return weightGuruOperation
}

func removeDeletedOperation(operation WeightGuruOperation, weightHistory []WeightGuruOperation) []WeightGuruOperation {
	for i, v := range weightHistory {
		if v.entryTimestamp == operation.entryTimestamp {
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
	deleteOperations := make([]WeightGuruOperation, 0)

	for _, operation := range operations {
		entry := convertResponseInterfaceToWeightGuruOperation(operation)
		if entry.operationType == "delete" {
			deleteOperations = append(deleteOperations, entry)
		} else {
			weightGuruEntries = append(weightGuruEntries, entry)
		}
	}

	for _, deleteOperation := range deleteOperations {
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
