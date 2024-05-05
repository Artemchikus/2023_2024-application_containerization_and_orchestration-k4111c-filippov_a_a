package main

import (
	"find-ship/cmd/app"
	"log"
)

func main() {
	var err error
	if err = app.Execute(); err != nil {
		log.Println(err)
	}
}
