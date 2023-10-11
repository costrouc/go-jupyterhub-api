package main

import (
	"fmt"

	"github.com/costrouc/go-jupyterhub-api/api"
	"github.com/costrouc/go-jupyterhub-api/utils"
)

func main() {
	client, err := api.CreateClient(&api.ClientConfig{Token: "397d165abaab49c6b5dea2bbad19712e"})
	if err != nil {
		panic(err.Error())
	}
	client.GetVersion()
	_, err = client.GetInfo()
	if err != nil {
		panic(err.Error())
	}
	client.GetCurrentUser()
	client.ListUsers(&api.ListUsersParams{Limit: 10})

	randomUsername1 := fmt.Sprintf("test-%s", utils.RandSeq(16))
	client.CreateUsers(&api.CreateUsersBody{Admin: true, Usernames: []string{randomUsername1}})
	client.GetUser(randomUsername1)

	randomUsername2 := fmt.Sprintf("test-%s", utils.RandSeq(16))
	randomUsername3 := fmt.Sprintf("test-%s", utils.RandSeq(16))
	client.CreateUser(randomUsername2)
	client.UpdateUser(randomUsername2, &api.UpdateUserBody{Name: randomUsername3})

	client.DeleteUser(randomUsername1)
	client.DeleteUser(randomUsername3)

	client.NotifyUserActivity("username", &api.UserActivityBody{LastActivity: "2023-10-10 01:10:20"})
	client.ListUserTokens("username")
	client.CreateUserToken("username", &api.CreateUserTokenBody{ExpiresIn: 100, Note: "A note"})
	client.ListGroups(&api.ListGroupsParams{})
	client.GetProxyTable(&api.GetProxyTableParams{})
}
