auth:
	go run ./cmd/auth/main.go -envVars

up:
    docker-compose up --build -d

test:
    go test ./...