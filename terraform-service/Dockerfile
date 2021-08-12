FROM golang:1.16-alpine AS builder

WORKDIR /app

ENV AWS_ACCESS_KEY_ID="<provide aws access key id>"
ENV AWS_SECRET_ACCESS_KEY="<provide aws secret access key>"
ENV AWS_SESSION_TOKEN="<provide aws session token>"

ENV GIT_USERNAME="<provide git username>"
ENV GIT_PASSWORD="<provide git access token>"
ENV GIT_REPOSITORY="<provide git repo, ex: 'github.com/clarshad/golang.git'>"

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

COPY --from=builder /app . 

EXPOSE 8080

CMD ["./main"]
