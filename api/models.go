package api

import (
	"fmt"
	"net/url"
	"strconv"
)

const (
	ListUsersStateInactive = "inactive"
	ListUsersStateActive   = "active"
	ListUsersStateReady    = "ready"
)

type ClientConfig struct {
	ApiToken                 string
	ServiceName              string
	ApiURL                   string
	BaseURL                  string
	ServicePrefix            string
	ServiceURL               string
	OAuthScopes              []string
	OAuthAccessScopes        []string
	OAuthClientAllowedScopes []string
	ClientId                 string
}

type VersionResponse struct {
	Version string `json:"version"`
}

type AuthenticatorClass struct {
	Class   string `json:"class"`
	Version string `json:"version"`
}

type SpawnerClass struct {
	Class   string `json:"class"`
	Version string `json:"version"`
}

type InfoResponse struct {
	Version       string             `json:"version"`
	Python        string             `json:"python"`
	SysExecutable string             `json:"sys_executable"`
	Authenticator AuthenticatorClass `json:"authenticator"`
	Spawner       SpawnerClass       `json:"spawner"`
}

type JupyterHubServer struct {
	Name         string `json:"name"`
	Ready        bool   `json:"ready"`
	Stopped      bool   `json:"stopped"`
	Pending      string `json:"pending"`
	Url          string `json:"url"`
	ProgressUrl  string `json:"progress_url"`
	Started      string `json:"started"`
	LastActivity string `json:"last_activity"`
	State        interface{}
	UserOptions  interface{}
}

type JupyterHubUser struct {
	SessionId    string                      `json:"session_id"`
	Scopes       []string                    `json:"scopes"`
	Name         string                      `json:"name"`
	Admin        bool                        `json:"admin"`
	Roles        []string                    `json:"roles"`
	Groups       []string                    `json:"groups"`
	Server       string                      `json:"server"`
	Pending      string                      `json:"pending"`
	LastActivity string                      `json:"last_activitiy"`
	Servers      map[string]JupyterHubServer `json:"servers"`
	AuthState    interface{}
}

type CurrentUserResponse JupyterHubUser

type ListUsersParams struct {
	State                 string
	Offset                int
	Limit                 int
	IncludeStoppedServers bool
}

func (r *ListUsersParams) Encode() string {
	v := url.Values{}
	if r.State == ListUsersStateInactive || r.State == ListUsersStateActive || r.State == ListUsersStateReady {
		v.Set("state", r.State)
	}
	if r.Offset != 0 {
		v.Set("offset", fmt.Sprint(r.Offset))
	}
	if r.Limit != 0 {
		v.Set("limit", fmt.Sprint(r.Limit))
	}
	v.Set("include_stopped_servers", strconv.FormatBool(r.IncludeStoppedServers))
	return v.Encode()
}

type ListUsersResponse []JupyterHubUser

type CreateUsersBody struct {
	Usernames []string `json:"usernames"`
	Admin     bool     `json:"admin"`
}

type CreateUsersResponse []JupyterHubUser

type GetUserResponse JupyterHubUser

type CreateUserResponse JupyterHubUser

type UpdateUserBody struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

type UpdateUserResponse JupyterHubUser

type UserActivityBody struct {
	LastActivity string            `json:"last_activity"`
	Servers      map[string]string `json:"servers"`
}

type JupyterHubToken struct {
	Token        string   `json:"token"`
	Id           string   `json:"id"`
	User         string   `json:"user"`
	Service      string   `json:"service"`
	Roles        []string `json:"roles"`
	Scopes       []string `json:"scopes"`
	Note         string   `json:"note"`
	Created      string   `json:"created"`
	ExpiresAt    string   `json:"expires_at"`
	LastActivity string   `json:"last_activity"`
	SessionId    string   `json:"session_id"`
}

type ListTokenResponse []JupyterHubToken

type CreateUserTokenBody struct {
	ExpiresIn int      `json:"expires_in"`
	Note      string   `json:"note"`
	Roles     []string `json:"roles"`
	Scopes    []string `json:"scopes"`
}

type CreateUserTokenResponse JupyterHubToken

type GetUserTokenResponse JupyterHubToken

type ListGroupsParams struct {
	Offset int
	Limit  int
}

func (r *ListGroupsParams) Encode() string {
	v := url.Values{}
	if r.Offset != 0 {
		v.Set("offset", fmt.Sprint(r.Offset))
	}
	if r.Limit != 0 {
		v.Set("limit", fmt.Sprint(r.Limit))
	}
	return v.Encode()
}

type JupyterHubGroup struct {
	Name       string      `json:"name"`
	Users      []string    `json:"users"`
	Properties interface{} `json:"properties"`
	Roles      []string    `json:"roles"`
}

type ListGroupsResponse []JupyterHubGroup

type GetGroupResponse JupyterHubGroup

type CreateGroupResponse JupyterHubGroup

type AddGroupUsersBody struct {
	Users []string `json:"users"`
}

type AddGroupUsersResponse JupyterHubGroup

type RemoveGroupUsersBody struct {
	Users []string `json:"users"`
}

type JupyterHubService struct {
	Name    string      `json:"name"`
	Admin   bool        `json:"admin"`
	Roles   []string    `json:"roles"`
	Url     string      `json:"url"`
	Prefix  string      `json:"prefix"`
	Pid     int         `json:"pid"`
	Command []string    `json:"command"`
	Info    interface{} `json:"info"`
}

type ListServicesResponse []JupyterHubService

type GetServiceResponse JupyterHubService

type GetProxyTableParams struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (r *GetProxyTableParams) Encode() string {
	v := url.Values{}
	if r.Offset != 0 {
		v.Set("offset", fmt.Sprint(r.Offset))
	}
	if r.Limit != 0 {
		v.Set("limit", fmt.Sprint(r.Limit))
	}
	return v.Encode()
}

type JupyterHubProxyRoute struct {
	RouteSpec string      `json:"routespec"`
	Target    string      `json:"target"`
	Data      interface{} `json:"data"`
}

type GetProxyTableResponse map[string]JupyterHubProxyRoute

type NotifyNewProxyBody struct {
	Ip        string `json:"ip"`
	Port      string `json:"port"`
	Protocol  string `json:"protocol"`
	AuthToken string `json:"auth_token"`
}

type NewTokenBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewTokenResponse struct {
	Token string `json:"token"`
}

type GetOAuth2EndpointParams struct {
	ClientId     string
	ResponseType string
	State        string
	RedirectUri  string
	Scope        string
}

func (p *GetOAuth2EndpointParams) Encode() string {
	v := url.Values{}
	v.Set("client_id", p.ClientId)
	v.Set("response_type", p.ResponseType)
	v.Set("state", p.State)
	v.Set("redirect_uri", p.RedirectUri)
	if p.Scope != "" {
		v.Set("scope", p.Scope)
	}
	return v.Encode()
}

type GetOAuth2TokenBody struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	Code         string `json:"code"`
	RedirectUri  string `json:"redirect_uri"`
}

type GetOAuth2TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

type ShutdownBody struct {
	Proxy   bool `json:"proxy"`
	Servers bool `json:"servers"`
}
