FROM golang:1.16-alpine

WORKDIR /app
COPY . .
COPY go.mod ./
RUN apt-get update
RUN apt-get install uuid-runtime
RUN go mod download
RUN go get -t
RUN go build
CMD [ "./aut" ]
