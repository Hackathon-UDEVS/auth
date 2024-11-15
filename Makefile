CURRENT_DIR=$(shell pwd)

create-mig:
 migrate create -ext sql -dir ./internal/migration -seq auth-service

mig-insert:
 migrate create -ext sql -dir migration -seq insert_table


proto-gen:
 ./internal/script/gen-proto.sh ${CURRENT_DIR}

mig-up:
 migrate -database 'postgres://postgres:1111@localhost:5432/cinemauserservice?sslmode=disable' -path migration up

mig-down:
 migrate -database 'postgres://postgres:1111@localhost:5432/cinemauserservice?sslmode=disable' -path migration down

mig-force:
 migrate -database 'postgres://postgres:1111@localhost:5432/cinema-user-service?sslmode=disable' -path migration force 1

SWAGGER := ~/go/bin/swag
SWAGGER_DOCS := docs
SWAGGER_INIT := $(SWAGGER) init -g ./api/router.go -o $(SWAGGER_DOCS)
