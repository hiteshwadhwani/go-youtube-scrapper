FROM golang:latest

LABEL maintainer="Your Name Hitesh Wadhwani"
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go
EXPOSE 8080
CMD ["./main"]