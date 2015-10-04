package main

import (
	"log"

	"dao"
	"dto"
	"fmt"
)

/*type Tweet struct {
	Timeline string     `cql:"timeline"`
	Id       gocql.UUID `cql:"id"`
	Text     string     `cql:"text"`
}

func (this Tweet) String() string {
	return fmt.Sprintln("Tweet:", this.Id, this.Text)
}

func main() {
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "example"

	cluster.Consistency = gocql.Quorum

	session, err := dao.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	if err = session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`, "me", gocql.TimeUUID(), "hello world").Exec(); err != nil {
		log.Fatal(err)
	}

	ids := make([]gocql.UUID, 0)
	var id gocql.UUID
	var text string

*/ /* Search for a specific set of records whose 'timeline' column matches
 * the value 'me'. The secondary index that we created earlier will be
 * used for optimizing the search */ /*
	if err := session.Query(`SELECT id, text FROM tweet WHERE timeline = ? LIMIT 1`, "me").Consistency(gocql.One).Scan(&id, &text); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Tweet:", id, text)

	iter := session.Query(`SELECT id, text FROM tweet WHERE timeline = ?`, "me").Iter()
	for iter.Scan(&id, &text) {
		ids = append(ids, id)
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	const insertMsg = "hello world wrong 5"
	//	if err = session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`, "me", ids[rand.Intn(len(ids))], insertMsg).Exec(); err != nil {
	if err = session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?) if not exists`, "me", ids[rand.Intn(len(ids))], insertMsg).Exec(); err != nil {
		log.Fatal(err)
	}

	// list all tweets
	iter = session.Query(`SELECT id, text FROM tweet WHERE timeline = ?`, "me").Iter()
	i := 0
	for iter.Scan(&id, &text) {
		fmt.Println("Tweet:#", i, id, text)
		ids = append(ids, id)
		i++
	}
	if err := iter.Close(); err != nil {
		log.Fatal(err)
	}

	q := session.Query(`SELECT text, id, timeline FROM tweet WHERE timeline = ?`, "me")
	b := cqlr.BindQuery(q)

	var t Tweet
	for b.Scan(&t) {
		// Application specific code goes here
		fmt.Print(t)
	}
}*/

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
