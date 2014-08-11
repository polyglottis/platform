package main

import (
	"flag"
	"log"

	contentServer "github.com/polyglottis/content_server/server"
	"github.com/polyglottis/frontend_server"
	"github.com/polyglottis/platform"
	"github.com/polyglottis/platform/backend"
	_ "github.com/polyglottis/platform/config"
)

var addr = flag.String("http", ":8080", "Address for the http server")
var contentDB = flag.String("content", "content.db", "Path to content database")

func main() {

	log.Println("Configuring Polyglottis...")

	content, err := contentServer.NewServer(*contentDB)
	if err != nil {
		log.Fatal(err)
	}

	config := &platform.Configuration{
		Frontend: frontend_server.New(),

		Backend: &backend.Configuration{
			Content: content,
		},
	}

	log.Println("Launching Polyglottis...")

	log.Fatal(platform.Launch(*addr, config))
}
