package config

import (
	"log"
	"path/filepath"
	"strings"
)

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	log.Println("Executing main configuration")
	defer log.Println("Main configuration executed")

	basePath, err := getBasePath()
	if err != nil {
		log.Fatalln(err)
	}

	Default = &Config{
		basePath: basePath,
	}
}

func getBasePath() (string, error) {
	abs, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}

	// TODO this is clearly not great :(
	return strings.SplitAfter(abs, "polyglottis")[0] + "/", nil
}

type Config struct {
	basePath string
}

var Default *Config

func (c *Config) TemplatePath() string {
	return c.basePath + "html/templates/"
}
