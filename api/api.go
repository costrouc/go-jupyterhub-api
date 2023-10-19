package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func CreateClient(config *ClientConfig) (*ClientConfig, error) {
	clientConfig := ClientConfig{
		ApiToken:                 "",
		ServiceName:              "",
		ApiURL:                   "http://localhost:8000/hub/api",
		BaseURL:                  "/",
		ServicePrefix:            "",
		ServiceURL:               "",
		OAuthScopes:              []string{},
		OAuthAccessScopes:        []string{},
		OAuthClientAllowedScopes: []string{},
		ClientId:                 "",
	}

	if config.ApiToken != "" {
		clientConfig.ApiToken = config.ApiToken
	} else {
		apiToken, ok := os.LookupEnv("JUPYTERHUB_API_TOKEN")
		if !ok {
			return nil, errors.New("api token not defined can be set via JUPYTERHUB_API_TOKEN")
		}
		clientConfig.ApiToken = apiToken
	}

	if config.ServiceName != "" {
		clientConfig.ServiceName = config.ServiceName
	} else {
		clientConfig.ServiceName = os.Getenv("JUPYTERHUB_SERVICE_NAME")
	}

	if config.ApiURL != "" {
		clientConfig.ApiURL = config.ApiURL
	} else {
		apiURL, ok := os.LookupEnv("JUPYTERHUB_API_URL")
		if ok {
			clientConfig.ApiURL = apiURL
		}
	}

	if config.BaseURL != "" {
		clientConfig.BaseURL = config.BaseURL
	} else {
		clientConfig.BaseURL = os.Getenv("JUPYTERHUB_BASE_URL")
	}

	if config.ServicePrefix != "" {
		clientConfig.ServicePrefix = config.ServicePrefix
	} else {
		clientConfig.ServicePrefix = os.Getenv("JUPYTERHUB_SERVICE_PREFIX")
	}

	if config.ServiceURL != "" {
		clientConfig.ServiceURL = config.ServiceURL
	} else {
		clientConfig.ServiceURL = os.Getenv("JUPYTERHUB_SERVICE_URL")
	}

	if len(config.OAuthScopes) != 0 {
		clientConfig.OAuthScopes = config.OAuthScopes
	} else {
		OAuthScopesString, ok := os.LookupEnv("JUPYTERHUB_OAUTH_SCOPES")
		if ok {
			err := json.Unmarshal([]byte(OAuthScopesString), &clientConfig.OAuthScopes)
			if err != nil {
				return nil, err
			}
		} else {
			clientConfig.OAuthScopes = []string{}
		}

	}

	if len(config.OAuthAccessScopes) != 0 {
		clientConfig.OAuthAccessScopes = config.OAuthAccessScopes
	} else {
		OAuthAccessScopesString, ok := os.LookupEnv("JUPYTERHUB_OAUTH_ACCESS_SCOPES")
		if ok {
			err := json.Unmarshal([]byte(OAuthAccessScopesString), &clientConfig.OAuthAccessScopes)
			if err != nil {
				return nil, err
			}
		} else {
			clientConfig.OAuthAccessScopes = []string{}
		}

	}

	if len(config.OAuthScopes) != 0 {
		clientConfig.OAuthClientAllowedScopes = config.OAuthClientAllowedScopes
	} else {
		OAuthClientAllowedScopesString, ok := os.LookupEnv("JUPYTERHUB_OAUTH_CLIENT_ALLOWED_SCOPES")
		if ok {
			err := json.Unmarshal([]byte(OAuthClientAllowedScopesString), &clientConfig.OAuthClientAllowedScopes)
			if err != nil {
				return nil, err
			}
		} else {
			clientConfig.OAuthClientAllowedScopes = []string{}
		}

	}

	if config.ClientId != "" {
		clientConfig.ClientId = config.ClientId
	} else {
		clientConfig.ClientId = os.Getenv("JUPYTERHUB_CLIENT_ID")
	}

	return &clientConfig, nil
}

func (c *ClientConfig) Request(method string, path string, requestBody []byte) ([]byte, error) {
	url := fmt.Sprintf("%s/%s", c.ApiURL, path)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiToken))
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
	data, err := c.Request(http.MethodGet, "info", nil)
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
	data, err := c.Request(http.MethodGet, "", nil)
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
	data, err := c.Request(http.MethodGet, "user", nil)
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
	url := "users"
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
	data, err := c.Request(http.MethodPost, "users", body)
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
	data, err := c.Request(http.MethodGet, fmt.Sprintf("users/%s", username), nil)
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
	data, err := c.Request(http.MethodPost, fmt.Sprintf("users/%s", username), nil)
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
	_, err := c.Request(http.MethodDelete, fmt.Sprintf("users/%s", username), nil)
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

	data, err := c.Request(http.MethodPatch, fmt.Sprintf("users/%s", username), body)
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

	_, err = c.Request(http.MethodPost, fmt.Sprintf("users/%s/activity", username), body)
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

	_, err = c.Request(http.MethodPost, fmt.Sprintf("users/%s/server", username), body)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) StopUserServer(username string) error {
	_, err := c.Request(http.MethodDelete, fmt.Sprintf("users/%s/server", username), nil)
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

	_, err = c.Request(http.MethodPost, fmt.Sprintf("users/%s/servers/%s", username, serverName), body)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) StopUserNamedServer(username string, serverName string) error {
	_, err := c.Request(http.MethodDelete, fmt.Sprintf("users/%s/servers/%s", username, serverName), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) ListUserTokens(username string) (*ListTokenResponse, error) {
	data, err := c.Request(http.MethodGet, fmt.Sprintf("users/%s/tokens", username), nil)
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
	data, err := c.Request(http.MethodPost, fmt.Sprintf("users/%s/tokens", username), body)
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
	data, err := c.Request(http.MethodGet, fmt.Sprintf("users/%s/tokens/%s", username, tokenId), nil)
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
	_, err := c.Request(http.MethodDelete, fmt.Sprintf("users/%s/tokens/%s", username, tokenId), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) ListGroups(options *ListGroupsParams) (*ListGroupsResponse, error) {
	url := "groups"
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
	data, err := c.Request(http.MethodGet, fmt.Sprintf("groups/%s", groupname), nil)
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
	data, err := c.Request(http.MethodPost, fmt.Sprintf("groups/%s", groupname), nil)
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
	_, err := c.Request(http.MethodDelete, fmt.Sprintf("groups/%s", groupname), nil)
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
	data, err := c.Request(http.MethodPost, fmt.Sprintf("groups/%s/users", groupname), body)
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
	_, err = c.Request(http.MethodDelete, fmt.Sprintf("groups/%s/users", groupname), body)
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
	_, err = c.Request(http.MethodPut, fmt.Sprintf("groups/%s/properties", groupname), body)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) ListServices() (*ListServicesResponse, error) {
	data, err := c.Request(http.MethodGet, "services", nil)
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
	data, err := c.Request(http.MethodGet, fmt.Sprintf("services/%s", servicename), nil)
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
	url := "proxy"
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
	_, err := c.Request(http.MethodPost, "proxy", nil)
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
	_, err = c.Request(http.MethodPost, "proxy", body)
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
	data, err := c.Request(http.MethodPost, "authorizations/token", body)
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
	_, err := c.Request(http.MethodGet, fmt.Sprintf("authorizations/token/%s", token), nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ClientConfig) GetOAuth2Endpoint(options *GetOAuth2EndpointParams) (string, error) {
	if options.ClientId == "" && c.ClientId != "" {
		options.ClientId = c.ClientId
	} else if options.ClientId == "" && c.ServiceName != "" {
		options.ClientId = fmt.Sprintf("service-%s", c.ServiceName)
	} else if options.ClientId == "" {
		return "", errors.New("ClientId not set via options or environment variable JUPYTERHUB_CLIENT_ID or JUPYTERHUB_SERVICE_NAME")
	}

	if options.ResponseType == "" {
		options.ResponseType = "code"
	}

	if options.State == "" {
		return "", errors.New("state not set for OAuth request and is required")
	}

	if options.RedirectUri == "" {
		return "", errors.New("RedirectUri not set for OAuth request and is required")
	}

	return fmt.Sprintf("%s/oauth2/authorize?%s", c.ApiURL, options.Encode()), nil
}

func (c *ClientConfig) ParseOAuthRequest(r *http.Request, state string) (string, error) {
	query := r.URL.Query()
	if state != query.Get("state") {
		return "", errors.New("state of request did not match expected state")
	}
	return query.Get("code"), nil
}

func (c *ClientConfig) GetOAuth2Token(options *GetOAuth2TokenBody) (*GetOAuth2TokenResponse, error) {
	if options.ClientId == "" && c.ServiceName != "" {
		options.ClientId = c.ServiceName
	} else {
		return nil, errors.New("ClientId not set via options or environment variable JUPYTERHUB_SERVICE_NAME")
	}

	if options.ClientSecret == "" {
		options.ClientSecret = c.ApiToken
	}

	if options.GrantType == "" {
		options.GrantType = "authorization_code"
	}

	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	data, err := c.Request(http.MethodPost, "oauth2/token", body)
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
	_, err = c.Request(http.MethodPost, "shutdown", body)
	if err != nil {
		return err
	}
	return nil
}
