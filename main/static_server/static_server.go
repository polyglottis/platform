package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

type Config struct {
	ProductionLog bool   // If true, date and full file name is also on log
	Port          string // Port of the file server
	PathPrefix    string
	MountPrefix   string
	Dirs          []struct {
		Path  string // on file system
		Mount string // mounting point on server
	}
}

var configPath = flag.String("config", "server_config.json", "Path to configuration file")

func main() {
	log.SetFlags(log.Lshortfile | log.Ltime)

	f, err := ioutil.ReadFile(*configPath)
	if err != nil {
		log.Fatalln(err)
	}

	var config Config
	err = json.Unmarshal(f, &config)
	if err != nil {
		log.Fatalln(err)
	}

	if config.ProductionLog {
		log.SetFlags(log.Flags() | log.Ldate | log.Llongfile)
	}

	log.Printf("Launching Polyglottis Prototype with config %+v", config)

	for _, dir := range config.Dirs {
		fsPath := path.Join(config.PathPrefix, dir.Path)
		fs := http.FileServer(http.Dir(fsPath))
		prefix := path.Join("/", config.MountPrefix, dir.Mount)
		if prefix[len(prefix)-1] != '/' {
			prefix += "/"
		}
		http.Handle(prefix, http.StripPrefix(prefix, fs))
		log.Printf("Mounting folder %s at %s", fsPath, prefix)
	}

	log.Printf("Launching server on port %s", config.Port)

	log.Fatal(http.ListenAndServe(config.Port, nil))
}
