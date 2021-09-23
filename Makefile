auth:
	go run ./cmd/auth/main.go -envVars
up:
	docker-compose up --build
test:
	go test ./...