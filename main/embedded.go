package main

import (
	"log"

	contentServer "github.com/polyglottis/content_server/server"
	"github.com/polyglottis/frontend_server"
	languageServer "github.com/polyglottis/language_server/server"
	"github.com/polyglottis/platform"
	"github.com/polyglottis/platform/backend"
	"github.com/polyglottis/platform/config"
)

func main() {

	log.Println("Configuring Polyglottis...")

	c := config.Get()

	content, err := contentServer.NewServer(c.ContentDB)
	if err != nil {
		log.Fatalln(err)
	}

	lang, err := languageServer.NewServer(c.LanguageDB)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Launching Polyglottis...")

	log.Fatal(platform.Launch(c.HttpServer, &platform.Configuration{
		Frontend: frontend_server.New(),

		Backend: &backend.Configuration{
			Content:  content,
			Language: lang,
		},
	}))
}
