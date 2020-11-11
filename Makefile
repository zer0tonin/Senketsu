build:
	go build -o senketsu github.com/zer0tonin/senketsu/src
	docker-compose build

run: build
	docker-compose up

down:
	docker-compose down --rmi=local --volumes

dev:
	go run github.com/zer0tonin/senketsu/src
