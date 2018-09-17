FROM golang:1.11-alpine AS build-env

MAINTAINER TFG Co <backend@tfgco.com>

RUN apk update && apk add git
RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p /go/src/github.com/topfreegames/pitaya-admin
ADD . /go/src/github.com/topfreegames/pitaya-admin

WORKDIR /go/src/github.com/topfreegames/pitaya-admin
RUN dep ensure
RUN go build -o main .

FROM alpine:3.8

RUN apk update && apk add ca-certificates

WORKDIR /app
COPY --from=build-env /go/src/github.com/topfreegames/pitaya-admin/main /app
COPY --from=build-env /go/src/github.com/topfreegames/pitaya-admin/config /app/config

CMD ["./main", "start"]
