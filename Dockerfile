FROM golang:1.21

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

COPY . .
COPY migrations/00001_init.up.sql .
COPY migrations/00001_init.down.sql .

RUN go build -o main ./cmd

ENV DB_CONNECTION_STRING="postgresql://bestuser:bestuser@host.docker.internal:4999/?sslmode=disable"
ENV HOST=":8080"
ENV SALT="my salt"

RUN apt-get update && apt-get install -y postgresql-client

RUN migrate -path /app/migrations -database $DB_CONNECTION_STRING up

CMD ["./main"]










