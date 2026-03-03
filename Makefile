include .env
export


service_run:
	go run cmd/todo_list/main.go


migrate_up:
	migrate -path ./migrations -database ${STORAGE_PATH} up

migrate_down:
	migrate -path ./migrations -database ${STORAGE_PATH}  down