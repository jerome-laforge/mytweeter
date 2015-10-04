package dto

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/relops/cqlr"
)

type Tweet struct {
	Timeline string     `cql:"timeline" json:"timeline"`
	Id       gocql.UUID `cql:"id" json:"id"`
	Text     string     `cql:"text" json:"text"`
}

func NewTweet(timeLine, text string) (tw *Tweet) {
	tw = new(Tweet)
	tw.Timeline = timeLine
	tw.Text = text
	tw.Id = gocql.TimeUUID()
	return
}

func (this Tweet) Insert(session *gocql.Session) error {
	b := cqlr.Bind(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`, this)
	return b.Exec(session)
}

func (this Tweet) Select(session *gocql.Session, timeLine string) *cqlr.Binding {
	q := session.Query(`SELECT text, id, timeline FROM tweet WHERE timeline = ?`, timeLine)
	return cqlr.BindQuery(q)
}

func (this *Tweet) Next(bind *cqlr.Binding) bool {
	return bind.Scan(this)
}

func (this Tweet) String() string {
	return fmt.Sprint("Tweet:", this.Id, this.Text, this.Timeline)
}
