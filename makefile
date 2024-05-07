cs_name = content-service

build-cs:
	go build -o publish/$(cs_name) ./src/$(cs_name)

run-cs:
	go run ./src/content-service

docker-cs:
	docker build -f .\docker\$(cs_name).dockerfile .

compose-cs:
	docker-compose -f .\docker\$(cs_name).docker-compose.yml up --build

grpc-gen-cs:
	protoc --go_out=src --go_opt=paths=import --go-grpc_out=src --go-grpc_opt=paths=import src/contracts/media-content/*.proto