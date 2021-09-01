.PHONY: start-back
start-back:
	docker-compose up -d database
	go run server/server.go

.PHONY: start-client
start-client:
	go run client/client.go

.PHONY: stop
stop:
	docker-compose down --volumes
