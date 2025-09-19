# образ контейнера Go
FROM golang:1.25.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# копировать весь проект в рабочую директорию /app
COPY . .

# компиляция проекта
RUN go build -o app .

CMD ["./app"]