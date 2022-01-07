run:
	docker-compose run --rm kvstore go run main.go

test:
	docker-compose run --rm kvstore go test -v ./...
