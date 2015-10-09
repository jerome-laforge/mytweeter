package main

import (
	"dao"
	"os"
	"www"

	"github.com/inconshreveable/log15"
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
	log15.Root().SetHandler(log15.CallerFileHandler(log15.StdoutHandler))
	log := log15.New("module", "main")
	session, err := dao.GetSession()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	defer session.Close()
	www.StartWebServer()
	os.Exit(0)
}
