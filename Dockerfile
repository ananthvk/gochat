# First stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum /

COPY vendor ./vendor

COPY . .

# Setting CGO_ENABLED=0 bundles all the libraries statically when building the executable
RUN GOOS=linux CGO_ENABLED=0 go build -mod vendor -ldflags="-w -s" -o gochat ./cmd/gochat

# Second stage
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/gochat /

CMD ["/gochat"]