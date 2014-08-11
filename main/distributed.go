package main

import (
	"flag"
	"log"

	"github.com/polyglottis/platform"
	"github.com/polyglottis/platform/backend"
	_ "github.com/polyglottis/platform/config"
	contentRpc "github.com/polyglottis/platform/content/rpc"
	frontendRpc "github.com/polyglottis/platform/frontend/rpc"
	languageRpc "github.com/polyglottis/platform/language/rpc"
	userRpc "github.com/polyglottis/platform/user/rpc"
)

var addr = flag.String("http", ":8080", "Address for the http server")

var frontendAddr = flag.String("frontend", ":18658", "Address of the frontend server")
var contentAddr = flag.String("content", ":18982", "Address of the content server")
var userAddr = flag.String("user", ":14773", "Address of the user server")
var languageAddr = flag.String("language", ":14342", "Address of the language server")

func main() {

	log.Println("Configuring Polyglottis...")

	frontend, err := frontendRpc.NewClient(*frontendAddr)
	if err != nil {
		log.Fatalln(err)
	}

	content, err := contentRpc.NewClient(*contentAddr)
	if err != nil {
		log.Fatalln(err)
	}

	user, err := userRpc.NewClient(*userAddr)
	if err != nil {
		log.Fatalln(err)
	}

	language, err := languageRpc.NewClient(*languageAddr)
	if err != nil {
		log.Fatalln(err)
	}

	config := &platform.Configuration{
		Frontend: frontend,

		Backend: &backend.Configuration{
			Content:  content,
			User:     user,
			Language: language,
		},
	}

	log.Println("Launching Polyglottis...")

	log.Fatal(platform.Launch(*addr, config))
}
