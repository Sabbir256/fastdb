# fastdb
An in-memory key/value store implementation using Golang. Data is persisted to disk every second using Append only File(AOF) technique.

Aceepts commands in Redis Serialization Protocol (RESP) format, and also respones using the same. That's why this server can communicate with the redis client server.

## Setup

Install redis on your system and stop the redis service to ensure that the redis client talks to our server only.
```bash
brew services stop redis
```

After that clone this repositiory and run the `main.go` file. If it does not work, you can try,
```bash
go run *.go
```

Now, we need to run the redis client to talk to our server
```bash
redis-cli
```
run the above command from your terminal to fire up the redis cli client and start testing redis commands.
```bash
127.0.0.1:6379> set name John
OK
127.0.0.1:6379> get name
"John"
```
