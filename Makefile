build:
	go build -o senketsu github.com/zer0tonin/senketsu/src
	docker-compose build senketsu

run: build
	docker-compose up senketsu

down:
	docker-compose down --rmi=local --volumes

dev:
	docker-compose up senketsu-dev

test:
	go test github.com/zer0tonin/senketsu/src/...
