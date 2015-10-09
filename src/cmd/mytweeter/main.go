package main

import (
	"os"
	"www"
)

//sudo docker pull spotify/cassandra
//sudo docker run --name cassandra -d -p 9042:9042 spotify/cassandra
//sudo docker exec -it cassandra bash

//Before you execute the program, Launch `cqlsh` and execute:
//create keyspace example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
//create table example.tweet(timeline text, id UUID, text text, PRIMARY KEY(id));
//create index on example.tweet(timeline);

//create table example.user(id UUID, login text, passwd text, PRIMARY KEY(id));
//create index on example.user(login);
//insert into example.user(id, login, passwd) values (now(), 'admin', 'f807c2b4caa8ca621298907e5372c975a6e07322');
func main() {
	www.StartWebServer()
	os.Exit(0)
}
