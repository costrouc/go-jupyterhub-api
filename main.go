package main

import (
	"fmt"

	"github.com/costrouc/go-jupyterhub-api/api"
	"github.com/costrouc/go-jupyterhub-api/utils"
)

func main() {
	api.GetVersion()
	api.GetInfo()
	api.GetCurrentUser()
	api.ListUsers(&api.JupyterHubListUsersParams{Limit: 10})

	randomUsername1 := fmt.Sprintf("test-%s", utils.RandSeq(16))
	api.CreateUsers(&api.JupyterHubCreateUsersBody{Admin: true, Usernames: []string{randomUsername1}})
	api.GetUser(randomUsername1)

	randomUsername2 := fmt.Sprintf("test-%s", utils.RandSeq(16))
	randomUsername3 := fmt.Sprintf("test-%s", utils.RandSeq(16))
	api.CreateUser(randomUsername2)
	api.UpdateUser(randomUsername2, &api.JupyterHubUpdateUserBody{Name: randomUsername3})

	api.DeleteUser(randomUsername1)
	api.DeleteUser(randomUsername3)

	api.NotifyUserActivity("username", &api.JupyterHubUserActivityBody{LastActivity: "2023-10-10 01:10:20"})
	api.ListUserTokens("username")
	api.CreateUserToken("username", &api.JupyterHubCreateUserTokenBody{ExpiresIn: 100, Note: "A note"})
}
