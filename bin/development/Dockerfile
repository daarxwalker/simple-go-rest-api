FROM golang:1.25.2-alpine

ARG NAME
ARG PORT

ENV APP_NAME=$NAME
ENV APP_PORT=$PORT
ENV APP_ENV=development

RUN apk update && apk upgrade && \
	apk add --no-cache bash git openssh

RUN go install github.com/cespare/reflex@latest

WORKDIR /app/${APP_NAME}

COPY . .

ENTRYPOINT reflex -r '(\.go$|go\.mod$|\.yaml)' -s go run ./app/${APP_NAME}/cmd