package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func baseURL() string {
	jupyterhubBaseURL, ok := os.LookupEnv("JUPYTERHUB_BASE_URL")
	if !ok {
		jupyterhubBaseURL = "http://localhost:8000"
	}
	return jupyterhubBaseURL
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func setAuthHeader(request *http.Request) {
	token, ok := os.LookupEnv("JUPYTERHUB_TOKEN")
	if ok {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		return
	}

	username, ok := os.LookupEnv("JUPYTERHUB_USERNAME")
	if !ok {
		panic("environment variable JUPYTERHUB_USERNAME not set")
	}

	password, ok := os.LookupEnv("JUPYTERHUB_USERNAME")
	if !ok {
		panic("environment variable JUPYTERHUB_USERNAME not set")
	}

	request.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth(username, password)))
}

func jupyterHubRequest(method string, path string, requestBody []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", baseURL(), path)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	setAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("response returned status code of %d instead of 200", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(string(body))
	return body, nil
}

func GetInfo() (*JupyterHubInfoResponse, error) {
	data, err := jupyterHubRequest(http.MethodGet, "hub/api/info", nil)
	if err != nil {
		return nil, err
	}

	var result JupyterHubInfoResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetVersion() (*JupyterHubVersionResponse, error) {
	data, err := jupyterHubRequest(http.MethodGet, "hub/api/", nil)
	if err != nil {
		return nil, err
	}

	var result JupyterHubVersionResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetCurrentUser() (*JupyterHubCurrentUserResponse, error) {
	data, err := jupyterHubRequest(http.MethodGet, "hub/api/user", nil)
	if err != nil {
		return nil, err
	}

	var result JupyterHubCurrentUserResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func ListUsers(options *JupyterHubListUsersParams) (*JupyterHubListUsersResponse, error) {
	url := "hub/api/users"
	if options != nil {
		url = fmt.Sprintf("%s?%s", url, options.Encode())
	}

	data, err := jupyterHubRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var result JupyterHubListUsersResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func CreateUsers(options *JupyterHubCreateUsersBody) (*JupyterHubListUsersResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	data, err := jupyterHubRequest(http.MethodPost, "hub/api/users", body)
	if err != nil {
		return nil, err
	}

	var result JupyterHubListUsersResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func GetUser(username string) (*JupyterHubGetUserResponse, error) {
	data, err := jupyterHubRequest(http.MethodGet, fmt.Sprintf("hub/api/users/%s", username), nil)
	if err != nil {
		return nil, err
	}

	var result JupyterHubGetUserResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func CreateUser(username string) (*JupyterHubCreateUserResponse, error) {
	data, err := jupyterHubRequest(http.MethodPost, fmt.Sprintf("hub/api/users/%s", username), nil)
	if err != nil {
		return nil, err
	}

	var result JupyterHubCreateUserResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func DeleteUser(username string) error {
	_, err := jupyterHubRequest(http.MethodDelete, fmt.Sprintf("hub/api/users/%s", username), nil)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUser(username string, options *JupyterHubUpdateUserBody) (*JupyterHubUpdateUserResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	data, err := jupyterHubRequest(http.MethodPatch, fmt.Sprintf("hub/api/users/%s", username), body)
	if err != nil {
		return nil, err
	}

	var result JupyterHubUpdateUserResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
