package config

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
)

var path = flag.String("c", "config.json", "Path to configuration file")

var config = new(Config)

func init() {
	flag.Parse()

	f, err := ioutil.ReadFile(*path)
	if err != nil {
		if abs, err := filepath.Abs(*path); err == nil {
			log.Printf("Unable to read config file [%s]: %v", abs, err)
		} else {
			log.Printf("Unable to read config file: %v", err)
		}
	}

	err = json.Unmarshal(f, config)
	if err != nil {
		log.Fatalln("Error parsing config file:", err)
	}
}

func Get() *Config {
	return config
}

type Config struct {
	// Databases
	ContentDB  string // path to content database
	LanguageDB string // path to language database
	UserDB     string // path to user database

	// Web
	HttpServer   string // main http server
	TemplateRoot string // path to templates files
	StaticDir    string // optional: if present, the server also serves static files from there

	// Optional, only for distributed version over rpc:
	Frontend   string // frontend server
	Content    string // content server
	User       string // user server
	Language   string // language server
	ContentOp  string // content operations server
	UserOp     string // user operations server
	LanguageOp string // language operations server
}
