# Go JupyterHub API

This is a golang module to implement the [JupyterHub API](https://jupyterhub.readthedocs.io/en/stable/reference/rest-api.html).

```go
client, err := CreateClient(&ClientConfig{ApiToken: "usertoken"})
if err != nil {
    fmt.Errorf(err)
}
ctx := context.Background()
data, err := client.GetCurrentUser(ctx)
if err != nil {
	t.Error(err)
}
fmt.Printf("User is %v!", data.Name)
```