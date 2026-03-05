include .env
export


service_run:
	go run cmd/todo_list/main.go


docker_compose_up:
	docker-compose up -d
docker_compose_down:
	docker-compose down

migrate_up:
	migrate -path ./migrations -database ${STORAGE_PATH} up

migrate_down:
	migrate -path ./migrations -database ${STORAGE_PATH}  down