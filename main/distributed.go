package main

import (
	"log"

	"github.com/polyglottis/platform"
	"github.com/polyglottis/platform/backend"
	"github.com/polyglottis/platform/config"
	contentRpc "github.com/polyglottis/platform/content/rpc"
	frontendRpc "github.com/polyglottis/platform/frontend/rpc"
	languageRpc "github.com/polyglottis/platform/language/rpc"
	userRpc "github.com/polyglottis/platform/user/rpc"
)

func main() {

	log.Println("Configuring Polyglottis...")

	c := config.Get()

	frontend, err := frontendRpc.NewClient(c.Frontend)
	if err != nil {
		log.Fatalln(err)
	}

	content, err := contentRpc.NewClient(c.Content)
	if err != nil {
		log.Fatalln(err)
	}

	user, err := userRpc.NewClient(c.User)
	if err != nil {
		log.Fatalln(err)
	}

	language, err := languageRpc.NewClient(c.Language)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Launching Polyglottis...")

	log.Fatal(platform.Launch(c.HttpServer, &platform.Configuration{
		Frontend: frontend,

		Backend: &backend.Configuration{
			Content:  content,
			User:     user,
			Language: language,
		},
	}))
}
