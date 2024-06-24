FROM golang:1.22-alpine3.18 AS build

RUN apk add --update --no-cache openssl-dev curl git openssh gcc musl-dev linux-headers util-linux

ENV CGO_ENABLED=0

WORKDIR $GOPATH/src

ADD ./go.mod $GOPATH/src
ADD ./go.sum $GOPATH/src

RUN go mod download

ADD ./ $GOPATH/src

RUN go install keight/cmd/api
FROM alpine:3.18 as keight

RUN apk --no-cache --update add ca-certificates && rm -rf /tmp/* /var/cache/apk/*
ENV APP_ROOT /opt/keight
RUN mkdir -p $APP_ROOT
ENV PATH "$PATH:$APP_ROOT"

WORKDIR $APP_ROOT
RUN adduser -D -u 1111 keight
COPY --from=build /go/bin/ $APP_ROOT/
USER keight
