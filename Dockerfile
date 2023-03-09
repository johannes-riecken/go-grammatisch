FROM golang:1.19-bullseye
WORKDIR /app
COPY . .
EXPOSE 8080
CMD go run main.go
