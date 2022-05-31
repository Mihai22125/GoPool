FROM golang:alpine

WORKDIR /GoPool

ADD . .

RUN go mod download

RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate

ENTRYPOINT go run ./cmd/web/*