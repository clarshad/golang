# Golang
## Terraform as a Service
Implements service to trigger terraform operations via HTTP requests

Run below curl commands to apply and destroy terraform configuration from a specific terraform version. 
To get the job status, query with request_id captured while running terraform configuration (POST request).

```
curl -X POST -d '{"version": "0.12.15"}' localhost:8080/apply
curl -X POST -d '{"version": "0.12.15"}' localhost:8080/destroy
curl -X GET localhost:8080/job/{request_id}
```