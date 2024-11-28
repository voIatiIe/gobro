FROM golang:1.21-alpine3.18

RUN apk add --no-cache bash
RUN apk add chromium

COPY . /server
WORKDIR /server

RUN go install

ENTRYPOINT ["go", "run", "main.go"]
