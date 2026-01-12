# First stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum /

COPY vendor ./vendor

COPY --exclude=frontend . .

# Setting CGO_ENABLED=0 bundles all the libraries statically when building the executable
RUN GOOS=linux CGO_ENABLED=0 go build -mod vendor -ldflags="-w -s" -o gochat ./cmd/gochat

# Build the frontend
FROM node:25-alpine AS frontend
ARG VITE_API_BASE_URL

ENV VITE_API_BASE_URL=$VITE_API_BASE_URL

WORKDIR /app

COPY ./frontend/package.json .
COPY ./frontend/package-lock.json .
RUN npm ci
COPY ./frontend/. .
RUN npm run build


# Production stage
FROM alpine
#gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/gochat /

COPY --from=frontend /app/dist /static

#RUN ls

#RUN ls /static

CMD ["/gochat"]