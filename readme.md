## nlogx

This is a demo project, it will collect nginx layer 7 data and layer 4 packet data and will be sent to kafka using protobuf.

## Todos 
- [ ] Installations & setup
    - [x] Setup kafka
    - [ ] Use docker container for kafka
    - [x] Decide on protobuf ( grpc / rpc )
    - [ ] Setup a nginx docker
- [x] Enable .env support     
- [x] Connect to kafka
- [x] Setup a kafka client to receive messages
- [x] Setup a kafka producer 
- [ ] Write tests
- [ ] Get some log data and test the whole system
- [ ] Dockerize full app
- [x] Prepare dummy NginxLogRequest, Packet & Headers
- [x] Prepare a client to test grpc and kafka producer with dummy data
- [x] Also see manually if kafka consumer is consuming messages 


#### At the moment, I did not include any docker image for any component.
#### Full app will be dockerized.

## Running the application

### Manually 

First run ZooKeeper and Kafka brokers. Go to your kafka directory and run 

- Start Zookeeper
```
bin/zookeeper-server-start.sh config/zookeeper.properties
```

- Start Kafka
```
bin/kafka-server-start.sh config/server.properties
```

- Then create a topic with 

```bash
bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic $topicName
```

- you can check your created topics with 
```bash
bin/kafka-topics.sh --list --zookeeper localhost:2181
```


### With Docker

At the moment I have not docker image ready for this but soon will.
If you want to run kafka inside docker, you can simply use your docker
container's address eg: localhost:39092 or something like that ( in the .env) .


- To run the app, you can run ```go run main.go``` . 
  
The app will start listening ,

You can run ```go run kafka/kafka_producer.go``` to generate random strings
and stream it to your brokers. 

Your consumer console should have an output like this 

![Consumer](images/consumer.png?raw=true "Kafka Consumers")


#### Flags

There are a few flags you can pass along when running the main.go file.
Those are, 

- *kafka-brokers* these are comma separated values for your brokers. eg: *localhost:9092,localhost:39093* . It has a default value of localhost:9092 
- *kafka-topic* eg: hello_world. Default hello_world
- *kafka-client-id* eg: test_client_id. Default will be a random string kafka_client_XYZ . XYZ is randomly generated string.

#### More to come


### Sources

[nginx log data format](http://nginx.org/en/docs/http/ngx_http_log_module.html#log_formata)