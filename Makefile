test: ### Run unit tests
	go test -v -cover -race ./internal/...

local-run:
	docker-compose up --build