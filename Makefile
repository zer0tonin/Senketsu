build:
	go build -o senketsu src/main.go
	docker-compose build

run: build
	docker-compose up senketsu

dev:
	go run src/main.go
