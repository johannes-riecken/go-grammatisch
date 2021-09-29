FROM golang:1.17.0-bullseye
WORKDIR /app
COPY . .
EXPOSE 8080
CMD go run cmd/go-grammatisch/main.go
