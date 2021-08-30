FROM golang:1.16

WORKDIR /app
COPY . .
COPY go.mod ./
RUN apt-get update
RUN apt-get install -y uuid-runtime


RUN go mod download
RUN go get -t
RUN go build
CMD [ "./aut" ]
