# Definitions
BINARY_NAME             := ddos
COMPOSE_FILE_NAME       := docker-compose.yml

.PHONY: migration create up down

MAKEFLAGS += --silent

build:
	 GOARCH=arm64 GOOS=linux go build -o bin/${BINARY_NAME} -v main.go

composeUp:
	docker-compose -f $(COMPOSE_FILE_NAME) up -d

composeUpFresh:
	docker-compose -f $(COMPOSE_FILE_NAME) up -d --build go-server --build go-client

composeRestart:
	docker-compose -f $(COMPOSE_FILE_NAME) restart

composeDown:
	docker-compose down

clean:
	rm ${BINARY_NAME}

run:
	./${BINARY_NAME}
