get-deps:
	dep ensure
up:
	@docker-compose up -d

down:
	@docker-compose down

init:
	@make get-deps
	@make up



run:
	go run example/main.go
