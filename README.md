# Golang
## Terraform as a Service
Implements service to trigger terraform operations via HTTP requests

Run below curl commands to apply terraform configuration from a specific terraform version and monitor it

```
curl -X POST -d '{"version": "0.12.15"}' localhost:8080
curl -X GET localhost:8080
```