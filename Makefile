build:
	go build -o senketsu src/*
	docker-compose build

run: build
	docker-compose up

down:
	docker-compose down --rmi=local --volumes

dev:
	go run src/*
