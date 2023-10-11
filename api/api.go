package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func CreateClient(config *ClientConfig) (*ClientConfig, error) {
	clientConfig := ClientConfig{
		Host:     "localhost:8000",
		Prefix:   "/",
		Protocol: "http",
		Token:    "",
		Username: "",
		Password: "",
	}

	if config.Host != "" {
		clientConfig.Host = config.Host
	}

	if config.Prefix != "" {
		clientConfig.Prefix = config.Prefix
	}

	if config.Protocol != "" {
		clientConfig.Protocol = config.Protocol
	}

	if config.Token != "" {
		clientConfig.Token = config.Token
	} else {
		clientConfig.Token, _ = os.LookupEnv("JUPYTERHUB_TOKEN")
	}

	if config.Username != "" {
		clientConfig.Username = config.Username
	} else {
		username, ok := os.LookupEnv("JUPYTERHUB_USERNAME")
		if !ok && clientConfig.Token == "" {
			return nil, errors.New("environment variable JUPYTERHUB_TOKEN and JUPYTERHUB_USERNAME not defined")
		}
		clientConfig.Username = username
	}

	if config.Password != "" {
		clientConfig.Password = config.Password
	} else {
		password, ok := os.LookupEnv("JUPYTERHUB_PASSWORD")
		if !ok && clientConfig.Token == "" {
			return nil, errors.New("environment variable JUPYTERHUB_TOKEN and JUPYTERHUB_PASSWORD not defined")
		}
		clientConfig.Password = password
	}

	return &clientConfig, nil
}

func (c *ClientConfig) BaseURL() string {
	return fmt.Sprintf("%s://%s%s", c.Protocol, c.Host, c.Prefix)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func (c *ClientConfig) setAuthHeader(request *http.Request) {
	if c.Token != "" {
		request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	} else {
		request.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth(c.Username, c.Password)))
	}
}

func (c *ClientConfig) Request(method string, path string, requestBody []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL(), path)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	c.setAuthHeader(req)
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
		return nil, err
	}
	return body, nil
}

func (c *ClientConfig) GetInfo() (*InfoResponse, error) {
	data, err := c.Request(http.MethodGet, "hub/api/info", nil)
	if err != nil {
		return nil, err
	}

	var result InfoResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetVersion() (*VersionResponse, error) {
	data, err := c.Request(http.MethodGet, "hub/api/", nil)
	if err != nil {
		return nil, err
	}

	var result VersionResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetCurrentUser() (*CurrentUserResponse, error) {
	data, err := c.Request(http.MethodGet, "hub/api/user", nil)
	if err != nil {
		return nil, err
	}

	var result CurrentUserResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) ListUsers(options *ListUsersParams) (*ListUsersResponse, error) {
	url := "hub/api/users"
	if options != nil {
		url = fmt.Sprintf("%s?%s", url, options.Encode())
	}

	data, err := c.Request(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var result ListUsersResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) CreateUsers(options *CreateUsersBody) (*ListUsersResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	data, err := c.Request(http.MethodPost, "hub/api/users", body)
	if err != nil {
		return nil, err
	}

	var result ListUsersResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetUser(username string) (*GetUserResponse, error) {
	data, err := c.Request(http.MethodGet, fmt.Sprintf("hub/api/users/%s", username), nil)
	if err != nil {
		return nil, err
	}

	var result GetUserResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) CreateUser(username string) (*CreateUserResponse, error) {
	data, err := c.Request(http.MethodPost, fmt.Sprintf("hub/api/users/%s", username), nil)
	if err != nil {
		return nil, err
	}

	var result CreateUserResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) DeleteUser(username string) error {
	_, err := c.Request(http.MethodDelete, fmt.Sprintf("hub/api/users/%s", username), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) UpdateUser(username string, options *UpdateUserBody) (*UpdateUserResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}

	data, err := c.Request(http.MethodPatch, fmt.Sprintf("hub/api/users/%s", username), body)
	if err != nil {
		return nil, err
	}

	var result UpdateUserResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) NotifyUserActivity(username string, options *UserActivityBody) error {
	body, err := json.Marshal(options)
	if err != nil {
		return err
	}

	_, err = c.Request(http.MethodPost, fmt.Sprintf("hub/api/users/%s/activity", username), body)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) StartUserServer(username string, options interface{}) error {
	body, err := json.Marshal(options)
	if err != nil {
		return err
	}

	_, err = c.Request(http.MethodPost, fmt.Sprintf("hub/api/users/%s/server", username), body)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) StopUserServer(username string) error {
	_, err := c.Request(http.MethodDelete, fmt.Sprintf("hub/api/users/%s/server", username), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) StartUserNamedServer(username string, serverName string, options interface{}) error {
	body, err := json.Marshal(options)
	if err != nil {
		return err
	}

	_, err = c.Request(http.MethodPost, fmt.Sprintf("hub/api/users/%s/servers/%s", username, serverName), body)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) StopUserNamedServer(username string, serverName string) error {
	_, err := c.Request(http.MethodDelete, fmt.Sprintf("hub/api/users/%s/servers/%s", username, serverName), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) ListUserTokens(username string) (*ListTokenResponse, error) {
	data, err := c.Request(http.MethodGet, fmt.Sprintf("hub/api/users/%s/tokens", username), nil)
	if err != nil {
		return nil, err
	}

	var result ListTokenResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) CreateUserToken(username string, options *CreateUserTokenBody) (*CreateUserTokenResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	data, err := c.Request(http.MethodPost, fmt.Sprintf("hub/api/users/%s/tokens", username), body)
	if err != nil {
		return nil, err
	}

	var result CreateUserTokenResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetUserToken(username string, tokenId string) (*GetUserTokenResponse, error) {
	data, err := c.Request(http.MethodGet, fmt.Sprintf("hub/api/users/%s/tokens/%s", username, tokenId), nil)
	if err != nil {
		return nil, err
	}

	var result GetUserTokenResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) DeleteUserToken(username string, tokenId string) error {
	_, err := c.Request(http.MethodDelete, fmt.Sprintf("hub/api/users/%s/tokens/%s", username, tokenId), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) ListGroups(options *ListGroupsParams) (*ListGroupsResponse, error) {
	url := "hub/api/groups"
	if options != nil {
		url = fmt.Sprintf("%s?%s", url, options.Encode())
	}

	data, err := c.Request(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var result ListGroupsResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetGroup(groupname string) (*GetGroupResponse, error) {
	data, err := c.Request(http.MethodGet, fmt.Sprintf("hub/api/groups/%s", groupname), nil)
	if err != nil {
		return nil, err
	}

	var result GetGroupResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) CreateGroup(groupname string) (*CreateGroupResponse, error) {
	data, err := c.Request(http.MethodPost, fmt.Sprintf("hub/api/groups/%s", groupname), nil)
	if err != nil {
		return nil, err
	}

	var result CreateGroupResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) DeleteGroup(groupname string) error {
	_, err := c.Request(http.MethodDelete, fmt.Sprintf("hub/api/groups/%s", groupname), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) AddGroupUsers(groupname string, options *AddGroupUsersBody) (*AddGroupUsersResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	data, err := c.Request(http.MethodPost, fmt.Sprintf("hub/api/groups/%s/users", groupname), body)
	if err != nil {
		return nil, err
	}

	var result AddGroupUsersResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) RemoveGroupUsers(groupname string, options *RemoveGroupUsersBody) error {
	body, err := json.Marshal(options)
	if err != nil {
		return err
	}
	_, err = c.Request(http.MethodDelete, fmt.Sprintf("hub/api/groups/%s/users", groupname), body)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) SetGroupProperties(groupname string, properties interface{}) error {
	body, err := json.Marshal(properties)
	if err != nil {
		return err
	}
	_, err = c.Request(http.MethodPut, fmt.Sprintf("hub/api/groups/%s/properties", groupname), body)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) ListServices() (*ListServicesResponse, error) {
	data, err := c.Request(http.MethodGet, "hub/api/services", nil)
	if err != nil {
		return nil, err
	}

	var result ListServicesResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetService(servicename string) (*GetServiceResponse, error) {
	data, err := c.Request(http.MethodGet, fmt.Sprintf("hub/api/services/%s", servicename), nil)
	if err != nil {
		return nil, err
	}

	var result GetServiceResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) GetProxyTable(options *GetProxyTableParams) (*GetProxyTableResponse, error) {
	url := "hub/api/proxy"
	if options != nil {
		url = fmt.Sprintf("%s?%s", url, options.Encode())
	}

	data, err := c.Request(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	var result GetProxyTableResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) ForceProxySync() error {
	_, err := c.Request(http.MethodPost, "hub/api/proxy", nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) NotifyNewProxy(options *NotifyNewProxyBody) error {
	body, err := json.Marshal(options)
	if err != nil {
		return err
	}
	_, err = c.Request(http.MethodPost, "hub/api/proxy", body)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) NewAPIToken(options *NewTokenBody) (*NewTokenResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	data, err := c.Request(http.MethodPost, "hub/api/authorizations/token", body)
	if err != nil {
		return nil, err
	}
	var result NewTokenResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) ValidateToken(token string) error {
	_, err := c.Request(http.MethodGet, fmt.Sprintf("hub/api/authorizations/token/%s", token), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) GetOAuth2Endpoint(options *GetOAuth2EndpointParams) string {
	return fmt.Sprintf("%s/oauth2/authorize?%s", c.BaseURL(), options.Encode())
}

func (c *ClientConfig) GetOAuth2Token(options *GetOAuth2TokenBody) (*GetOAuth2TokenResponse, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	data, err := c.Request(http.MethodPost, "hub/api/oauth2/token", body)
	if err != nil {
		return nil, err
	}
	var result GetOAuth2TokenResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *ClientConfig) Shutdown(options *ShutdownBody) error {
	body, err := json.Marshal(options)
	if err != nil {
		return err
	}
	_, err = c.Request(http.MethodPost, "hub/api/shutdown", body)
	if err != nil {
		return err
	}
	return nil
}
