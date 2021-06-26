FROM golang:latest

RUN mkdir -p /go/src/logflow

WORKDIR /go/src/logflow

#COPY . /go/src/logflow

ENV PORT 8080

EXPOSE 8080
EXPOSE 5053
EXPOSE 6060/udp

CMD ["sh", "-c", "go run main.go --action kafka:topics && go run main.go --action serve & go run server/udp_server.go; echo 'started udp server';"]
