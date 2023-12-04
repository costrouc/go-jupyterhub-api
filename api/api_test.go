package api

import (
	"context"
	"regexp"
	"testing"
)

func TestJupyterHubStatus(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "usertoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	data, err := client.GetInfo(ctx)
	if err != nil {
		t.Error(err)
	}
	expectedSpawner := "jupyterhub.spawner.SimpleLocalProcessSpawner"
	if data.Spawner.Class != expectedSpawner {
		t.Errorf("Expected spawner %v, got %v", expectedSpawner, data.Spawner.Class)
	}
}

func TestJupyterHubGetVersion(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "usertoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	data, err := client.GetVersion(ctx)
	if err != nil {
		t.Error(err)
	}
	if matched, _ := regexp.Match("[0-9]+\\.[0-9]+\\.[0-9]+", []byte(data.Version)); !matched {
		t.Errorf("Version not in format [0-9]+\\.[0-9]+\\.[0-9]+ , got %v", data.Version)
	}
}

func TestJupyterHubGetCurrentUser(t *testing.T) {
	client, err := CreateClient(&ClientConfig{ApiToken: "usertoken"})
	if err != nil {
		t.Error(err)
	}
	ctx := context.Background()
	data, err := client.GetCurrentUser(ctx)
	if err != nil {
		t.Error(err)
	}
	if data.Name != "username" {
		t.Errorf("Expected authenticated user 'username', got %v", data.Name)
	}
}
