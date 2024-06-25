COMPOSE_BASE := docker compose -f docker-compose.yml

app-build:
	GOARCH=amd64 GOOS=linux go build -tags musl -o backend-engineer-takehome

build:
	docker build -t shiva5128/backend-engineer-takehome:latest .

start:
	$(COMPOSE_BASE)	up -d;

restart:
	$(COMPOSE_BASE)	restart;

stop:
	$(COMPOSE_BASE)	down;

start-server:
	go run main.go

test:
	go test ./...
