FROM golang:1.16
RUN apt update
COPY . .
RUN go mode
RUN go get
