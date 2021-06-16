## nlogx

This is a demo project, it will collect nginx layer 7 data and layer 4 packet data and will be sent to kafka using protobuf.

## Todos 
- [ ] Installations & setup
    - [ ] Setup kafka ( probably using docker )
    - [ ] Decide on protobuf / grpc / rpc
    - [ ] Setup a nginx docker 

- [ ] Enable .env     
- [ ] Connect to kafka
- [ ] Setup a kafka client to receive messages
- [ ] Setup a kafka producer 
- [ ] Write tests
- [ ] Get some log data and test the whole system




## Start Kafka
* Kafka uses ZooKeeper as a distributed backend.

### Start Zookeeper
```
bin/zookeeper-server-start.sh config/zookeeper.properties
```

### Start Kafka
```
bin/kafka-server-start.sh config/server.properties
```

## Topics

### Create Topic
```
bin/kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test
```

### List Topics
```
bin/kafka-topics.sh --list --zookeeper localhost:2181
```

## Messages
### Send Message
```
bin/kafka-console-producer.sh --broker-list localhost:9092 --topic test
```


## Consumers
### Start Consumer
```
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning
```




