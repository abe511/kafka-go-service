## kafka go service

a pipeline which let's you send a message to kafka, store it in a db, consume the message and update its status in the db

---
send a message:
```
curl -X POST -H "Content-Type: application/json" -d '{"content": "test message"}' https://kafka-go-service.onrender.com/message
```

[**see stats here**](https://kafka-go-service.onrender.com/stats)

---
Deploy locally by running:

`docker compose up`\
this will spin up three containers:
- the service
- postgres
- kafka

*you may need to add **sudo** before the docker command*

### Usage

To produce a Kafka message
send a POST request to\
`localhost:8080/message`

with a json formatted text\
example:
```json
{"content": "test message"}
```

The response should look like:
```json
{"id":1,"content":"test message","processed":false}
```
The response is sent back before the message gets processed by Kafka.


To see the total number of sent and processed message\
send a GET request to\
`localhost:8080/stats`

response:
```json
{"total_messages":1,"processed_messages":1}
```

### Clean up

`docker compose down`\
or Ctrl+C to shut down the containers

Remove the containers:\
`docker rm kafka-go-service postgres kafka`

Remove the image:\
`docker rmi kafka-go-service`
