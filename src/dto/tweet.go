package dto

import (
	"dao"
	"fmt"

	"github.com/gocql/gocql"
	"github.com/relops/cqlr"
)

type Tweet struct {
	Timeline string     `cql:"timeline" json:"timeline"`
	Id       gocql.UUID `cql:"id"       json:"id"`
	Text     string     `cql:"text"     json:"text"`
}

func (this *Tweet) Insert() error {
	session, err := dao.GetSession()
	if err != nil {
		return err
	}

	b := cqlr.Bind(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`, this)
	return b.Exec(session)
}

func (this Tweet) String() string {
	return fmt.Sprintf(`Tweet id: "%s" text: "%s" timeline: "%s"`, this.Id, this.Text, this.Timeline)
}

func (this *Tweet) GenerateId() {
	this.Id = gocql.TimeUUID()
}

func NewTweet(timeLine, text string) (tw *Tweet) {
	tw = new(Tweet)
	tw.Timeline = timeLine
	tw.Text = text
	tw.Id = gocql.TimeUUID()
	return
}

func GetAllTweetsForTimeLine(timeLine string) ([]Tweet, error) {
	session, err := dao.GetSession()
	if err != nil {
		return nil, err
	}

	q := session.Query(`SELECT text, id, timeline FROM tweet WHERE timeline = ?`, timeLine)
	bind := cqlr.BindQuery(q)
	defer bind.Close()

	var tweets []Tweet
	t := Tweet{}
	for bind.Scan(&t) {
		tweets = append(tweets, t)
	}
	return tweets, nil
}
