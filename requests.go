package weightgurus

import (
	"io"
	"io/ioutil"
	"net/http"
)

func AddBearerTokenToRequest(req *http.Request, bearerToken string) *http.Request {
	req.Header.Set("Authorization", "Bearer "+bearerToken)
	return req
}

func DoRequestReturnBody(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// handle outside to defer than log fatal
	if err != nil {
		return nil, err
	}
	return body, err
}

func CreateNewGetRequest(endpointUrl string) (*http.Request, error) {
	return http.NewRequest(http.MethodGet, endpointUrl, nil)
}

func CreateNewPostRequest(endpointUrl, contentType string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodPost, endpointUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return req, nil
}

func getEndpointUrl(params weightHistoryParams) string {
	if params.startDate == "" {
		return "https://api.weightgurus.com/v3/operation/?"
	}

	return "https://api.weightgurus.com/v3/operation/?" + params.startDate
}

func prepareWeightHistoryRequest(params weightHistoryParams) (*http.Request, error) {
	req, err := CreateNewGetRequest(getEndpointUrl(params))
	if err != nil {
		return nil, err
	}
	req = AddBearerTokenToRequest(req, params.bearerToken)
	return req, nil
}
