package main

import (
	"dao"
	"dto"
	"fmt"
	"log"
)

//Before you execute the program, Launch `cqlsh` and execute:
//create keyspace example with replication = { 'class' : 'SimpleStrategy', 'replication_factor' : 1 };
//create table example.tweet(timeline text, id UUID, text text, PRIMARY KEY(id));
//create index on example.tweet(timeline);
func main() {
	session, err := dao.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	tw := dto.NewTweet("Jerome LAFORGE", "Hello world")
	tw.Insert(session)

	binding := tw.Select(session, "Jerome LAFORGE")
	defer binding.Close()

	for tw.Next(binding) {
		fmt.Println(tw)
	}
}
