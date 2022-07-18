package WeightGurus

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

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
