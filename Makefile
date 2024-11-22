CURRENT_DIR=$(shell pwd)

run:
	go run $(CURRENT_DIR)/cmd/main.go

create-mig:
	migrate create -ext sql -dir ./internal/migration -seq app-service

mig-insert:
	migrate create -ext sql -dir ./internal/migration -seq insert_table

proto-gen:
	./internal/script/gen-proto.sh ${CURRENT_DIR}

mig-up:
	migrate -database 'postgres://postgres:1234@localhost:5432/tender?sslmode=disable' -path ./internal/migration up

mig-down:
	migrate -database 'postgres://postgres:1234@localhost:5432/tender?sslmode=disable' -path ./internal/migration down

mig-force:
	migrate -database 'postgres://postgres:1234@localhost:5432/tender?sslmode=disable' -path ./internal/migration force 1

SWAGGER := ~/go/bin/swag
SWAGGER_DOCS := docs
SWAGGER_INIT := $(SWAGGER) init -g ./api/router.go -o $(SWAGGER_DOCS)