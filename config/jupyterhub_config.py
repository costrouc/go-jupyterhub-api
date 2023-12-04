c.JupyterHub.authenticator_class = "jupyterhub.auth.DummyAuthenticator"
c.JupyterHub.spawner_class = "jupyterhub.spawner.SimpleLocalProcessSpawner"
c.DummyAuthenticator.password = "password"
c.Authenticator.admin_users = ["username"]

c.JupyterHub.api_tokens = {'usertoken': 'username'}

c.JupyterHub.services = [{
    "name": "my-service",
    "api_token": "servicetoken",
    "admin": True,
}]