# Используем официальный образ Golang
FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .
COPY migrations/00001_init.up.sql .
COPY migrations/00001_init.down.sql .

RUN go build -o main ./cmd

ENV DB_CONNECTION_STRING="postgresql://bestuser:bestuser@host.docker.internal:4999/ewallet?sslmode=disable"
ENV HOST=":8080"

RUN apt-get update && apt-get install -y postgresql-client

RUN migrate -path /app/migrations -database $DB_CONNECTION_STRING up

CMD ["./main"]











# Используйте официальный образ PostgreSQL

#FROM golang:1.21 as builder
#
## Установите рабочую директорию внутри контейнера
#WORKDIR /app
#
## Скопируйте go mod и sum файлы
#COPY go.mod go.sum ./
#
## Загрузите все зависимости. Их можно кэшировать, если go.mod и go.sum не изменяются
#RUN go mod download
#
## Установите инструмент миграции
#RUN go get -u github.com/lib/pq
#RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
#
## Скопируйте исходный код и файлы миграции в рабочую директорию
#COPY . .
#COPY migrations migrations
#
## Соберите приложение
#RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd
#
## Используйте официальный образ alpine для исполняемого файла
#FROM alpine:latest
#
#RUN apk --no-cache add ca-certificates
#
#WORKDIR /root/
#
## Скопируйте исполняемый файл из билдера
#COPY --from=builder /app/main .
#
## Установите переменные окружения для строки подключения к базе данных
#ENV DB_HOST="127.0.0.1"
#ENV DB_PORT="5432"
#ENV DB_USER="bestuser"
#ENV DB_PASSWORD="bestuser"
#ENV DB_NAME="ewallet"
#
#EXPOSE 8080
#
## Запустите приложение
#CMD ["./main"]









