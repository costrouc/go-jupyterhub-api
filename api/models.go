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

type JupyterHubVersionResponse struct {
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

type JupyterHubInfoResponse struct {
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

type JupyterHubCurrentUserResponse JupyterHubUser

type JupyterHubListUsersParams struct {
	State                 string
	Offset                int
	Limit                 int
	IncludeStoppedServers bool
}

func (r *JupyterHubListUsersParams) Encode() string {
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

type JupyterHubListUsersResponse []JupyterHubUser

type JupyterHubCreateUsersBody struct {
	Usernames []string `json:"usernames"`
	Admin     bool     `json:"admin"`
}

type JupyterHubCreateUsersResponse []JupyterHubUser

type JupyterHubGetUserResponse JupyterHubUser

type JupyterHubCreateUserResponse JupyterHubUser

type JupyterHubUpdateUserBody struct {
	Name  string `json:"name"`
	Admin bool   `json:"admin"`
}

type JupyterHubUpdateUserResponse JupyterHubUser

type JupyterHubUserActivityBody struct {
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

type JupyterHubListTokenResponse []JupyterHubToken

type JupyterHubCreateUserTokenBody struct {
	ExpiresIn int      `json:"expires_in"`
	Note      string   `json:"note"`
	Roles     []string `json:"roles"`
	Scopes    []string `json:"scopes"`
}

type JupyterHubCreateUserTokenResponse JupyterHubToken
