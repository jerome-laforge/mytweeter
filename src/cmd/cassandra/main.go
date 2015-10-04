package main

import (
	"dao"
	"dto"
	"fmt"
	"log"
)

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
