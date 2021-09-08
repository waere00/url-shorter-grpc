# Тестовый проект URL shorter

Сервер и клиент запускаются отдельно друг от друга

Для работы предлагаются команды:
```
make start-back    ## запустит сервер
make start-client  ## запустит клиент
make stop          ## остановит сервер
```
Либо:
```
(from src dir)$ docker-compose up -d database
(from src dir)$ cd server
go run server.go

(from src dir)$ cd client
go run client.go
```