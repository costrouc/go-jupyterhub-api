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
	api.CreateUser(randomUsername2)

	api.DeleteUser(randomUsername1)
	api.DeleteUser(randomUsername2)
}
