FROM golang:1.18.1-bullseye
WORKDIR /app
COPY . .
EXPOSE 8080
CMD go run main.go
