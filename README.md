# Employee service
Веб-сервис на Go, который записывает данные сотрудника (ФИО, телефон, город) в базу PostgreSQL.
Весь проект обернут в Докер и также содержит тесты.

# Стек
- local go version: go1.25.1 windows/amd64
- docker container go version: golang:1.25.1-alpine
- postgresql image version: 15

- github.com/gin-gonic/gin – для REST API
- github.com/jackc/pgx/v5/pgxpool@v5.7.6 – для подключения к Postgres
- github.com/stretchr/testify/assert@v1.11.1 - для тестов
- github.com/joho/godotenv - для загрузки .env

# Чтобы запустить контейнер (в powershell):
cd employee-service
docker-compose up --build -d

# Чтобы запустить тесты (локально через powershell терминал)
## Сперва дождаться пока статус контейнера employees_db станет healthy
go test -v
