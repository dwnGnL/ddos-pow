FROM golang:1.21

WORKDIR /app

COPY . /app

ENV CONFIG_FILE=config.yaml

RUN go mod tidy

RUN GOARCH=arm64 GOOS=linux go build -o bin/ddos -v main.go

EXPOSE 8740

CMD ["./bin/ddos", "server"]