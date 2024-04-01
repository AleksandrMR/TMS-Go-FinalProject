
do-migration:
	go run ./cmd/migrator --storage-path=./storage/hashService.db --migrations-path=./migrations

run-local:
	go run ./cmd/hashService/main.go --config=./config/local.yaml
