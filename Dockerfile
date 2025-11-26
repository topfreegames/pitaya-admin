FROM golang:1.23-alpine AS build-env

MAINTAINER TFG Co <backend@tfgco.com>

RUN apk update && apk add git gcc musl-dev

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o main .

FROM alpine:3.8

RUN apk update && apk add ca-certificates

WORKDIR /app
COPY --from=build-env /app/main /app
COPY --from=build-env /app/config /app/config

CMD ["./main", "start"]
