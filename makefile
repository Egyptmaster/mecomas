cs_name = content-service
scylla_db_data = ../data/scylla-db

build-cs:
	go build -o publish/$(cs_name) ./src/$(cs_name)

run-cs:
	go run ./src/$(cs_name) ./settings/$(cs_name).settings.yml

docker-cs:
	docker build -f .\docker\$(cs_name).dockerfile .

compose-cs:
	docker-compose -f .\docker\$(cs_name).docker-compose.yml up --build

grpc-gen-cs:
	protoc --go_out=src --go_opt=paths=import --go-grpc_out=src --go-grpc_opt=paths=import src/contracts/media-content/*.proto

scylla-db:
	docker-compose -f .\docker\scylla-db.docker-compose.yml up