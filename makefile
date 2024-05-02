build-cs:
	go build -o publish/content-service ./src/content-service

run-cs:
	go run ./src/content-service

docker-cs:
	docker build -f .\docker\content-service.dockerfile .

compose-cs:
	docker-compose -f .\docker\content-service.docker-compose.yml up --build