# portainerUtility

portainerUtility is a small cli-based utility to interact with Portainer API.
It allows for creating and deleting of Portainer Swarm stack based on GIT repository and is aimed to help with creating the stack automatically in a CI/CD pipeline.

## Installation

```bash
go build portainerUtility.go
```

## Usage
To create a new Docker Swarm stack from a GIT repository:
```bash
portainerUtility createStack --portainer-url=<PORTAINER_URL> --endpoint-id=<ENDPOINT_ID> --swarm-id="<SWARM_ID>" --name=<STACK_NAME> --portainer-api-key=<API_KEY> --tls-skip-verify=<true|false> --compose-file "<COMPOSE_FILE.YML>" --repository-url="<GIT_REPOSITORY_URL>" --repository-password="<REPOSITORY_PASSWORD>" --repository-username="<REPOSITORY_USERNAME>" --repository-reference-name="<REPOSITORY_REFERENCE_NAME>" --env=<var1=value1,var2=value2...>
```

Repository reference means e.g.: "refs/heads/main" for a latest commit of a "main" branch.

To delete a stack:
```bash
portainerUtility deleteStack --portainer-url=<PORTAINER_URL> --endpoint-id=<ENDPOINT_ID> --swarm-id="<SWARM_ID>" --name=<STACK_NAME> --portainer-api-key=<API_KEY> --tls-skip-verify=<true|false>
```

## Contributing

Pull requests are welcome. For major changes, please open an issue first
to discuss what you would like to change.

## License

GPL 3.0
