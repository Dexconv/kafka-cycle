this is a test project showcasing kafka producing and consuming cycles.

for the sake of simplicity, both these services have been implemented in the same app to avoid external storage.

in this project 2 services have been mocked:

- data producer to kafka.
- data consumer from kafka.

there is a queue (of type []byte and cap of 3075) being filled periodically, and it gets emptied if one of two conditions are met:

- the new data and existing queue data exceeds the predetermined 3075 length.
- the length of queue is a multiply of 1024.

after the condition is met, the queue data is split in 4 parts and produced to kafka partitions, then they are read and conjoined.

the conjoined data and origionally produced data are then compared to show no data was lost from kafka.

## steps to run the project:

_please use an empty topic for each run_

#### 1. start kafka using docker compose

```bash
docker compose up -d
```

#### 2. create a topic with 4 partitions

```bash
bin/kafka-topics.sh --create --topic test --bootstrap-server localhost:29092 --partitions 4
```

#### 3. run the project

```bash
go run .
```

#### 4. starting over

preferably down kafka and start from the top, alteratively use a new topic

```bash
docker compose down
```

### output:

logs if the message has been verified, and the shortened hexadecimal representations of each compared []byte

```bash
2023/10/03 18:34:46 message is verified
2023/10/03 18:34:46 sent: 2a314d6afa...282d9a9e16 , recieved: 2a314d6afa...282d9a9e16
```
