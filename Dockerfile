## Используйте официальный образ Go
#FROM golang:1.21
#
## Установите рабочую директорию в контейнере
#WORKDIR /app
#
## Копируйте go.mod и go.sum в рабочую директорию
#COPY go.mod go.sum ./
#
## Загрузите все зависимости
#RUN go mod download
#
## Копируйте исходный код в рабочую директорию
#COPY . .
#
## Установите инструмент миграции
#RUN GOBIN=/app/bin go get -u github.com/golang-migrate/migrate/v4/cmd/migrate
#
## Добавьте /app/bin в PATH
#ENV PATH="/app/bin:${PATH}"
#
## Соберите приложение
#RUN go build -o main ./cmd
#
## Выполните миграции базы данных
#RUN migrate -source file://migrations -database 'postgres://bestuser:bestuser@localhost:5432/ewallet?sslmode=disable' up
#
## Экспортируйте порт, который будет слушать приложение
#EXPOSE 8080
#
## Запустите приложение
#CMD ["./main"]
# Используйте официальный образ Golang 1.21 для сборки исполняемого файла
# Используйте официальный образ Golang 1.21 для сборки исполняемого файла
FROM golang:1.21 as builder

# Установите рабочую директорию внутри контейнера
WORKDIR /app

# Скопируйте go mod и sum файлы
COPY go.mod go.sum ./

# Загрузите все зависимости. Их можно кэшировать, если go.mod и go.sum не изменяются
RUN go mod download

# Скопируйте исходный код и файлы миграции в рабочую директорию
COPY . .
COPY migrations migrations

# Соберите приложение
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd


# Используйте официальный образ alpine для исполняемого файла
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Скопируйте исполняемый файл из билдера
COPY --from=builder /app/main .

# Установите переменную окружения для строки подключения к базе данных
ENV DB_CONNECTION_STRING="postgres://bestuser:bestuser@localhost:5432/ewallet?sslmode=disable"

# Запустите приложение
CMD ["./main"]




