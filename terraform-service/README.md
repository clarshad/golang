# Terraform as a Service

Implements service to trigger terraform operations via HTTP requests

### Getting Started

#### Prerequisite

Set environment variables to access AWS resources and access git repository where terraform configurations are located

```
export AWS_ACCESS_KEY_ID="<provide aws access key id>"
export AWS_SECRET_ACCESS_KEY="<provide aws secret access key>"
export AWS_SESSION_TOKEN="<provide aws session token>"
export git_username = "<provide git username>"
export git_password = "<provide git access token>"
export git_repo = "<provide git repo, ex: 'github.com/clarshad/golang.git'>"
```

#### Run Locally

After cloning this repository, change directory to terraform-service `cd terraform-service` and run `go build`. This should create `terraform-service` binary, run the file `./terraform-service` to start API server

#### Run as Container

After cloning this repository, change directory to terraform-service `cd terraform-service` and edit `Dockerfile`. Update environment variables as mentioned in prerequisites

Run below commands to build and run the docker image
- `docker build -t terraform-service-image:1.0`
- `docker run -d -p 8080:8080 --name terraform-service terraform-service-image:1.0`

### Examples

Run below curl commands to test the functionality.

```
curl -X POST -d '{"version": "0.12.15", "path": "path/to/terraform/config/file/in/repo"}' localhost:8080/apply
curl -X POST -d '{"version": "0.12.15", "path": "path/to/terraform/config/file/in/repo"}' localhost:8080/destroy
curl -X GET localhost:8080/job/{request_id}
```

- `"version"` represents terrafom version to be installed before running the configuration
- `"path"` represents the path or directory folder where terraform configuation files are located in the git repository. Could be any git repository that's set as environment variable
- `request_id` is the unique ID generated for each POST request. This ID could then be used to GET status of POST job triggered