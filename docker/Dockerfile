FROM golang:latest

RUN mkdir -p /go/src/logflow

WORKDIR /go/src/logflow

#COPY . /go/src/logflow

ENV PORT 8080

EXPOSE 8080

CMD ["sh", "-c", "go run server/udp_server.go; echo 'started udp server'; go run main.go --action serve"]
