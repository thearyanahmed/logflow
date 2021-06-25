## Logflow 

### What it does?
Logflow collects nginx layer 7 data and layer 4 packet data and will be sent to kafka using protobuf.
As of using protocol buffer, it works very fast. We will be improving on it. It is still in very early stage.

It takes the data from the gRPC request and uses publishes into kafka . Logflow uses `client stream` to collect the data.

It can publish a single request into multiple topics and topics can be changed on every request, as the time of writing this,
all possible topics must be given (set in the `.env` file) before you start the server. 

### Dependencies
1. [Kafka](https://kafka.apache.org/downloads)
2. GO
3. [Protocol Buffer](https://grpc.io/docs/protoc-installation/) for development purpose.

*At the moment we don't have the docker environments ready yet. The dependencies need to be installed manually.*

## Running the Application

First, git clone

```bash
// https 
git clone https://github.com/thearyanahmed/logflow.git

// or ssh
git clone git@github.com:thearyanahmed/logflow.git

// github cli
gh repo clone thearyanahmed/logflow
```
then cd into the directory

```bash
cd logflow
```

If you don't have ZooKeeper and Kafka running,
[you can follow these instructions](https://github.com/thearyanahmed/logflow#start-zookeeper-and-kafka "Starting kafka")



### .ENV
Once you have ZooKeeper and Kafka running, in the `logflow` directory

```bash
cp .env.example to .env
```
Make sure you have the .env values setup correctly, in case you have changed any config for kafka or 
if some default ports/values are already in use on your machine.

To start the server, run 

```bash
go run main.go --action serve
```

This is start tcp connect server at `RPC_PORT` ( from .env) . The default is 5053.

To start the client, you can type

```bash
go run main.go --action client
```

**Note** Before running the client, make sure you have ZooKeeper and Kafka running.

If everything was successful, you'll have an output like

![Consumer](images/consumer.png?raw=true "Kafka Consumers")

**Note** The program is at an very early stage.


### Architecture
![Logflow Architecture](images/logflow-architecture.png?raw=true "Logflow Architecture")


#### Start ZooKeeper and Kafka

First run ZooKeeper and Kafka brokers. Go to your kafka installation directory and run

```
bin/zookeeper-server-start.sh config/zookeeper.properties
```

```
bin/kafka-server-start.sh config/server.properties
```

If this is your first time running them, you probably don't have topics. You'll need at least one.
To create one, use the following command.

It will create a topic name `hello_world`
```bash
bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic hello_world
```

- Running kafka consumer (optional)
  In case you want to see the data transmission directly from kafka.

```bash
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic hello_world --from-beginning```
```

### With Docker

At the moment I have not docker image ready for this but soon will.
If you want to run kafka inside docker, you can simply use your docker
container's address eg: localhost:39092 or something like that ( in the .env) .

### Flags

There is at the moment only one flag you can pass along when running the main.go file.

1.  `--action` Possible values: serve, client, help. eg: `go run main.go --action serve`


### Useful Resources

[Enterprise Network Flow Collector (IPFIX, sFlow, Netflow) from Verizon Media](https://github.com/VerizonDigital/vflow)

[The high-scalability sFlow/NetFlow/IPFIX collector used internally at Cloudflare](https://github.com/cloudflare/goflow)

[GopherCon 2016: John Leon - Packet Capture, Analysis, and Injection with Go](https://www.youtube.com/watch?v=APDnbmTKjgM)

[Capturing HTTP packets the hard way](https://medium.com/@cjoudrey/capturing-http-packets-the-hard-way-b9c799bfb6)

[LISA16: Linux 4.X Tracing Tools: Using BPF Superpowers](https://www.youtube.com/watch?v=UmOU3I36T2U)

[Sniffing Creds with Go, A Journey with libpcap](https://itnext.io/sniffing-creds-with-go-a-journey-with-libpcap-73bc3e74966)

[Collecting NGINX Plus Monitoring Statistics with Go
](https://www.nginx.com/blog/collecting-nginx-plus-monitoring-statistics-with-go/)

[GoPacket by Google](https://pkg.go.dev/github.com/google/gopacket)

[Packet Capture, Injection, and Analysis with Go Packet](https://www.devdungeon.com/content/packet-capture-injection-and-analysis-gopacket)

[BCC HTTP Filter](https://github.com/iovisor/bcc/tree/master/examples/networking/http_filter)

[Logkit by Qiniu](https://github.com/qiniu/logkit)

[NGINX log data format](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_formata)

These ^ are gems

## Todos 
- [ ] Installations & setup
    - [x] Setup kafka
    - [ ] Use docker container for kafka
    - [x] Decide on protobuf ( grpc / rpc )
    - [ ] Setup a nginx docker
- [x] Enable .env support     
- [ ] Draw system architecture  
- [x] Connect to kafka
- [x] Setup a kafka client to receive messages
- [x] Setup a kafka producer 
- [ ] Write tests
- [ ] Get some log data and test the whole system
- [ ] Dockerize full app
- [x] Prepare dummy NginxLogRequest, Packet & Headers
- [x] Prepare a client to test grpc and kafka producer with dummy data
- [x] Also see manually if kafka consumer is consuming messages
- [ ] Document everything

#### At the moment, I did not include any docker image for any component.

#### More to come
