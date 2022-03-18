FROM golang:1.18.0-alpine3.14
RUN mkdir /app
ADD . /app
WORKDIR /app/pkg/
RUN apk add git
RUN apk add --no-cache gcc musl-dev
RUN go build  -o main . 
