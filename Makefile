postgres:
	docker run -d -p 5432:5432 --name postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:13.2-alpine

psqlshell:
	docker exec -it postgres psql -U root -d orijinplus

createdb:
	docker exec -it postgres createdb --username=root --owner=root orijinplus

dropdb:
	docker exec -it postgres dropdb --username=root orijinplus

migrateup-local:
	migrate -path sql/migrations -database "postgresql://root:secret@localhost:5432/orijinplus?sslmode=disable" -verbose up

migratedown-local:
	migrate -path sql/migrations -database "postgresql://root:secret@localhost:5432/orijinplus?sslmode=disable" -verbose down

sqlc:
	sqlc generate

graphql:
	go get github.com/99designs/gqlgen
	go run github.com/99designs/gqlgen

dbseed:
	./run-seed.sh

run:
	./run-local.sh

build:
	go build -o main cmd/main.go

.PHONY:
	postgres
	psqlshell
	createdb
	dropdb
	migrateup-local
	migratedown-local
	sqlc
	graphql
	dbseed
	run
	build