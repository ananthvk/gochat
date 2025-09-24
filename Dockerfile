FROM golang:1.25-bookworm

WORKDIR /app

COPY go.mod go.sum /

COPY vendor ./vendor

COPY . .

RUN go build -mod vendor -o gochat ./cmd/gochat

CMD ["./gochat"]