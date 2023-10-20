package main

import (
	"fmt"

	"github.com/costrouc/go-jupyterhub-api/api"
	"github.com/costrouc/go-jupyterhub-api/utils"
)

func main() {
	client, err := api.CreateClient(&api.ClientConfig{ApiToken: "d92fb08258574bc0b27aed7b6883e370"})
	if err != nil {
		panic(err.Error())
	}

	_, err = client.GetVersion()
	if err != nil {
		panic(err.Error())
	}

	_, err = client.GetInfo()
	if err != nil {
		panic(err.Error())
	}

	_, err = client.GetCurrentUser()
	if err != nil {
		panic(err.Error())
	}

	_, err = client.ListUsers(&api.ListUsersParams{Limit: 10})
	if err != nil {
		panic(err.Error())
	}

	randomUsername1 := fmt.Sprintf("test-%s", utils.RandSeq(16))

	_, err = client.CreateUsers(&api.CreateUsersBody{Admin: true, Usernames: []string{randomUsername1}})
	if err != nil {
		panic(err.Error())
	}
	_, err = client.GetUser(randomUsername1)
	if err != nil {
		panic(err.Error())
	}

	randomUsername2 := fmt.Sprintf("test-%s", utils.RandSeq(16))
	randomUsername3 := fmt.Sprintf("test-%s", utils.RandSeq(16))
	_, err = client.CreateUser(randomUsername2)
	if err != nil {
		panic(err.Error())
	}
	_, err = client.UpdateUser(randomUsername2, &api.UpdateUserBody{Name: randomUsername3})
	if err != nil {
		panic(err.Error())
	}

	err = client.DeleteUser(randomUsername1)
	if err != nil {
		panic(err.Error())
	}
	err = client.DeleteUser(randomUsername3)
	if err != nil {
		panic(err.Error())
	}

	client.NotifyUserActivity("username", &api.UserActivityBody{LastActivity: "2023-10-10 01:10:20"})
	client.ListUserTokens("username")
	client.CreateUserToken("username", &api.CreateUserTokenBody{ExpiresIn: 100, Note: "A note"})
	client.ListGroups(&api.ListGroupsParams{})
	client.GetProxyTable(&api.GetProxyTableParams{})
}
