test:
	go test ./...

build_server:
	go build -o server ./cmd/server

build_client:
	go build -o client ./cmd/client

run_server:
	go build -o server ./cmd/server && ./server

run_client:
	go build -o client ./cmd/client && ./client

docker-up:
	docker-compose up