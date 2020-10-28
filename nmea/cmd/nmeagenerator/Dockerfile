FROM golang:alpine
RUN apk add git

RUN go get github.com/RaaLabs/shipsimulator/nmea

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .
RUN adduser -S -D -H -h /app appuser
USER appuser