package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

var (
	_client *http.Client
)

type HttpResponseBody struct {
	Success bool        `json:"success"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

type APIError struct {
	ResponseCode int
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API returned code: %d", e.ResponseCode)
}

func SendGetRequest(url string, returnModel interface{}) error {

	if _client == nil {
		_client = http.DefaultClient
	}

	fmt.Printf("#### GET #### url: %s\r\n", url)

	req, _ := http.NewRequest("GET", url, nil)

	r, err := _client.Do(req)

	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return &APIError{ResponseCode: r.StatusCode}
	}
	defer r.Body.Close()

	responseObject := HttpResponseBody{}

	json.NewDecoder(r.Body).Decode(&responseObject)

	if !responseObject.Success {
		return errors.New("api returned an error: " + responseObject.Error)
	}

	b, _ := json.Marshal(responseObject.Data)
	json.Unmarshal(b, returnModel)

	return nil
}

func SendPostRequest(url string, bodyModel interface{}, returnModel interface{}) error {

	if _client == nil {
		_client = http.DefaultClient
	}

	bodyJSON, err := json.Marshal(bodyModel)

	if err != nil {
		return err
	}

	fmt.Printf("#### POST #### url: %s, data: %s\r\n", url, bytes.NewBuffer(bodyJSON))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	r, err := _client.Do(req)

	if err != nil {
		return err
	}

	if r.StatusCode != http.StatusOK {
		return &APIError{ResponseCode: r.StatusCode}
	}

	defer r.Body.Close()

	responseObject := HttpResponseBody{}

	json.NewDecoder(r.Body).Decode(&responseObject)

	if !responseObject.Success {
		fmt.Println(r.Body)
		return errors.New("api returned an error: " + responseObject.Error)
	}

	b, _ := json.Marshal(responseObject.Data)
	err = json.Unmarshal(b, returnModel)

	if err != nil {
		return err
	}

	return nil
}
